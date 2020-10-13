package container

import (
	"bytes"
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func runContainer(ctx context.Context, cli *client.Client, ID string) (string, error) {
	if err := cli.ContainerStart(ctx, ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
		return "INTERNAL_ERROR", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
			return "INTERNAL_ERROR", err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatal(err)
		return "INTERNAL_ERROR", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	res := buf.String()

	cli.ContainerRemove(ctx, ID, types.ContainerRemoveOptions{})
	return res, err
}

// RunCpp exported
func RunCpp(cli *client.Client, image string, path string, memory int64, time uint64) (string, error) {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"bash", "/tests/evaluate"},
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
		log.Fatal(err)
		return "INTERNAL_ERROR", err
	}

	out, err := runContainer(ctx, cli, resp.ID)
	return out, err
}

//CompileCpp exported
func CompileCpp(cli *client.Client, image string, path string) (string, error) {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", "g++ /tests/data/a.cpp -o /tests/data/a.out -O2 2>&1"},
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
		log.Fatal(err)
		return "INTERNAL_ERROR", err
	}

	out, err := runContainer(ctx, cli, resp.ID)
	return out, err
}
