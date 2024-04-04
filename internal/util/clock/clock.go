package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (c *RealClocker) Now() time.Time {
	return time.Now()
}

type MockClocker struct{}

func (c *MockClocker) Now() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
