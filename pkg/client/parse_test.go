package client

import (
	"reflect"
	"testing"
	"time"

	"github.com/guidewire-oss/fern-junit-client/pkg/models/fern"
	"github.com/guidewire-oss/fern-junit-client/pkg/models/junit"
)

func Test_parseReports(t *testing.T) {
	type args struct {
		testRun     *fern.TestRun
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
				testRun:     &fern.TestRun{},
				filePattern: reportsCombinedPattern,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "failed report",
			args: args{
				testRun:     &fern.TestRun{},
				filePattern: reportFailedPath,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "passed report",
			args: args{
				testRun:     &fern.TestRun{},
				filePattern: reportPassedPath,
				verbose:     true,
			},
			wantErr: false,
		},
		{
			name: "no reports",
			args: args{
				testRun:     &fern.TestRun{},
				filePattern: nonExistentFilePath,
				verbose:     true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseReports(tt.args.testRun, tt.args.filePattern, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("parseReports() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseReport(t *testing.T) {
	type args struct {
		filePath string
		verbose  bool
	}
	tests := []struct {
		name    string
		args    args
		want    []fern.SuiteRun
		wantErr bool
	}{
		{
			name: "failed report",
			args: args{
				filePath: reportFailedPath,
				verbose:  true,
			},
			want:    fernTestRunFailed.SuiteRuns,
			wantErr: false,
		},
		{
			name: "passed report",
			args: args{
				filePath: reportPassedPath,
				verbose:  true,
			},
			want:    fernTestRunPassed.SuiteRuns,
			wantErr: false,
		},
		{
			name: "non-existent report",
			args: args{
				filePath: nonExistentFilePath,
				verbose:  true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReport(tt.args.filePath, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTestSuite(t *testing.T) {
	type args struct {
		testSuite junit.TestSuite
		verbose   bool
	}
	tests := []struct {
		name         string
		args         args
		wantSuiteRun fern.SuiteRun
		wantErr      bool
	}{
		{
			name: "failed suite",
			args: args{
				testSuite: junitTestSuiteFailed,
				verbose:   true,
			},
			wantSuiteRun: fernTestRunFailed.SuiteRuns[0],
			wantErr:      false,
		},
		{
			name: "passed suite",
			args: args{
				testSuite: junitTestSuitePassed,
				verbose:   true,
			},
			wantSuiteRun: fernTestRunPassed.SuiteRuns[0],
			wantErr:      false,
		},
		{
			name: "empty suite",
			args: args{
				testSuite: junit.TestSuite{},
				verbose:   true,
			},
			wantSuiteRun: fern.SuiteRun{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSuiteRun, err := parseTestSuite(tt.args.testSuite, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTestSuite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSuiteRun, tt.wantSuiteRun) {
				t.Errorf("parseTestSuite() = %v, want %v", gotSuiteRun, tt.wantSuiteRun)
			}
		})
	}
}

func Test_getEndTime(t *testing.T) {
	type args struct {
		startTime       time.Time
		durationSeconds string
	}
	tests := []struct {
		name        string
		args        args
		wantEndTime time.Time
		wantErr     bool
	}{
		{
			name: "10 seconds",
			args: args{
				startTime:       time.Unix(0, 0),
				durationSeconds: "10",
			},
			wantEndTime: time.Unix(10, 0),
			wantErr:     false,
		},
		{
			name: "10.5 seconds",
			args: args{
				startTime:       time.Unix(0, 0),
				durationSeconds: "10.5",
			},
			wantEndTime: time.Unix(10, 500000000),
			wantErr:     false,
		},
		{
			name: "invalid duration",
			args: args{
				startTime:       time.Unix(0, 0),
				durationSeconds: "foo",
			},
			wantEndTime: time.Unix(0, 0),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEndTime, err := getEndTime(tt.args.startTime, tt.args.durationSeconds)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEndTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEndTime, tt.wantEndTime) {
				t.Errorf("getEndTime() = %v, want %v", gotEndTime, tt.wantEndTime)
			}
		})
	}
}
