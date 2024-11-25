package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

func sendTestRun(testRun fern.TestRun, fernUrl string, verbose bool) error {
	payload, err := json.Marshal(testRun)
	if err != nil {
		return fmt.Errorf("failed to serialize TestRun: %w", err)
	}

	endpoint, err := url.JoinPath(fernUrl, "api", "testrun")
	if err != nil {
		return fmt.Errorf("failed to construct endpoint URL: %w", err)
	}

	if verbose {
		log.Default().Printf("Sending POST request to %s...\n", endpoint)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
	return nil
}
