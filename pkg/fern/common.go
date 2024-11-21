package fern

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guidewire-oss/fern-junit-client/pkg/models"
)

func sendTestRun(testRun models.TestRun, serviceUrl string) error {
	payload, err := json.Marshal(testRun)
	if err != nil {
		return fmt.Errorf("Failed to serialize test run: %v", err)
	}

	resp, err := http.Post(serviceUrl+"/api/testrun", "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Unexpected response code: %d", resp.StatusCode)
	}
	return nil
}
