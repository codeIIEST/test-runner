package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func buildImage(cli *Client, ctx Context) (ImageBuildResponse, error) {

	opt := types.ImageBuildOptions{
		CPUSetCPUs: "1",
		CPUSetMems: "12",
		CPUShares:  20,
		CPUQuota:   10,
		CPUPeriod:  30,
		Memory:     256,
		ShmSize:    10,
		Dockerfile: "dockerfile/Dockerfile",
	}
	res, err := cli.ImageBuild(ctx, nil, opt)
	return res, err
}

// CreateContainer builds several containers for running the code
func CreateContainer() {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	res, err := buildImage(cli, ctx)
	fmt.Println(res)
}
