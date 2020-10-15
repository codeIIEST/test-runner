package container

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	ct "github.com/raydwaipayan/test-runner/internal/types"
)

func runContainer(ctx context.Context, cli *client.Client, ID string) (string, error) {
	if err := cli.ContainerStart(ctx, ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	res := buf.String()

	cli.ContainerRemove(ctx, ID, types.ContainerRemoveOptions{})
	return res, err
}

func runCpp(ctx context.Context, cli *client.Client, image string, path string, count int, time int64, memory int64) error {
	eval := fmt.Sprintf("/tests/evaluate %v %v %v", count, time, memory)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", eval},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/tests/data",
			},
		},
	}, nil, nil, "")

	if err != nil {
		return err
	}

	_, err = runContainer(ctx, cli, resp.ID)
	return err
}

func compileCpp(ctx context.Context, cli *client.Client, image string, path string) (int, string) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", "g++ -w -O2 /tests/data/a.cpp -o /tests/data/a.out 2>&1"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/tests/data",
			},
		},
	}, nil, nil, "")

	if err != nil {
		return 2, err.Error()
	}

	out, err := runContainer(ctx, cli, resp.ID)

	if err != nil {
		return 2, err.Error()
	}

	r, _ := regexp.Compile("error")
	if r.MatchString(out) {
		return 0, out
	}

	return 1, ""
}

// RunTests runs tests on the given data and returns evaluation result
func RunTests(data ct.TestData) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ctx := context.Background()

	if err != nil {
		log.Print(err)
		return
	}
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "tests", data.ID)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Print(err)
		}
	}

	f, err := os.Create(filepath.Join(path, data.Filename))
	if err != nil {
		log.Print(err)
		return
	}
	f.WriteString(data.Code)
	f.Close()

	for i := 1; i <= data.TestCount; i++ {
		f, err := os.Create(filepath.Join(path, fmt.Sprintf("in%v.txt", i)))
		if err != nil {
			log.Print(err)
			return
		}
		f2, err := os.Create(filepath.Join(path, fmt.Sprintf("out%v.txt", i)))
		if err != nil {
			log.Print(err)
			return
		}
		f.WriteString(data.InputData[i-1])
		f2.WriteString(data.OutputData[i-1])

		f.Close()
		f2.Close()
	}

	switch data.Lang {
	case "cpp":
		status, mesg := compileCpp(ctx, cli, "runner:latest", path)
		log.Print(status, mesg)
		err := runCpp(ctx, cli, "runner:latest", path, data.TestCount,
			data.TimeLimit, data.MemLimit)
		log.Print(err)
	}
}
