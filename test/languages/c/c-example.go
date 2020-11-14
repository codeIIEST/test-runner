package c

import (
	"io/ioutil"
	"log"

	"github.com/codeiiest/test-runner/runner/run"
)

// Test is used to test C
func Test() {
	dat, _ := ioutil.ReadFile("samples/c/a.c")
	in := []string{""}
	out := []string{"20"}

	res := run.Evaluate(string(dat), "c", "a.c", in, out, 1, 1, 500*1024*1024)
	log.Print(res)
}
