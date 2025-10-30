package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/guidewire-oss/fern-junit-client/pkg/auth"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

func sendTestRun(testRun fern.TestRun, fernUrl string, verbose bool) error {
	payload, err := json.Marshal(testRun)
	if err != nil {
		return fmt.Errorf("failed to serialize TestRun: %w", err)
	}

	// Log the request payload in verbose mode
	if verbose {
		// Pretty print the JSON for better readability
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, payload, "", "  "); err == nil {
			log.Default().Printf("Request Payload:\n%s\n", prettyJSON.String())
		} else {
			// If pretty print fails, just print raw
			log.Default().Printf("Request Payload: %s\n", string(payload))
		}
	}

	// Check for custom API endpoint path from environment variable
	// Default to "api/v1/test-runs" if not set
	apiPath := os.Getenv("FERN_API_ENDPOINT_PATH")
	if apiPath == "" {
		apiPath = "api/v1/test-runs"
	}

	endpoint, err := url.JoinPath(fernUrl, apiPath)
	if err != nil {
		return fmt.Errorf("failed to construct endpoint URL: %w", err)
	}

	if verbose {
		log.Default().Printf("Using API endpoint path: %s\n", apiPath)
		log.Default().Printf("Full endpoint URL: %s\n", endpoint)
	}

	// Create OAuth client
	oauthClient := auth.NewOAuthClient()

	// Create HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add OAuth authentication if enabled
	if oauthClient != nil && oauthClient.IsEnabled() {
		if verbose {
			log.Default().Println("OAuth authentication is enabled, fetching token...")
		}
		if err := oauthClient.AddAuthHeader(req); err != nil {
			return fmt.Errorf("failed to add OAuth authentication: %w", err)
		}
	} else if verbose {
		log.Default().Println("OAuth authentication is not configured, proceeding without authentication")
	}

	if verbose {
		log.Default().Printf("Sending POST request to %s...\n", endpoint)
	}

	// Use the OAuth client's HTTP client if OAuth is enabled, otherwise use default
	var httpClient *http.Client
	if oauthClient != nil && oauthClient.IsEnabled() {
		httpClient = oauthClient.HTTPClient()
	} else {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Log response details when verbose is enabled
	if verbose {
		log.Default().Printf("Response Status: %s\n", resp.Status)
		log.Default().Printf("Response Status Code: %d\n", resp.StatusCode)

		// Read the entire response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Default().Printf("Error reading response body: %v\n", err)
		} else {
			// Replace the response body so it can still be read later if needed
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			if len(bodyBytes) > 0 {
				// Try to pretty print JSON if possible
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
					log.Default().Printf("Response Body:\n%s\n", prettyJSON.String())
				} else {
					// If not JSON, print as is
					log.Default().Printf("Response Body:\n%s\n", string(bodyBytes))
				}
			} else {
				log.Default().Println("Response Body: (empty)")
			}
		}

		// Log response headers
		log.Default().Println("Response Headers:")
		for key, values := range resp.Header {
			for _, value := range values {
				log.Default().Printf("  %s: %s\n", key, value)
			}
		}
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
	return nil
}
