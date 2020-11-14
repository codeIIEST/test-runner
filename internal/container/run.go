package container

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
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

// Execute executes the test program against gives test cases
func Execute(ctx context.Context, cli *client.Client, image string, lang string,
	path string, count int, time int, mem int64) error {

	eval := ""
	switch lang {
	case "c":
		eval = fmt.Sprintf("/tests/evaluate %v %v %v", count, time, mem)
	case "cpp":
		eval = fmt.Sprintf("/tests/evaluate %v %v %v", count, time, mem)
	}

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
			mount.Mount{
				Type:   mount.TypeBind,
				Source: "/sys/fs/cgroup",
				Target: "/sys/fs/cgroup",
			},
		},
		Privileged: true,
		CapDrop:    []string{"all"},
	}, nil, nil, "")

	if err != nil {
		return err
	}

	_, err = runContainer(ctx, cli, resp.ID)
	return err
}

// Compile compiles the program and generates the executable
func Compile(ctx context.Context, cli *client.Client, image string, lang string,
	path string) (int, string) {

	eval := ""
	switch lang {
	case "c":
		eval = "gcc /tests/data/a.c -o /tests/data/a.out 2>&1"
	case "cpp":
		eval = "g++ -w -O2 /tests/data/a.cpp -o /tests/data/a.out 2>&1"
	}
	fmt.Println(eval)
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
