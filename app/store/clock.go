package store

import "time"

type Clock interface {
	CurrentMillis() int
}

type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

func (c *SystemClock) CurrentMillis() int {
	return int(time.Now().UnixMilli())
}
