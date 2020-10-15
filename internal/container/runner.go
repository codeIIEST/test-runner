package container

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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

func execute(ctx context.Context, cli *client.Client, image string, lang string, path string, count int, time int64, memory int64) error {
	eval := ""
	switch lang {
	case "cpp":
		eval = fmt.Sprintf("/tests/evaluate %v %v %v", count, time, memory)
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

func compile(ctx context.Context, cli *client.Client, image string, lang string, path string) (int, string) {
	eval := ""
	switch lang {
	case "cpp":
		eval = "g++ -w -O2 /tests/data/a.cpp -o /tests/data/a.out 2>&1"
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

func createFiles(data *ct.TestData, path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			return err
		}
	}

	f, err := os.Create(filepath.Join(path, data.Filename))
	if err != nil {
		return err
	}
	f.WriteString(data.Code)
	f.Close()

	for i := 1; i <= data.TestCount; i++ {
		f, err := os.Create(filepath.Join(path, fmt.Sprintf("in%v.txt", i)))
		if err != nil {
			return err
		}
		f2, err := os.Create(filepath.Join(path, fmt.Sprintf("out%v.txt", i)))
		if err != nil {
			return err
		}
		f.WriteString(data.InputData[i-1])
		f2.WriteString(data.OutputData[i-1])

		f.Close()
		f2.Close()
	}
	return nil
}

func strToMap(data string) map[string]string {
	mp := make(map[string]string, 0)
	for _, line := range strings.Split(strings.TrimSuffix(data, "\n"), "\n") {
		idx := strings.Index(line, "=")
		mp[line[0:idx]] = line[idx+1:]
	}
	return mp
}

func readOutputFiles(path string, data *ct.TestData, res *ct.TestResult) error {
	for i := 1; i <= data.TestCount; i++ {
		dat, err := ioutil.ReadFile(filepath.Join(path, fmt.Sprintf("diff%v.txt", i)))
		diff := strings.TrimSpace(string(dat))
		if err != nil {
			return err
		}

		dat, err = ioutil.ReadFile(filepath.Join(path, fmt.Sprintf("stats%v.txt", i)))
		stats := string(dat)
		stmap := strToMap(stats)

		returnValue, _ := strconv.Atoi(stmap["returnvalue"])
		termination := strings.TrimSpace(stmap["terminationreason"])
		time, _ := strconv.ParseFloat(strings.TrimSuffix(stmap["cputime"], "s"), 64)
		mem, _ := strconv.ParseFloat(strings.TrimSuffix(stmap["memory"], "B"), 64)

		res.Time[i-1] = time
		res.Memory[i-1] = mem

		if returnValue != 0 {
			switch returnValue {
			case 9, 15:
				{
					if termination == "cputime" {
						res.Result[i-1] = "TIME LIMIT EXCEEDED"
					} else if termination == "memory" {
						res.Result[i-1] = "MEMORY LIMIT EXCEEDED"
					} else {
						res.Result[i-1] = "ILLEGAL INSTRUCTIONS"
					}
				}
			default:
				{
					res.Result[i-1] = "RUNTIME ERROR"
				}
			}
		} else {
			if diff == "" {
				res.Result[i-1] = "ACCEPTED"
			} else {
				res.Result[i-1] = "WRONG ANSWER"
			}
		}
	}
	return nil
}

// RunTests runs tests on the given data and returns evaluation result
func RunTests(data ct.TestData) ct.TestResult {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ctx := context.Background()

	if err != nil {
		log.Print(err)
		return ct.TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "temp", data.ID)

	err = createFiles(&data, path)
	if err != nil {
		log.Print(err)
		return ct.TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	status, mesg := compile(ctx, cli, "runner:latest", data.Lang, path)
	if status != 1 {
		log.Print(mesg)
		if status == 0 {
			return ct.TestResult{
				ID:      data.ID,
				Status:  "COMPILATION ERROR",
				Message: mesg,
			}
		}
		return ct.TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	err = execute(ctx, cli, "runner:latest", data.Lang, path, data.TestCount,
		data.TimeLimit, data.MemLimit)
	if err != nil {
		log.Print(err)
		return ct.TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	res := ct.TestResult{
		ID:     data.ID,
		Time:   make([]float64, data.TestCount),
		Memory: make([]float64, data.TestCount),
		Result: make([]string, data.TestCount),
		Error:  make([]string, data.TestCount),
		Status: "OK",
	}

	err = readOutputFiles(path, &data, &res)
	if err != nil {
		log.Print(err)
		return ct.TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}
	err = os.RemoveAll(path)
	if err != nil {
		log.Print(err)
	}
	return res
}
