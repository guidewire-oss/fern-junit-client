package client

import (
	"encoding/xml"
	"fmt"
	"github.com/guidewire-oss/fern-junit-client/pkg/util"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/junit"
)

func parseReports(testRun *fern.TestRun, filePattern string, tags string, verbose bool) error {
	files, err := filepath.Glob(filePattern)
	if err != nil {
		return fmt.Errorf("failed to parse file pattern %s: %w", filePattern, err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found for pattern %s", filePattern)
	}
	for _, file := range files {
		suiteRun, err := parseReport(file, tags, verbose)
		if err != nil {
			return fmt.Errorf("failed to parse report %s: %w", file, err)
		}
		testRun.SuiteRuns = append(testRun.SuiteRuns, suiteRun...)
	}
	for _, suiteRun := range testRun.SuiteRuns {
		// Set testRun.StartTime to the earliest suite start time
		if testRun.StartTime.IsZero() || suiteRun.StartTime.Compare(testRun.StartTime) < 0 {
			testRun.StartTime = suiteRun.StartTime
		}
		// Set testRun.EndTime to the latest suite end time
		if testRun.EndTime.IsZero() || suiteRun.EndTime.Compare(testRun.EndTime) > 0 {
			testRun.EndTime = suiteRun.EndTime
		}
	}
	if verbose {
		log.Default().Printf("TestRun start time: %s\n", testRun.StartTime.String())
		log.Default().Printf("TestRun end time: %s\n", testRun.EndTime.String())
	}
	return nil
}

func parseReport(filePath string, tags string, verbose bool) ([]fern.SuiteRun, error) {
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
		run, err := parseTestSuite(util.RealClock{}, suite, tags, verbose)
		if err != nil {
			return nil, err
		}
		suiteRuns = append(suiteRuns, run)
	}
	return suiteRuns, err
}

func parseTestSuite(clock util.Clock, testSuite junit.TestSuite, tags string, verbose bool) (suiteRun fern.SuiteRun, err error) {
	if verbose {
		log.Default().Printf("Parsing TestSuite %s\n", testSuite.Name)
	}

	suiteRun.SuiteName = testSuite.Name

	if testSuite.Timestamp == "" {
		suiteRun.StartTime = clock.Now()
	} else {
		suiteRun.StartTime, err = time.Parse(time.RFC3339, testSuite.Timestamp)
		if err != nil {
			// Attempt to parse with a "Z" suffix for UTC time if the initial parsing fails
			suiteRun.StartTime, err = time.Parse(time.RFC3339, testSuite.Timestamp+"Z")
			if err != nil {
				err = fmt.Errorf("failed to parse suite start time: %w", err)
				return
			}
		}
	}

	suiteRun.EndTime, err = getEndTime(suiteRun.StartTime, testSuite.Time)
	if err != nil {
		err = fmt.Errorf("failed to calculate suite end time: %w", err)
		return
	}

	if verbose {
		log.Default().Printf("Suite start time: %s\n", suiteRun.StartTime.String())
		log.Default().Printf("Suite end time: %s\n", suiteRun.EndTime.String())
	}

	startTime := suiteRun.StartTime
	var endTime time.Time
	for _, testCase := range testSuite.TestCases {
		status := ""
		message := ""
		if len(testCase.Failures) > 0 {
			status = "failed"
			message = testCase.Failures[0].Message + "\n" + testCase.Failures[0].Content
		} else if len(testCase.Errors) > 0 {
			status = "failed"
			message = testCase.Errors[0].Message + "\n" + testCase.Errors[0].Content
		} else if len(testCase.Skips) > 0 {
			status = "skipped"
		} else {
			status = "passed"
		}

		endTime, err = getEndTime(startTime, testCase.Time)
		if err != nil {
			err = fmt.Errorf("failed to calculate test end time: %w", err)
			return
		}

		specRun := fern.SpecRun{
			SpecDescription: testCase.Name,
			Status:          status,
			Message:         message,
			Tags:            convertToTags(tags),
			StartTime:       startTime,
			EndTime:         endTime,
		}
		suiteRun.SpecRuns = append(suiteRun.SpecRuns, specRun)

		startTime = endTime

		if verbose {
			log.Default().Printf("%#v\n", specRun)
		}
	}
	return
}

func getEndTime(startTime time.Time, durationSeconds string) (endTime time.Time, err error) {
	ms, err := time.ParseDuration(durationSeconds + "s")
	endTime = startTime.Add(ms)
	return
}

func convertToTags(tagString string) (tags []fern.Tag) {
	for _, tag := range strings.Split(tagString, ",") {
		tags = append(tags, fern.Tag{Name: tag})
	}
	return
}
