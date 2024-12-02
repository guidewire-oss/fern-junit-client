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
		{
			name: "combined reports",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectName: testProjectName,
				filePattern: reportsCombinedPattern,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "failed report",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectName: testProjectName,
				filePattern: reportFailedPath,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "passed report",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectName: testProjectName,
				filePattern: reportPassedPath,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "no reports",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectName: testProjectName,
				filePattern: nonExistentFilePath,
				verbose:     true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendReports(tt.args.fernUrl, tt.args.projectName, tt.args.filePattern, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("SendReports() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
