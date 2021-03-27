package tests

import (
	"io/ioutil"
	"testing"

	"github.com/codeiiest/test-runner/runner/run"
	"github.com/codeiiest/test-runner/test/utils"
)

func TestJava(t *testing.T) {
	t.Helper()
	t.Run("Java test with all ACCEPTED", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("./testdata/java/Main.java")
		testIn := []string{" "}
		testOut := []string{"4"}
		var testCases int = 1
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "java", "Main.java", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		wantStatus := "OK"
		wantMessage := ""
		wantError := []string{""}
		wantResult := []string{"ACCEPTED"}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})
}
