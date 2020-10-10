package container

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func buildImage(ctx context.Context, cli *client.Client) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := "Dockerfile"
	cwd, _ := os.Getwd()
	dockerFileReader, err := os.Open(cwd + "/internal/container/dockerfile/Dockerfile")
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
		return err
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
		return err
	}
	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
		return err
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
		return err
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	opt := types.ImageBuildOptions{
		Tags:       []string{"runner"},
		CPUSetCPUs: "1",
		Memory:     512 * 1024 * 1024,
		ShmSize:    64,
		Context:    dockerFileTarReader,
		Dockerfile: dockerFile,
	}
	resp, err := cli.ImageBuild(ctx, dockerFileTarReader, opt)
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

// CreateContainer builds several containers for running the code
func CreateContainer() {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	err := buildImage(ctx, cli)
	if err != nil {
		return
	}
}
