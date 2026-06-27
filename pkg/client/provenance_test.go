package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

// captureServer records the TestRun payload of the last request it receives.
func captureServer(t *testing.T, got *fern.TestRun) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, got)
		w.WriteHeader(http.StatusOK)
	}))
}

// SendReports records the provenance it is given onto the payload (env-var
// resolution is the caller's job — see cmd).
func TestSendReportsRecordsProvenance(t *testing.T) {
	var got fern.TestRun
	srv := captureServer(t, &got)
	defer srv.Close()

	if err := SendReports(SendOptions{FernURL: srv.URL, ProjectID: testProjectId, FilePattern: reportPassedPath, Tags: exampleTags, Branch: "feature/x", CommitSha: "abc123", Environment: "staging", Verbose: false}); err != nil {
		t.Fatalf("SendReports() error = %v", err)
	}
	if got.GitBranch != "feature/x" || got.GitSha != "abc123" || got.Environment != "staging" {
		t.Errorf("payload provenance = {branch:%q sha:%q env:%q}, want {feature/x abc123 staging}",
			got.GitBranch, got.GitSha, got.Environment)
	}
}
