package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/codeiiest/test-runner/runner/tester"
)

// GetFilesInPwd lists files in the pwd
func GetFilesInPwd() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

// CompareUtils is Helper function to assert corrected of expected and received results
func CompareUtils(got tester.TestResult, wantStatus string, wantMessage string, wantError []string, wantResult []string, t *testing.T) {
	if got.Status != wantStatus {
		t.Errorf("got %s want %s", got.Status, wantStatus)
	}
	// Checking whether want is empty or not, if it is, it means either the value is irrelevant
	// OR it is a specific test where comparing that value is not feasible in realtime (say, long compilation errors)
	if wantMessage != "" {
		if got.Message != wantMessage {
			t.Errorf("got %s want %s", got.Message, wantMessage)
		}
	}
	for index := range got.Error {
		if len(wantError) > 0 {
			if (got.Error)[index] != wantError[index] {
				t.Errorf("got %s want %s", (got.Error)[index], wantError[index])
			}
		}
		if len(wantResult) > 0 {
			if (got.Result)[index] != wantResult[index] {
				t.Errorf("got %s want %s", (got.Result)[index], wantResult[index])
			}
		}
	}
}
