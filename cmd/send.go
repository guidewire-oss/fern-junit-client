package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/guidewire-oss/fern-junit-client/pkg/client"
)

var (
	fernUrl     string
	projectId   string
	filePattern string
	tags        string
	branch      string
	commitSha   string
	environment string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send JUnit test reports to Fern",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		b, c, e := resolveProvenance(branch, commitSha, environment)
		if err := client.SendReports(client.SendOptions{
			FernURL: fernUrl, ProjectID: projectId, FilePattern: filePattern, Tags: tags,
			Branch: b, CommitSha: c, Environment: e, Verbose: verbose,
		}); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	},
}

// resolveProvenance fills in git provenance from common CI env vars when the
// corresponding flag was not provided (an explicit flag value always wins).
func resolveProvenance(branch, commit, environment string) (string, string, string) {
	// GITHUB_HEAD_REF is the PR source branch (set only on pull_request events) and
	// takes precedence over GITHUB_REF_NAME, which on PR builds is the synthetic
	// refs/pull/<n>/merge ref rather than a real branch name.
	return firstNonEmpty(branch, os.Getenv("GITHUB_HEAD_REF"), os.Getenv("GITHUB_REF_NAME"), os.Getenv("CI_COMMIT_REF_NAME")),
		firstNonEmpty(commit, os.Getenv("GITHUB_SHA"), os.Getenv("CI_COMMIT_SHA")),
		firstNonEmpty(environment, os.Getenv("CI_ENVIRONMENT_NAME"), os.Getenv("FERN_ENVIRONMENT"))
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

func init() {
	sendCmd.PersistentFlags().StringVarP(&fernUrl, "fern-url", "u", "", "base URL of the Fern Platform instance to send test reports to (required)")
	sendCmd.PersistentFlags().StringVarP(&projectId, "project-id", "p", "", "Id of the project to associate test reports with (required). You must register the application first in Fern Platform")
	sendCmd.PersistentFlags().StringVarP(&filePattern, "file-pattern", "f", "", "file name pattern of test reports to send to Fern (required)")
	sendCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "comma-separated tags to be included on runs")
	sendCmd.PersistentFlags().StringVar(&branch, "branch", "", "git branch name for this run (falls back to $GITHUB_HEAD_REF / $GITHUB_REF_NAME / $CI_COMMIT_REF_NAME)")
	sendCmd.PersistentFlags().StringVar(&commitSha, "commit", "", "git commit SHA for this run (falls back to $GITHUB_SHA / $CI_COMMIT_SHA)")
	sendCmd.PersistentFlags().StringVar(&environment, "environment", "", "environment label, e.g. ci, staging (falls back to $CI_ENVIRONMENT_NAME / $FERN_ENVIRONMENT)")
	if err := sendCmd.MarkPersistentFlagRequired("fern-url"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if err := sendCmd.MarkPersistentFlagRequired("project-id"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if err := sendCmd.MarkPersistentFlagRequired("file-pattern"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	rootCmd.AddCommand(sendCmd)
}
