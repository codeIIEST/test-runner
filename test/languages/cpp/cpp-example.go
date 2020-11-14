package cpp

import (
	"io/ioutil"
	"log"

	"github.com/codeiiest/test-runner/runner/run"
	util "github.com/codeiiest/test-runner/test/utils"
)

// Test is used to test CPP
func Test() {

	dat, _ := ioutil.ReadFile("samples/cpp/a.cpp")
	in := []string{"2", "4", "5"}
	out := []string{"4", "16", "25"}
	util.GetFilesInPwd()
	res := run.Evaluate(string(dat), "cpp", "a.cpp", in, out, 3, 1, 500*1024*1024)
	log.Print(res)
}
