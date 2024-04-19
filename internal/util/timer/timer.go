package timer

import "time"

type Timer interface {
	Now() time.Time
}

type RealTimer struct{}

func (t *RealTimer) Now() time.Time {
	return time.Now()
}

type MockTimer struct{}

func (t *MockTimer) Now() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
