package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	TokenURL       string
	ClientID       string
	ClientPassword string
	Scopes         string
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// OAuthClient handles OAuth authentication
type OAuthClient struct {
	config      *OAuthConfig
	token       *TokenResponse
	tokenExpiry time.Time
}

// NewOAuthClient creates a new OAuth client
// Returns nil if OAuth is not configured (no AUTH_URL set)
// Returns an error if OAuth is partially configured (missing required parameters)
func NewOAuthClient() (*OAuthClient, error) {
	config := &OAuthConfig{
		TokenURL:       os.Getenv("AUTH_URL"),
		ClientID:       os.Getenv("FERN_AUTH_CLIENT_ID"),
		ClientPassword: os.Getenv("FERN_AUTH_CLIENT_SECRET"),
		Scopes:         os.Getenv("FERN_CLIENT_SCOPE"),
	}

	// If token URL is not set, OAuth is disabled - this is OK
	if config.TokenURL == "" {
		return nil, nil
	}

	// If AUTH_URL is set, validate that we have all required OAuth parameters
	var missingParams []string
	if config.ClientID == "" {
		missingParams = append(missingParams, "FERN_AUTH_CLIENT_ID")
	}
	if config.ClientPassword == "" {
		missingParams = append(missingParams, "FERN_AUTH_CLIENT_SECRET")
	}

	if len(missingParams) > 0 {
		return nil, fmt.Errorf("OAuth configuration error: AUTH_URL is set but missing required parameters: %s", strings.Join(missingParams, ", "))
	}

	return &OAuthClient{
		config: config,
	}, nil
}

// GetToken fetches a new OAuth token or returns the cached one if still valid
func (c *OAuthClient) GetToken() (string, error) {
	if c == nil {
		return "", nil
	}

	// Check if we have a valid cached token
	if c.token != nil && time.Now().Before(c.tokenExpiry) {
		return c.token.AccessToken, nil
	}

	// Fetch new token
	if err := c.fetchToken(); err != nil {
		return "", fmt.Errorf("failed to fetch OAuth token: %w", err)
	}

	return c.token.AccessToken, nil
}

// fetchToken fetches a new OAuth token from the authorization server
func (c *OAuthClient) fetchToken() error {
	// Prepare the token request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientPassword)

	// Add scopes if provided
	if c.config.Scopes != "" {
		data.Set("scope", c.config.Scopes)
	}

	req, err := http.NewRequest("POST", c.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	// Parse the response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	// Store the token and calculate expiry
	c.token = &tokenResp
	// Subtract 30 seconds from expiry to ensure we refresh before it actually expires
	expiryDuration := time.Duration(tokenResp.ExpiresIn-30) * time.Second
	c.tokenExpiry = time.Now().Add(expiryDuration)

	return nil
}

// AddAuthHeader adds the OAuth bearer token to the request header if OAuth is enabled
func (c *OAuthClient) AddAuthHeader(req *http.Request) error {
	if c == nil {
		// OAuth is not enabled
		return nil
	}

	token, err := c.GetToken()
	if err != nil {
		return err
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return nil
}

// IsEnabled returns whether OAuth is enabled
func (c *OAuthClient) IsEnabled() bool {
	return c != nil
}

// HTTPClient returns an http.Client with OAuth authentication if enabled
func (c *OAuthClient) HTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &OAuthTransport{
			Base:        http.DefaultTransport,
			OAuthClient: c,
		},
	}
}

// OAuthTransport is an http.RoundTripper that adds OAuth authentication
type OAuthTransport struct {
	Base        http.RoundTripper
	OAuthClient *OAuthClient
}

// RoundTrip implements the http.RoundTripper interface
func (t *OAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqCopy := req.Clone(req.Context())

	// Add OAuth header
	if err := t.OAuthClient.AddAuthHeader(reqCopy); err != nil {
		return nil, err
	}

	// Use the base transport to send the request
	return t.Base.RoundTrip(reqCopy)
}
