package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/guidewire-oss/fern-junit-client/pkg/auth"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

// Test constants
const (
	testClientID        = "test-client-id"
	testClientSecret    = "test-client-secret"
	testProjectID       = "test-project-123"
	testAccessToken     = "test-access-token-123"
	testScopes          = "read write admin"
	integrationClientID = "integration-client"
	integrationSecret   = "integration-secret"
	testAPIPath         = "api/v1/test-runs"
	customAPIPath       = "custom/api/path"
)

// Helper function to save and restore OAuth environment variables
func withCleanOAuthEnv(t *testing.T, fn func()) {
	t.Helper()
	origTokenURL := os.Getenv("AUTH_URL")
	origClientID := os.Getenv("FERN_AUTH_CLIENT_ID")
	origClientPassword := os.Getenv("FERN_AUTH_CLIENT_SECRET")
	origScopes := os.Getenv("FERN_CLIENT_SCOPE")
	origAPIPath := os.Getenv("FERN_API_ENDPOINT_PATH")

	defer func() {
		os.Setenv("AUTH_URL", origTokenURL)
		os.Setenv("FERN_AUTH_CLIENT_ID", origClientID)
		os.Setenv("FERN_AUTH_CLIENT_SECRET", origClientPassword)
		os.Setenv("FERN_CLIENT_SCOPE", origScopes)
		os.Setenv("FERN_API_ENDPOINT_PATH", origAPIPath)
	}()

	fn()
}

// Helper function to set OAuth environment variables
func setOAuthTestEnv(tokenURL, clientID, clientSecret, scopes string) {
	os.Setenv("AUTH_URL", tokenURL)
	os.Setenv("FERN_AUTH_CLIENT_ID", clientID)
	os.Setenv("FERN_AUTH_CLIENT_SECRET", clientSecret)
	if scopes != "" {
		os.Setenv("FERN_CLIENT_SCOPE", scopes)
	}
}

// Helper function to clear OAuth environment variables
func clearOAuthTestEnv() {
	os.Unsetenv("AUTH_URL")
	os.Unsetenv("FERN_AUTH_CLIENT_ID")
	os.Unsetenv("FERN_AUTH_CLIENT_SECRET")
	os.Unsetenv("FERN_CLIENT_SCOPE")
}

// Helper function to create standard OAuth response
func createTokenResponse(token string, scopes string) auth.TokenResponse {
	return auth.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       scopes,
	}
}


// Helper function to verify Fern server request
func verifyFernRequest(t *testing.T, r *http.Request, expectedAuth string, expectedPath string) {
	t.Helper()

	// Verify authorization header
	if expectedAuth != "" && r.Header.Get("Authorization") != expectedAuth {
		t.Errorf("Expected Authorization: %s, got %s", expectedAuth, r.Header.Get("Authorization"))
	} else if expectedAuth == "" && r.Header.Get("Authorization") != "" {
		t.Errorf("Should not have Authorization header, got %s", r.Header.Get("Authorization"))
	}

	// Verify content type
	if r.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
	}

	// Verify method
	if r.Method != "POST" {
		t.Errorf("Expected POST request, got %s", r.Method)
	}

	// Verify path
	if expectedPath != "" && !strings.Contains(r.URL.Path, expectedPath) {
		t.Errorf("Expected path to contain %s, got %s", expectedPath, r.URL.Path)
	}
}


