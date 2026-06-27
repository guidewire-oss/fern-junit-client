package client

import (
	"log"
	"os"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/util"
)

// SendReports parses JUnit XML reports and posts a TestRun to Fern Platform.
// branch, commitSha, and environment may be empty strings; when empty the
// function falls back to GITHUB_REF_NAME / GITHUB_SHA / CI_ENVIRONMENT_NAME
// environment variables so that callers in CI pipelines get automatic population
// without explicit flags.
func SendReports(fernUrl, projectId, filePattern, tags, branch, commitSha, environment string, verbose bool) error {
	var testRun fern.TestRun
	testRun.TestProjectID = projectId
	testRun.TestSeed = uint64(util.GlobalClock.Now().Nanosecond())

	// Resolve git provenance: flag value wins; fall back to common CI env vars.
	if branch == "" {
		branch = firstNonEmpty(os.Getenv("GITHUB_REF_NAME"), os.Getenv("CI_COMMIT_REF_NAME"))
	}
	if commitSha == "" {
		commitSha = firstNonEmpty(os.Getenv("GITHUB_SHA"), os.Getenv("CI_COMMIT_SHA"))
	}
	if environment == "" {
		environment = firstNonEmpty(os.Getenv("CI_ENVIRONMENT_NAME"), os.Getenv("FERN_ENVIRONMENT"))
	}

	testRun.GitBranch = branch
	testRun.GitSha = commitSha
	testRun.Environment = environment

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

// firstNonEmpty returns the first non-empty string from the provided candidates.
func firstNonEmpty(candidates ...string) string {
	for _, s := range candidates {
		if s != "" {
			return s
		}
	}
	return ""
}
