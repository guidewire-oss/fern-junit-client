package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/junit"
)

func parseReports(testRun *fern.TestRun, filePattern string, verbose bool) error {
	files, err := filepath.Glob(filePattern)
	if err != nil {
		return fmt.Errorf("failed to parse file pattern %s: %w", filePattern, err)
	}
	for _, file := range files {
		suiteRun, err := parseReport(file, verbose)
		if err != nil {
			return fmt.Errorf("failed to parse report %s: %w", file, err)
		}
		testRun.SuiteRuns = append(testRun.SuiteRuns, suiteRun...)
	}
	return nil
}

func parseReport(filePath string, verbose bool) ([]fern.SuiteRun, error) {
	var testSuites junit.TestSuites
	var testSuite junit.TestSuite
	var suiteRuns []fern.SuiteRun

	if verbose {
		log.Default().Printf("Reading %s\n", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if verbose {
		log.Default().Printf("Unmarshaling %s\n", filePath)
	}

	if err := xml.Unmarshal(byteValue, &testSuites); err != nil {
		if err = xml.Unmarshal(byteValue, &testSuite); err != nil {
			return nil, err
		} else {
			testSuites.TestSuites = append(testSuites.TestSuites, testSuite)
		}
	}

	for _, suite := range testSuites.TestSuites {
		run, err := parseTestSuite(suite)
		if err != nil {
			return nil, err
		}
		suiteRuns = append(suiteRuns, run)
	}

	return suiteRuns, err
}

func parseTestSuite(testSuite junit.TestSuite) (suiteRun fern.SuiteRun, err error) {
	suiteRun.SuiteName = testSuite.Name
	suiteRun.StartTime, err = time.Parse(time.RFC3339, testSuite.Timestamp)
	suiteRun.EndTime, err = getEndTime(suiteRun.StartTime, testSuite.Time)
	if err != nil {
		err = fmt.Errorf("failed to parse TestSuite %s: %w", testSuite.Name, err)
		return
	}

	runningTime := suiteRun.StartTime
	for _, testcase := range testSuite.TestCases {
		status := ""
		message := ""
		if len(testcase.Failures) > 0 {
			status = "failed"
			message = testcase.Failures[0].Message + "\n" + testcase.Failures[0].Content
		} else if len(testcase.Errors) > 0 {
			status = "failed"
			message = testcase.Errors[0].Message + "\n" + testcase.Errors[0].Content
		} else if len(testcase.Skips) > 0 {
			status = "skipped"
		} else {
			status = "passed"
		}

		after, erri := getEndTime(runningTime, testSuite.Time)
		if erri != nil {
			err = fmt.Errorf("failed to parse TestSuite %s: %w", testSuite.Name, erri)
			return
		}

		sr := fern.SpecRun{
			SpecDescription: testcase.Name,
			Status:          status,
			Message:         message,
			Tags:            []fern.Tag{},
			StartTime:       runningTime,
			EndTime:         after,
		}
		runningTime = after

		suiteRun.SpecRuns = append(suiteRun.SpecRuns, sr)
	}

	return
}

func getEndTime(startTime time.Time, duration string) (endTime time.Time, err error) {
	ms, err := time.ParseDuration(duration + "s")
	endTime = startTime.Add(ms)
	return
}
