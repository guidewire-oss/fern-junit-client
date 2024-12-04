package client

import (
	"log"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

func SendReports(fernUrl, projectName, filePattern string, tags string, verbose bool) error {
	var testRun fern.TestRun
	testRun.TestProjectName = projectName

	log.Default().Println("Parsing reports...")
	if err := parseReports(&testRun, filePattern, tags, verbose); err != nil {
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
