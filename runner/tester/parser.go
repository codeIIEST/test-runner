package tester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func strToMap(data string) map[string]string {
	mp := make(map[string]string, 0)
	for _, line := range strings.Split(strings.TrimSuffix(data, "\n"), "\n") {
		idx := strings.Index(line, "=")
		mp[line[0:idx]] = line[idx+1:]
	}
	return mp
}

func parseOutput(path string, data *TestData, res *TestResult) error {
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
