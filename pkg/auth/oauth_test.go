package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// Test constants
const (
	testClientID       = "test-client"
	testClientSecret   = "test-secret"
	testTokenURL       = "https://auth.example.com/token"
	testAccessToken    = "test-access-token"
	testBearerToken    = "test-bearer-token"
	testScopes         = "read write admin"
)

// Helper function to save and restore environment variables
func withCleanEnv(t *testing.T, fn func()) {
	t.Helper()
	origTokenURL := os.Getenv("AUTH_URL")
	origClientID := os.Getenv("FERN_AUTH_CLIENT_ID")
	origClientPassword := os.Getenv("FERN_AUTH_CLIENT_SECRET")
	origScopes := os.Getenv("FERN_CLIENT_SCOPE")

	defer func() {
		os.Setenv("AUTH_URL", origTokenURL)
		os.Setenv("FERN_AUTH_CLIENT_ID", origClientID)
		os.Setenv("FERN_AUTH_CLIENT_SECRET", origClientPassword)
		os.Setenv("FERN_CLIENT_SCOPE", origScopes)
	}()

	fn()
}

// Helper function to set OAuth environment variables
func setOAuthEnv(tokenURL, clientID, clientPassword, scopes string) {
	os.Setenv("AUTH_URL", tokenURL)
	os.Setenv("FERN_AUTH_CLIENT_ID", clientID)
	os.Setenv("FERN_AUTH_CLIENT_SECRET", clientPassword)
	os.Setenv("FERN_CLIENT_SCOPE", scopes)
}


// Helper function to create a mock OAuth server with standard response
func createMockOAuthServer(t *testing.T, token string, expiresIn int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TokenResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   expiresIn,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}

// Helper function to verify error expectations
func assertError(t *testing.T, err error, expectError bool, errorContains string) {
	t.Helper()
	if expectError && err == nil {
		t.Errorf("Expected error but got nil")
	} else if !expectError && err != nil {
		t.Errorf("Unexpected error = %v", err)
	} else if expectError && err != nil && errorContains != "" {
		if !strings.Contains(err.Error(), errorContains) {
			t.Errorf("Error = %v, expected to contain %s", err, errorContains)
		}
	}
}

func TestNewOAuthClient(t *testing.T) {
	withCleanEnv(t, func() {
		tests := []struct {
			name           string
			tokenURL       string
			clientID       string
			clientPassword string
			scopes         string
			expectNil      bool
			expectError    bool
			errorContains  string
		}{
			{
				name:           "all OAuth vars set with scopes",
				tokenURL:       testTokenURL,
				clientID:       testClientID,
				clientPassword: testClientSecret,
				scopes:         testScopes,
				expectNil:      false,
				expectError:    false,
			},
			{
				name:           "OAuth vars set without scopes",
				tokenURL:       testTokenURL,
				clientID:       testClientID,
				clientPassword: testClientSecret,
				scopes:         "",
				expectNil:      false,
				expectError:    false,
			},
			{
				name:           "missing token URL (OAuth disabled)",
				tokenURL:       "",
				clientID:       testClientID,
				clientPassword: testClientSecret,
				expectNil:      true,
				expectError:    false,
			},
			{
				name:           "missing client ID",
				tokenURL:       testTokenURL,
				clientID:       "",
				clientPassword: testClientSecret,
				expectNil:      true,
				expectError:    true,
				errorContains:  "FERN_AUTH_CLIENT_ID",
			},
			{
				name:           "missing client password",
				tokenURL:       testTokenURL,
				clientID:       testClientID,
				clientPassword: "",
				expectNil:      true,
				expectError:    true,
				errorContains:  "FERN_AUTH_CLIENT_SECRET",
			},
			{
				name:      "all vars empty",
				tokenURL:  "",
				clientID:  "",
				expectNil: true,
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				setOAuthEnv(tt.tokenURL, tt.clientID, tt.clientPassword, tt.scopes)
				client, err := NewOAuthClient()

				if tt.expectError {
					if err == nil {
						t.Errorf("Expected error but got nil")
					} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("Error = %v, should contain %v", err, tt.errorContains)
					}
				} else if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if tt.expectNil && client != nil {
					t.Errorf("Expected nil but got client")
				} else if !tt.expectNil && client == nil {
					t.Errorf("Expected client but got nil")
				}

				// Verify scopes are set correctly when client is created
				if !tt.expectNil && client != nil && client.config.Scopes != tt.scopes {
					t.Errorf("Scopes = %v, want %v", client.config.Scopes, tt.scopes)
				}
			})
		}
	})
}

