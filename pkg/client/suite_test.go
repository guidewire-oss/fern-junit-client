package client

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/junit"
)

const (
	testProjectName     = "TestProject"
	nonExistentFilePath = "this_file_does_not_exist"

	reportsCombinedPattern = "../../test/static/*.xml"
	reportFailedPath       = "../../test/static/junit_report_failed.xml"
	reportPassedPath       = "../../test/static/junit_report_passed.xml"

	fernTestRunCombinedPath = "../../test/static/fern_test_run_combined.json"
	fernTestRunFailedPath   = "../../test/static/fern_test_run_failed.json"
	fernTestRunPassedPath   = "../../test/static/fern_test_run_passed.json"

	exampleTags = "test,tagtest,9=-+_"
)

var (
	mockFernReporter *httptest.Server

	junitTestSuiteFailed junit.TestSuite
	junitTestSuitePassed junit.TestSuite

	fernTestRunCombined fern.TestRun
	fernTestRunFailed   fern.TestRun
	fernTestRunPassed   fern.TestRun
)

func init() {
	// Create mock Fern reporter
	mockFernReporter = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if data, err := io.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if err := json.Unmarshal(data, &fern.TestRun{}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))

	// Unmarshal JUnit test data
	junitTestSuiteFailed = junit.TestSuite{}
	if bytes, err := os.ReadFile(reportFailedPath); err != nil {
		panic(err)
	} else if err := xml.Unmarshal(bytes, &junitTestSuiteFailed); err != nil {
		panic(err)
	}

	junitTestSuitePassed = junit.TestSuite{}
	if bytes, err := os.ReadFile(reportPassedPath); err != nil {
		panic(err)
	} else if err := xml.Unmarshal(bytes, &junitTestSuitePassed); err != nil {
		panic(err)
	}

	// Unmarshal Fern test data
	fernTestRunCombined = fern.TestRun{}
	if bytes, err := os.ReadFile(fernTestRunCombinedPath); err != nil {
		panic(err)
	} else if err := json.Unmarshal(bytes, &fernTestRunCombined); err != nil {
		panic(err)
	}

	fernTestRunFailed = fern.TestRun{}
	if bytes, err := os.ReadFile(fernTestRunFailedPath); err != nil {
		panic(err)
	} else if err := json.Unmarshal(bytes, &fernTestRunFailed); err != nil {
		panic(err)
	}

	fernTestRunPassed = fern.TestRun{}
	if bytes, err := os.ReadFile(fernTestRunPassedPath); err != nil {
		panic(err)
	} else if err := json.Unmarshal(bytes, &fernTestRunPassed); err != nil {
		panic(err)
	}
}