func Test_sendTestRun_WithAuthentication(t *testing.T) {
	withCleanOAuthEnv(t, func() {
		tests := []struct {
			name           string
			setupEnv       func(oauthURL string)
			oauthStatus    int
			oauthToken     string
			fernStatus     int
			expectedAuth   string
			expectedPath   string
			expectError    bool
			errorContains  string
		}{
			{
				name: "successful authentication and send",
				setupEnv: func(oauthURL string) {
					setOAuthTestEnv(oauthURL, testClientID, testClientSecret, "")
				},
				oauthStatus:  http.StatusOK,
				oauthToken:   testAccessToken,
				fernStatus:   http.StatusOK,
				expectedAuth: "Bearer " + testAccessToken,
				expectedPath: testAPIPath,
				expectError:  false,
			},
			{
				name: "successful authentication with scopes",
				setupEnv: func(oauthURL string) {
					setOAuthTestEnv(oauthURL, testClientID, testClientSecret, testScopes)
				},
				oauthStatus:  http.StatusOK,
				oauthToken:   "scoped-token",
				fernStatus:   http.StatusOK,
				expectedAuth: "Bearer scoped-token",
				expectedPath: testAPIPath,
				expectError:  false,
			},
			{
				name: "authentication disabled - no oauth env vars",
				setupEnv: func(oauthURL string) {
					clearOAuthTestEnv()
				},
				fernStatus:   http.StatusOK,
				expectedAuth: "",
				expectedPath: testAPIPath,
				expectError:  false,
			},
			{
				name: "oauth token fetch fails",
				setupEnv: func(oauthURL string) {
					setOAuthTestEnv(oauthURL, testClientID, "wrong-password", "")
				},
				oauthStatus:   http.StatusUnauthorized,
				expectError:   true,
				errorContains: "failed to add OAuth authentication",
			},
			{
				name: "custom API endpoint path",
				setupEnv: func(oauthURL string) {
					setOAuthTestEnv(oauthURL, testClientID, testClientSecret, "")
					os.Setenv("FERN_API_ENDPOINT_PATH", customAPIPath)
				},
				oauthStatus:  http.StatusOK,
				oauthToken:   "custom-endpoint-token",
				fernStatus:   http.StatusOK,
				expectedAuth: "Bearer custom-endpoint-token",
				expectedPath: customAPIPath,
				expectError:  false,
			},
			{
				name: "fern server returns 401 with valid oauth",
				setupEnv: func(oauthURL string) {
					setOAuthTestEnv(oauthURL, testClientID, testClientSecret, "")
				},
				oauthStatus:   http.StatusOK,
				oauthToken:    "valid-token",
				fernStatus:    http.StatusUnauthorized,
				expectedAuth:  "Bearer valid-token",
				expectError:   true,
				errorContains: "unexpected response code: 401",
			},
			{
				name: "partial oauth config - missing client id",
				setupEnv: func(oauthURL string) {
					os.Setenv("AUTH_URL", oauthURL)
					os.Unsetenv("FERN_AUTH_CLIENT_ID")
					os.Setenv("FERN_AUTH_CLIENT_SECRET", testClientSecret)
				},
				fernStatus:   http.StatusOK,
				expectedAuth: "",
				expectedPath: testAPIPath,
				expectError:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create OAuth server
				var oauthCalled bool
				mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					oauthCalled = true

					if tt.oauthStatus != http.StatusOK {
						w.WriteHeader(tt.oauthStatus)
						w.Write([]byte(`{"error": "invalid_client"}`))
						return
					}

					// For scopes test, verify scopes are sent
					if tt.name == "successful authentication with scopes" {
						r.ParseForm()
						if r.FormValue("scope") != testScopes {
							t.Errorf("Expected scope='%s', got %s", testScopes, r.FormValue("scope"))
						}
					}

					response := createTokenResponse(tt.oauthToken, "")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
				defer mockOAuthServer.Close()

				// Create Fern server
				var fernCalled bool
				mockFernServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fernCalled = true

					// Read and verify body
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Failed to read request body: %v", err)
					}
					defer r.Body.Close()

					var testRun fern.TestRun
					if err := json.Unmarshal(body, &testRun); err != nil {
						t.Errorf("Failed to unmarshal test run: %v", err)
					}

					// Verify request
					verifyFernRequest(t, r, tt.expectedAuth, tt.expectedPath)

					// Send response
					w.WriteHeader(tt.fernStatus)
					if tt.fernStatus == http.StatusOK {
						w.Write([]byte(`{"success": true}`))
					} else {
						w.Write([]byte(`{"error": "unauthorized"}`))
					}
				}))
				defer mockFernServer.Close()

				// Setup environment
				tt.setupEnv(mockOAuthServer.URL)

				// Create and send test run
				testRun := fern.TestRun{
					TestProjectID: testProjectID,
					TestSeed:      12345,
				}

				err := sendTestRun(testRun, mockFernServer.URL, true)

				// Check error expectations
				if tt.expectError && err == nil {
					t.Errorf("Expected error but got nil")
				} else if !tt.expectError && err != nil {
					t.Errorf("Unexpected error = %v", err)
				} else if tt.expectError && err != nil && tt.errorContains != "" {
					if !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("Error = %v, expected to contain %s", err, tt.errorContains)
					}
				}

				// Verify OAuth server was called when expected
				if tt.name == "authentication disabled - no oauth env vars" && oauthCalled {
					t.Error("OAuth server should not be called when disabled")
				}

				// Verify Fern was called when expected
				if !tt.expectError && !fernCalled && tt.name != "oauth token fetch fails" {
					t.Error("Fern server was not called when expected")
				}
			})
		}
	})
}

func TestSendReports_Authentication(t *testing.T) {
	withCleanOAuthEnv(t, func() {
		t.Run("with authentication", func(t *testing.T) {
			// Create OAuth server
			oauthCallCount := 0
			mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				oauthCallCount++
				response := createTokenResponse(fmt.Sprintf("integration-token-%d", oauthCallCount), "")
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}))
			defer mockOAuthServer.Close()

			// Create Fern server
			var receivedAuthHeaders []string
			mockFernServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				receivedAuthHeaders = append(receivedAuthHeaders, authHeader)

				// Verify body
				body, _ := io.ReadAll(r.Body)
				defer r.Body.Close()

				var testRun fern.TestRun
				if err := json.Unmarshal(body, &testRun); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer mockFernServer.Close()

			// Setup OAuth
			setOAuthTestEnv(mockOAuthServer.URL, integrationClientID, integrationSecret, "")

			// Run SendReports
			err := SendReports(mockFernServer.URL, testProjectId, reportPassedPath, "auth-test", true)
			if err != nil {
				t.Fatalf("SendReports() with authentication failed: %v", err)
			}

			// Verify OAuth was called
			if oauthCallCount < 1 {
				t.Errorf("Expected at least 1 OAuth call, got %d", oauthCallCount)
			}

			// Verify auth header
			if len(receivedAuthHeaders) < 1 {
				t.Fatal("No auth headers received by Fern server")
			}
			expectedAuthHeader := "Bearer integration-token-1"
			if receivedAuthHeaders[0] != expectedAuthHeader {
				t.Errorf("Expected auth header %s, got %s", expectedAuthHeader, receivedAuthHeaders[0])
			}
		})

		t.Run("without authentication", func(t *testing.T) {
			// Clear OAuth env vars
			clearOAuthTestEnv()

			// Create Fern server that verifies no auth
			mockFernServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify no auth header
				if r.Header.Get("Authorization") != "" {
					t.Errorf("Should not have Authorization header when OAuth disabled, got %s", r.Header.Get("Authorization"))
				}

				// Parse body
				body, _ := io.ReadAll(r.Body)
				defer r.Body.Close()

				var testRun fern.TestRun
				if err := json.Unmarshal(body, &testRun); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer mockFernServer.Close()

			// Run SendReports
			err := SendReports(mockFernServer.URL, testProjectId, reportPassedPath, "no-auth-test", false)
			if err != nil {
				t.Fatalf("SendReports() without authentication failed: %v", err)
			}
		})
	})
}