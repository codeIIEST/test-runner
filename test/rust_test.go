package tests

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/codeiiest/test-runner/runner/run"
	"github.com/codeiiest/test-runner/test/utils"
)

func TestRust(t *testing.T) {
	t.Helper()
	t.Run("Rust test with all ACCEPTED", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("testdata/rust/main.rs")
		testIn := []string{""}
		testOut := []string{"x x y: 12"}
		var testCases int = 1
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "rust", "main.rs", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		log.Print(got, "\n", got.Status, got.Message, got.Error, got.Result)
		wantStatus := "OK"
		wantMessage := ""
		wantError := []string{""}
		wantResult := []string{"ACCEPTED"}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})

	t.Run("Rust test with 1 WRONG ANSWER", func(t *testing.T) {
		testDat, _ := ioutil.ReadFile("testdata/rust/main.rs")
		testIn := []string{""}
		testOut := []string{"-1"}
		var testCases int = 1
		var testTimeLimit int = 1
		var testMemoryLimit int64 = 500 * 1024 * 1024
		got := run.Evaluate(string(testDat), "rust", "main.rs", testIn, testOut, testCases, testTimeLimit, testMemoryLimit)
		log.Print(got, "\n", got.Status, got.Message, got.Error, got.Result)
		wantStatus := "OK"
		wantMessage := ""
		wantError := []string{"", "", ""}
		wantResult := []string{"WRONG ANSWER"}
		utils.CompareUtils(got, wantStatus, wantMessage, wantError, wantResult, t)
	})
}
