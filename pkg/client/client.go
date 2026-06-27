package client

import (
	"log"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/util"
)

// SendOptions carries the parameters for SendReports as named fields, so callers
// can't transpose the many same-typed string arguments by position. Branch,
// CommitSha, and Environment are recorded as run provenance; leave them empty to
// omit. Resolving them from CI env vars is the caller's concern (see cmd), which
// keeps this client free of CI-provider specifics.
type SendOptions struct {
	FernURL     string
	ProjectID   string
	FilePattern string
	Tags        string
	Branch      string
	CommitSha   string
	Environment string
	Verbose     bool
}

// SendReports parses JUnit XML reports and posts a TestRun to Fern Platform.
func SendReports(opts SendOptions) error {
	var testRun fern.TestRun
	testRun.TestProjectID = opts.ProjectID
	testRun.TestSeed = uint64(util.GlobalClock.Now().Nanosecond())
	testRun.GitBranch = opts.Branch
	testRun.GitSha = opts.CommitSha
	testRun.Environment = opts.Environment

	log.Default().Println("Parsing reports...")
	if err := parseReports(&testRun, opts.FilePattern, opts.Tags, opts.Verbose); err != nil {
		return err
	}
	log.Default().Println("Parsing reports succeeded!")

	log.Default().Println("Sending reports to Fern...")
	if err := sendTestRun(testRun, opts.FernURL, opts.Verbose); err != nil {
		return err
	}
	log.Default().Println("Sending reports succeeded!")
	return nil
}
