package client

import (
	"testing"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
)

func Test_sendTestRun(t *testing.T) {
	type args struct {
		testRun fern.TestRun
		fernUrl string
		verbose bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendTestRun(tt.args.testRun, tt.args.fernUrl, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("sendTestRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
