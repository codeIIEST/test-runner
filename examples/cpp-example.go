package main

import (
	"io/ioutil"
	"log"

	"github.com/codeiiest/test-runner/runner/run"
)

func main() {

	dat, _ := ioutil.ReadFile("a.cpp")
	in := []string{"2", "4", "5"}
	out := []string{"4", "16", "25"}

	res := run.Evaluate(string(dat), "cpp", "a.cpp", in, out, 3, 1, 500*1024*1024)
	log.Print(res)
}
