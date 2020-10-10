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
		return "internal_err", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
			return "internal_err", err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatal(err)
		return "internal_err", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	res := buf.String()

	return res, err
}

// RunCpp exported
func RunCpp(cli *client.Client, image string, path string) (string, error) {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", "./a.out <in.txt >o.txt 2>&1 && diff o.txt out.txt"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/tests",
			},
		},
	}, nil, nil, "")

	if err != nil {
		log.Fatal(err)
		return "internal_err", err
	}

	out, err := runContainer(ctx, cli, resp.ID)
	return out, err
}

//CompileCpp exported
func CompileCpp(cli *client.Client, image string, path string) (string, error) {
	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   []string{"/bin/bash", "-c", "g++ a.cpp -O2 2>&1"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/tests",
			},
		},
	}, nil, nil, "")

	if err != nil {
		log.Fatal(err)
		return "internal_err", err
	}

	out, err := runContainer(ctx, cli, resp.ID)
	return out, err
}