func TestOAuthClient_GetToken(t *testing.T) {
	t.Run("basic token operations", func(t *testing.T) {
		tokenCallCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCallCount++

			// Verify request
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("Expected Content-Type: application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
			}

			// Parse and verify form
			r.ParseForm()
			if r.FormValue("grant_type") != "client_credentials" {
				t.Errorf("Expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
			}
			if r.FormValue("client_id") != testClientID {
				t.Errorf("Expected client_id=%s, got %s", testClientID, r.FormValue("client_id"))
			}
			if r.FormValue("client_secret") != testClientSecret {
				t.Errorf("Expected client_secret=%s, got %s", testClientSecret, r.FormValue("client_secret"))
			}

			// Return different tokens for each call
			response := TokenResponse{
				AccessToken: testAccessToken + "-" + string(rune(tokenCallCount+'0')),
				TokenType:   "Bearer",
				ExpiresIn:   3600,
				Scope:       "test-scope",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer mockServer.Close()

		client := &OAuthClient{
			config: &OAuthConfig{
				TokenURL:       mockServer.URL,
				ClientID:       testClientID,
				ClientPassword: testClientSecret,
			},
		}

		// Test initial token fetch
		token, err := client.GetToken()
		if err != nil {
			t.Fatalf("GetToken() error = %v", err)
		}
		if token != testAccessToken+"-1" {
			t.Errorf("Expected token %s-1, got %s", testAccessToken, token)
		}
		if tokenCallCount != 1 {
			t.Errorf("Expected 1 token call, got %d", tokenCallCount)
		}

		// Test cached token (should not make another call)
		token2, err := client.GetToken()
		if err != nil {
			t.Fatalf("GetToken() cached error = %v", err)
		}
		if token2 != testAccessToken+"-1" {
			t.Errorf("Expected cached token %s-1, got %s", testAccessToken, token2)
		}
		if tokenCallCount != 1 {
			t.Errorf("Expected still 1 token call for cached token, got %d", tokenCallCount)
		}

		// Test token expiry and refresh
		client.tokenExpiry = time.Now().Add(-1 * time.Hour)
		token3, err := client.GetToken()
		if err != nil {
			t.Fatalf("GetToken() refresh error = %v", err)
		}
		if token3 != testAccessToken+"-2" {
			t.Errorf("Expected new token %s-2, got %s", testAccessToken, token3)
		}
		if tokenCallCount != 2 {
			t.Errorf("Expected 2 token calls after expiry, got %d", tokenCallCount)
		}
	})

	t.Run("with scopes", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			receivedScopes := r.FormValue("scope")

			if receivedScopes != testScopes {
				t.Errorf("Expected scope=%s, got %s", testScopes, receivedScopes)
			}

			response := TokenResponse{
				AccessToken: "token-with-scopes",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
				Scope:       receivedScopes,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer mockServer.Close()

		client := &OAuthClient{
			config: &OAuthConfig{
				TokenURL:       mockServer.URL,
				ClientID:       testClientID,
				ClientPassword: testClientSecret,
				Scopes:         testScopes,
			},
		}

		token, err := client.GetToken()
		if err != nil {
			t.Fatalf("GetToken() with scopes error = %v", err)
		}
		if token != "token-with-scopes" {
			t.Errorf("Expected token-with-scopes, got %s", token)
		}
	})

	t.Run("without scopes", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			if r.FormValue("scope") != "" {
				t.Errorf("Expected no scope parameter, but got: %s", r.FormValue("scope"))
			}

			response := TokenResponse{
				AccessToken: "token-no-scopes",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer mockServer.Close()

		client := &OAuthClient{
			config: &OAuthConfig{
				TokenURL:       mockServer.URL,
				ClientID:       testClientID,
				ClientPassword: testClientSecret,
				Scopes:         "",
			},
		}

		token, err := client.GetToken()
		if err != nil {
			t.Fatalf("GetToken() without scopes error = %v", err)
		}
		if token != "token-no-scopes" {
			t.Errorf("Expected token-no-scopes, got %s", token)
		}
	})

	t.Run("nil client", func(t *testing.T) {
		var client *OAuthClient
		token, err := client.GetToken()
		if err != nil {
			t.Errorf("GetToken() on nil client should not return error, got %v", err)
		}
		if token != "" {
			t.Errorf("GetToken() on nil client should return empty string, got %s", token)
		}
	})

	t.Run("error cases", func(t *testing.T) {
		tests := []struct {
			name           string
			serverResponse func(w http.ResponseWriter, r *http.Request)
			expectError    bool
			errorContains  string
		}{
			{
				name: "server returns 401",
				serverResponse: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				},
				expectError:   true,
				errorContains: "token request failed with status: 401",
			},
			{
				name: "server returns 500",
				serverResponse: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				},
				expectError:   true,
				errorContains: "token request failed with status: 500",
			},
			{
				name: "server returns invalid JSON",
				serverResponse: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("not valid json"))
				},
				expectError:   true,
				errorContains: "failed to decode token response",
			},
			{
				name: "server timeout",
				serverResponse: func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(35 * time.Second)
				},
				expectError:   true,
				errorContains: "failed to send token request",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.name == "server timeout" && testing.Short() {
					t.Skip("Skipping timeout test in short mode")
				}

				mockServer := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
				defer mockServer.Close()

				client := &OAuthClient{
					config: &OAuthConfig{
						TokenURL:       mockServer.URL,
						ClientID:       testClientID,
						ClientPassword: testClientSecret,
					},
				}

				_, err := client.GetToken()
				assertError(t, err, tt.expectError, tt.errorContains)
			})
		}
	})
}

