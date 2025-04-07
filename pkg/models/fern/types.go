package fern

import (
	"time"
)

type TestRun struct {
	ID            uint64     `json:"id"`
	TestProjectID string     `json:"test_project_id"`
	TestSeed      uint64     `json:"test_seed"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	SuiteRuns     []SuiteRun `json:"suite_runs"`
}

type SuiteRun struct {
	ID        uint64    `json:"id"`
	TestRunID uint64    `json:"test_run_id"`
	SuiteName string    `json:"suite_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SpecRuns  []SpecRun `json:"spec_runs"`
}

type SpecRun struct {
	ID              uint64    `json:"id"`
	SuiteID         uint64    `json:"suite_id"`
	SpecDescription string    `json:"spec_description"`
	Status          string    `json:"status"`
	Message         string    `json:"message"`
	Tags            []Tag     `json:"tags"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
}

type Tag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
