package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/guidewire-oss/fern-junit-client/pkg/client"
)

var (
	fernUrl     string
	projectName string
	filePattern string
	tags        string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send JUnit test reports to Fern",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.SendReports(fernUrl, projectName, filePattern, tags, verbose); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	sendCmd.PersistentFlags().StringVarP(&fernUrl, "fern-url", "u", "", "base URL of the Fern Reporter instance to send test reports to (required)")
	sendCmd.PersistentFlags().StringVarP(&projectName, "project-name", "p", "", "name of the project to associate test reports with (required)")
	sendCmd.PersistentFlags().StringVarP(&filePattern, "file-pattern", "f", "", "file name pattern of test reports to send to Fern (required)")
	sendCmd.PersistentFlags().StringVarP(&tags, "tags", "t", "", "comma-separated tags to be included on runs")
	err := sendCmd.MarkPersistentFlagRequired("fern-url")
	err = sendCmd.MarkPersistentFlagRequired("project-name")
	err = sendCmd.MarkPersistentFlagRequired("file-pattern")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
	rootCmd.AddCommand(sendCmd)
}