func TestOAuthClient_AddAuthHeader(t *testing.T) {
	t.Run("adds auth header", func(t *testing.T) {
		mockServer := createMockOAuthServer(t, testBearerToken, 3600)
		defer mockServer.Close()

		client := &OAuthClient{
			config: &OAuthConfig{
				TokenURL:       mockServer.URL,
				ClientID:       testClientID,
				ClientPassword: testClientSecret,
			},
		}

		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		err = client.AddAuthHeader(req)
		if err != nil {
			t.Fatalf("AddAuthHeader() error = %v", err)
		}

		authHeader := req.Header.Get("Authorization")
		expected := "Bearer " + testBearerToken
		if authHeader != expected {
			t.Errorf("Expected Authorization header %s, got %s", expected, authHeader)
		}
	})

	t.Run("nil client", func(t *testing.T) {
		var client *OAuthClient
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		err = client.AddAuthHeader(req)
		if err != nil {
			t.Errorf("AddAuthHeader() on nil client should not return error, got %v", err)
		}

		if req.Header.Get("Authorization") != "" {
			t.Errorf("AddAuthHeader() on nil client should not add header, got %s", req.Header.Get("Authorization"))
		}
	})
}

func TestOAuthClient_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		client   *OAuthClient
		expected bool
	}{
		{
			name:     "nil client",
			client:   nil,
			expected: false,
		},
		{
			name: "initialized client",
			client: &OAuthClient{
				config: &OAuthConfig{
					TokenURL:       testTokenURL,
					ClientID:       testClientID,
					ClientPassword: testClientSecret,
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.client.IsEnabled(); result != tt.expected {
				t.Errorf("IsEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOAuthClient_HTTPClient(t *testing.T) {
	mockServer := createMockOAuthServer(t, "http-client-token", 3600)
	defer mockServer.Close()

	client := &OAuthClient{
		config: &OAuthConfig{
			TokenURL:       mockServer.URL,
			ClientID:       testClientID,
			ClientPassword: testClientSecret,
		},
	}

	httpClient := client.HTTPClient()
	if httpClient == nil {
		t.Fatal("HTTPClient() returned nil")
	}

	if httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout of 30s, got %v", httpClient.Timeout)
	}

	transport, ok := httpClient.Transport.(*OAuthTransport)
	if !ok {
		t.Fatal("HTTPClient() transport is not OAuthTransport")
	}
	if transport.OAuthClient != client {
		t.Error("OAuthTransport does not reference the correct OAuth client")
	}
}

func TestOAuthTransport_RoundTrip(t *testing.T) {
	t.Run("adds auth header to request", func(t *testing.T) {
		oauthCallCount := 0
		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			oauthCallCount++
			response := TokenResponse{
				AccessToken: "transport-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer mockOAuthServer.Close()

		authHeaderReceived := ""
		mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderReceived = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}))
		defer mockAPIServer.Close()

		oauthClient := &OAuthClient{
			config: &OAuthConfig{
				TokenURL:       mockOAuthServer.URL,
				ClientID:       testClientID,
				ClientPassword: testClientSecret,
			},
		}

		transport := &OAuthTransport{
			Base:        http.DefaultTransport,
			OAuthClient: oauthClient,
		}

		req, err := http.NewRequest("GET", mockAPIServer.URL, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := transport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip() error = %v", err)
		}
		defer resp.Body.Close()

		if oauthCallCount != 1 {
			t.Errorf("Expected 1 OAuth call, got %d", oauthCallCount)
		}

		expected := "Bearer transport-token"
		if authHeaderReceived != expected {
			t.Errorf("Expected Authorization header %s, got %s", expected, authHeaderReceived)
		}

		// Verify original request was not modified
		if req.Header.Get("Authorization") != "" {
			t.Error("Original request should not be modified")
		}
	})

	t.Run("handles OAuth failure", func(t *testing.T) {
		oauthClient := &OAuthClient{
			config: &OAuthConfig{
				TokenURL:       "http://invalid-url-that-does-not-exist",
				ClientID:       testClientID,
				ClientPassword: testClientSecret,
			},
		}

		transport := &OAuthTransport{
			Base:        http.DefaultTransport,
			OAuthClient: oauthClient,
		}

		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		_, err = transport.RoundTrip(req)
		if err == nil {
			t.Error("RoundTrip() expected error when OAuth fails")
		}
	})
}