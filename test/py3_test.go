package tests

import (
	"io/ioutil"
	"testing"

	"github.com/codeiiest/test-runner/runner/run"
	"github.com/codeiiest/test-runner/test/utils"
)

func TestPython3(t *testing.T) {
	t.Helper()
	t.Run("Python3 test with all ACCEPTED", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("./testdata/python/simple.py")
		testIn := []string{"1", "2", "6969"}
		testOut := []string{"1", "4", "48566961"}
		var testCases int = 3
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "py", "a.py", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		wantStatus := "OK"
		wantMessage := ""
		wantError := []string{"", "", ""}
		wantResult := []string{"ACCEPTED", "ACCEPTED", "ACCEPTED"}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})
}
