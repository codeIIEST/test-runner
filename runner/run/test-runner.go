package run

import (
	"log"

	"github.com/codeiiest/test-runner/internal/container"
	"github.com/codeiiest/test-runner/runner/tester"
	"github.com/docker/docker/client"
)

var cli *client.Client

func init() {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := container.BuildImage(cli, "runner:latest")

	if err != nil {
		log.Fatal(err)
	}
}

// Evaluate runs the program against given test cases
// and returns the result as a TestData struct
func Evaluate(code string, lang string, file string, in []string, out []string,
	count int, time int, mem int64) tester.TestResult {

	data := tester.TestData{
		ID:         "1",
		Lang:       lang,
		Code:       code,
		Filename:   file,
		InputData:  in,
		OutputData: out,
		TestCount:  count,
		TimeLimit:  time,
		MemLimit:   mem,
	}
	return data.Run(cli)
}
