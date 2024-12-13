package client

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/util"
)

func SendReports(fernUrl, projectName, filePattern string, tags string, verbose bool, metricsFilePath string) error {
	var testRun fern.TestRun
	testRun.TestProjectName = projectName
	testRun.TestSeed = uint64(util.GlobalClock.Now().Nanosecond())

	log.Default().Println("Parsing reports...")
	if err := parseReports(&testRun, filePattern, tags, verbose); err != nil {
		return err
	}
	log.Default().Println("Parsing reports succeeded!")

	// Run the html request and test metric logging separately
	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		log.Default().Println("Sending reports to Fern...")
		if err := sendTestRun(testRun, fernUrl, verbose); err != nil {
			errChan <- err
		}
		log.Default().Println("Sending reports succeeded!")
	}()

	go func() {
		defer wg.Done()
		if err := recordTestRunMetrics(testRun, metricsFilePath); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}

	return nil
}

func recordTestRunMetrics(testRun fern.TestRun, metricsFilePath string) error {
	passed, failed, skipped := 0, 0, 0
	for _, suiteRun := range testRun.SuiteRuns {
		for _, specRun := range suiteRun.SpecRuns {
			switch specRun.Status {
			case "passed":
				passed++
			case "failed":
				failed++
			case "skipped":
				skipped++
			}
		}
	}

	log.Default().Printf("Total tests passed: %d\n", passed)
	log.Default().Printf("Total tests failed: %d\n", failed)
	log.Default().Printf("Total tests skipped: %d\n", skipped)

	if file, err := os.Create(metricsFilePath); err != nil {
		return err
	} else {
		file.WriteString(fmt.Sprintf("passed: %d\nfailed: %d\nskipped: %d\n", passed, failed, skipped))
		file.Close()
	}

	return nil
}
