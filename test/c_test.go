package tests

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/codeiiest/test-runner/runner/run"
	"github.com/codeiiest/test-runner/test/utils"
)

func TestC(t *testing.T) {
	t.Helper()
	t.Run("C test with all ACCEPTED", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("testdata/c/simple.c")
		testIn := []string{"1", "2", "6969"}
		testOut := []string{"1", "4", "48566961"}
		var testCases int = 3
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "c", "a.c", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		log.Print(got, "\n", got.Status, got.Message, got.Error, got.Result)
		wantStatus := "OK"
		wantMessage := ""
		wantError := []string{"", "", ""}
		wantResult := []string{"ACCEPTED", "ACCEPTED", "ACCEPTED"}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})

	t.Run("C test with 1 WRONG ANSWER", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("testdata/c/simple.c")
		testIn := []string{"12", "2", "1"}
		testOut := []string{"48566961", "4", "1"}
		var testCases int = 3
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "c", "a.c", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		log.Print(got, "\n", got.Status, got.Message, got.Error, got.Result)
		wantStatus := "OK"
		wantMessage := ""
		wantError := []string{"", "", ""}
		wantResult := []string{"WRONG ANSWER", "ACCEPTED", "ACCEPTED"}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})

	t.Run("C test with COMPILATION ERROR", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("testdata/c/compilationerror.c")
		testIn := []string{"12", "2", "1"}
		testOut := []string{"48566961", "4", "1"}
		var testCases int = 3
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "c", "a.c", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		log.Print(got)
		log.Print(got, "\n", got.Status, got.Message, got.Error, got.Result)
		wantStatus := "COMPILATION ERROR"
		wantMessage := ""
		wantError := []string{}
		wantResult := []string{}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})
}
