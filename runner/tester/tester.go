package tester

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/raydwaipayan/test-runner/internal/container"
)

// TestData contains all test data related to the program
// to be evaluated
type TestData struct {
	ID         string
	Lang       string
	Filename   string
	Code       string
	Path       string
	Image      string
	TestCount  int
	InputData  []string
	OutputData []string
	TimeLimit  int
	MemLimit   int64
}

// TestResult contains the evaluation result
type TestResult struct {
	ID      string
	Status  string
	Message string
	Time    []float64
	Memory  []float64
	Result  []string
	Error   []string
}

// Run executes tests on the given data and returns evaluation result
func (data *TestData) Run(cli *client.Client) TestResult {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ctx := context.Background()

	if err != nil {
		log.Print(err)
		return TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	cwd, _ := os.Getwd()
	data.Path = filepath.Join(cwd, "temp", data.ID)
	data.Image = "runner:latest"

	err = data.toFile()
	if err != nil {
		log.Print(err)
		return TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	status, mesg := container.Compile(ctx, cli, data.Image, data.Lang, data.Path)
	if status != 1 {
		log.Print(mesg)
		if status == 0 {
			return TestResult{
				ID:      data.ID,
				Status:  "COMPILATION ERROR",
				Message: mesg,
			}
		}
		return TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	err = container.Execute(ctx, cli, data.Image, data.Lang, data.Path,
		data.TestCount, data.TimeLimit, data.MemLimit)
	if err != nil {
		log.Print(err)
		return TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}

	res := TestResult{
		ID:     data.ID,
		Time:   make([]float64, data.TestCount),
		Memory: make([]float64, data.TestCount),
		Result: make([]string, data.TestCount),
		Error:  make([]string, data.TestCount),
		Status: "OK",
	}

	err = parseOutput(data.Path, data, &res)
	if err != nil {
		log.Print(err)
		return TestResult{
			ID:     data.ID,
			Status: "INTERNAL SERVER ERROR",
		}
	}
	err = os.RemoveAll(data.Path)
	if err != nil {
		log.Print(err)
	}
	return res
}

func (data *TestData) toFile() error {
	_, err := os.Stat(data.Path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(data.Path, 0755)
		if errDir != nil {
			return err
		}
	}

	f, err := os.Create(filepath.Join(data.Path, data.Filename))
	if err != nil {
		return err
	}
	f.WriteString(data.Code)
	f.Close()

	for i := 1; i <= data.TestCount; i++ {
		f, err := os.Create(filepath.Join(data.Path, fmt.Sprintf("in%v.txt", i)))
		if err != nil {
			return err
		}
		f2, err := os.Create(filepath.Join(data.Path, fmt.Sprintf("out%v.txt", i)))
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
