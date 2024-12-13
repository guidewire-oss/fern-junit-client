package client

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/junit"
)

const (
	testProjectName = "TestProject"
	exampleTags     = "test,tagtest,9=-+_"

	nonExistentFilePath = "this_file_does_not_exist"

	reportsCombinedPattern     = "../../test/static/*.xml"
	reportFailedPath           = "../../test/static/junit_report_failed.xml"
	reportPassedPath           = "../../test/static/junit_report_passed.xml"
	reportNoOptionalFieldsPath = "../../test/static/junit_report_no_optional_passed.xml"

	fernTestRunCombinedPath = "../../test/static/fern_test_run_combined.json"
	fernTestRunFailedPath   = "../../test/static/fern_test_run_failed.json"
	fernTestRunPassedPath   = "../../test/static/fern_test_run_passed.json"
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

func TestGenerateStaticJsonFiles(t *testing.T) {
	// Skip on normal test execution
	if os.Getenv("GENERATE_STATIC_FILES") != "true" {
		t.SkipNow()
	}
	// Generate static JSON files used in tests
	generateFernTestRunJson := func(testCase, filePattern, outputFilename string) {
		// Create mock Fern reporter that writes JSON payload to a file
		mockFernReporter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			testRun := &fern.TestRun{}
			if data, err := io.ReadAll(r.Body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else if err := json.Unmarshal(data, testRun); err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
				// Format JSON
				formattedJson, err := json.MarshalIndent(testRun, "", "  ")
				if err != nil {
					panic(err)
				}
				// Save JSON
				f, err := os.Create(outputFilename)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				_, err = f.Write(formattedJson)
				if err != nil {
					panic(err)
				}
			}
		}))
		defer mockFernReporter.Close()
		// Run SendReports with input
		if err := SendReports(mockFernReporter.URL, testProjectName, filePattern, exampleTags, true); err != nil {
			panic(fmt.Errorf("failed to generate Fern test run JSON for '%s' test case (%s): %w", testCase, filePattern, err))
		}
	}
	// Call generateFernTestRunJson for each test case
	generateFernTestRunJson("reports combined", reportsCombinedPattern, fernTestRunCombinedPath)
	generateFernTestRunJson("failed report", reportFailedPath, fernTestRunFailedPath)
	generateFernTestRunJson("passed report", reportPassedPath, fernTestRunPassedPath)
}
