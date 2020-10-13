package container

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func getContext(path string) io.Reader {
	ctx, _ := archive.TarWithOptions(path, &archive.TarOptions{})
	return ctx
}

// BuildImage creates the docker image to run the test programs
func BuildImage(cli *client.Client, imageName string) error {
	ctx := context.Background()
	if cli == nil {
		cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	}

	cwd, _ := os.Getwd()
	dockerCtx := getContext(cwd + "/internal/container/dockerfile")

	opt := types.ImageBuildOptions{
		Tags:       []string{imageName},
		CPUSetCPUs: "1",
		Memory:     512 * 1024 * 1024,
		ShmSize:    64,
		Dockerfile: "Dockerfile",
	}

	resp, err := cli.ImageBuild(ctx, dockerCtx, opt)
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
		return err
	}

	return err
}
