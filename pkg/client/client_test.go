package client

import "testing"

func TestSendReports(t *testing.T) {
	type args struct {
		fernUrl     string
		projectName string
		filePattern string
		verbose     bool
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
			if err := SendReports(tt.args.fernUrl, tt.args.projectName, tt.args.filePattern, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("SendReports() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
