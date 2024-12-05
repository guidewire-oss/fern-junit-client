package client

import (
	"log"
	"time"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

func SendReports(fernUrl, projectName, filePattern string, verbose bool) error {
	var testRun fern.TestRun
	testRun.TestProjectName = projectName
	testRun.TestSeed = uint64(time.Now().Nanosecond())

	log.Default().Println("Parsing reports...")
	if err := parseReports(&testRun, filePattern, verbose); err != nil {
		return err
	}
	log.Default().Println("Parsing reports succeeded!")

	log.Default().Println("Sending reports to Fern...")
	if err := sendTestRun(testRun, fernUrl, verbose); err != nil {
		return err
	}
	log.Default().Println("Sending reports succeeded!")
	return nil
}
