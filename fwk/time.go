package fwk

import (
	"time"
)

type Timer struct {
	last time.Time
}

func (t Timer) Delta() time.Duration {
	now := time.Now()
	delta := now.Sub(t.last)
	t.last = now
	return delta
}

func NewTimer() Timer {
	return Timer{time.Now()}
}
