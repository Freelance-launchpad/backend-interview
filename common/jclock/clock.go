package jclock

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	now time.Time
}

func NewFakeClock(now time.Time) FakeClock {
	return FakeClock{now: now}
}

func (f FakeClock) Now() time.Time {
	return f.now
}
