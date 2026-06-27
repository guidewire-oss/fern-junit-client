package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

func TestFirstNonEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"first non-empty wins", []string{"a", "b"}, "a"},
		{"skips leading empties", []string{"", "", "c"}, "c"},
		{"all empty", []string{"", ""}, ""},
		{"no candidates", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstNonEmpty(tt.in...); got != tt.want {
				t.Errorf("firstNonEmpty(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// captureServer records the TestRun payload of the last request it receives.
func captureServer(t *testing.T, got *fern.TestRun) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, got)
		w.WriteHeader(http.StatusOK)
	}))
}

func TestSendReportsProvenance(t *testing.T) {
	t.Run("explicit flags populate the payload", func(t *testing.T) {
		var got fern.TestRun
		srv := captureServer(t, &got)
		defer srv.Close()

		if err := SendReports(srv.URL, testProjectId, reportPassedPath, exampleTags, "feature/x", "abc123", "staging", false); err != nil {
			t.Fatalf("SendReports() error = %v", err)
		}
		if got.GitBranch != "feature/x" || got.GitSha != "abc123" || got.Environment != "staging" {
			t.Errorf("payload provenance = {branch:%q sha:%q env:%q}, want {feature/x abc123 staging}",
				got.GitBranch, got.GitSha, got.Environment)
		}
	})

	t.Run("empty flags fall back to CI env vars", func(t *testing.T) {
		t.Setenv("GITHUB_REF_NAME", "main")
		t.Setenv("GITHUB_SHA", "deadbeef")
		t.Setenv("CI_ENVIRONMENT_NAME", "prod")

		var got fern.TestRun
		srv := captureServer(t, &got)
		defer srv.Close()

		if err := SendReports(srv.URL, testProjectId, reportPassedPath, exampleTags, "", "", "", false); err != nil {
			t.Fatalf("SendReports() error = %v", err)
		}
		if got.GitBranch != "main" || got.GitSha != "deadbeef" || got.Environment != "prod" {
			t.Errorf("env fallback = {branch:%q sha:%q env:%q}, want {main deadbeef prod}",
				got.GitBranch, got.GitSha, got.Environment)
		}
	})

	t.Run("falls back to secondary (GitLab/FERN) env vars", func(t *testing.T) {
		// Null the primary GitHub vars so the secondary candidates win — robust
		// even on a GitHub Actions runner where GITHUB_REF_NAME/SHA are set.
		t.Setenv("GITHUB_REF_NAME", "")
		t.Setenv("GITHUB_SHA", "")
		t.Setenv("CI_ENVIRONMENT_NAME", "")
		t.Setenv("CI_COMMIT_REF_NAME", "feature/gl")
		t.Setenv("CI_COMMIT_SHA", "cafef00d")
		t.Setenv("FERN_ENVIRONMENT", "review")

		var got fern.TestRun
		srv := captureServer(t, &got)
		defer srv.Close()

		if err := SendReports(srv.URL, testProjectId, reportPassedPath, exampleTags, "", "", "", false); err != nil {
			t.Fatalf("SendReports() error = %v", err)
		}
		if got.GitBranch != "feature/gl" || got.GitSha != "cafef00d" || got.Environment != "review" {
			t.Errorf("secondary fallback = {branch:%q sha:%q env:%q}, want {feature/gl cafef00d review}",
				got.GitBranch, got.GitSha, got.Environment)
		}
	})

	t.Run("explicit flag overrides env var", func(t *testing.T) {
		t.Setenv("GITHUB_REF_NAME", "main")

		var got fern.TestRun
		srv := captureServer(t, &got)
		defer srv.Close()

		if err := SendReports(srv.URL, testProjectId, reportPassedPath, exampleTags, "release/1.0", "", "", false); err != nil {
			t.Fatalf("SendReports() error = %v", err)
		}
		if got.GitBranch != "release/1.0" {
			t.Errorf("branch = %q, want release/1.0 (explicit flag must override env)", got.GitBranch)
		}
	})
}
