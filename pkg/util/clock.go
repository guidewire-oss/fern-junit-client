package util

import "time"

var GlobalClock Clock = RealClock{}

type Clock interface {
	Now() time.Time
	// Add other time methods as necessary
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// MockClock uses a fixed time for Now(), which enables deep equality testing for
// report parsing functions when input suites don't include timestamp values
type MockClock struct {
	fixedTime time.Time
}

func (m MockClock) Now() time.Time {
	return m.fixedTime
}

func NewMockClock(t ...time.Time) *MockClock {
	if len(t) > 0 {
		return &MockClock{fixedTime: t[0]}
	}
	fixedTime, _ := time.Parse(time.RFC3339, time.RFC3339)
	return &MockClock{fixedTime: fixedTime}
}
