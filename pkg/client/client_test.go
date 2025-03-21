package client

import (
	"testing"
)

func TestSendReports(t *testing.T) {
	type args struct {
		fernUrl     string
		projectId   string
		filePattern string
		tags        string
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
				projectId:   testProjectId,
				filePattern: reportsCombinedPattern,
				tags:        exampleTags,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "failed report",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectId:   testProjectId,
				filePattern: reportFailedPath,
				tags:        exampleTags,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "passed report",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectId:   testProjectId,
				filePattern: reportPassedPath,
				tags:        exampleTags,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "no reports",
			args: args{
				fernUrl:     mockFernReporter.URL,
				projectId:   testProjectId,
				filePattern: nonExistentFilePath,
				tags:        exampleTags,
				verbose:     true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendReports(tt.args.fernUrl, tt.args.projectId, tt.args.filePattern, tt.args.tags, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("SendReports() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
