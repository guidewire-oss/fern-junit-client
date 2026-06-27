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
		if err := client.SendReports(fernUrl, projectId, filePattern, tags, branch, commitSha, environment, verbose); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	sendCmd.PersistentFlags().StringVarP(&fernUrl, "fern-url", "u", "", "base URL of the Fern Platform instance to send test reports to (required)")
	sendCmd.PersistentFlags().StringVarP(&projectId, "project-id", "p", "", "Id of the project to associate test reports with (required). You must register the application first in Fern Platform")
	sendCmd.PersistentFlags().StringVarP(&filePattern, "file-pattern", "f", "", "file name pattern of test reports to send to Fern (required)")
	sendCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "comma-separated tags to be included on runs")
	sendCmd.PersistentFlags().StringVar(&branch, "branch", "", "git branch name for this run (falls back to $GITHUB_REF_NAME / $CI_COMMIT_REF_NAME)")
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
