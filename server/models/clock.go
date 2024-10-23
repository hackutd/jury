package models

import (
	"sync"
	"time"
)

type ClockState struct {
	StartTime int64 `json:"start_time" bson:"start_time"`
	PauseTime int64 `json:"pause_time" bson:"pause_time"`
	Running   bool  `json:"running" bson:"running"`
}

func NewClockState() *ClockState {
	return &ClockState{
		StartTime: 0,
		PauseTime: 0,
		Running:   false,
	}
}

// Gets the current clock time in milliseconds
func GetCurrTime() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}

func (c *ClockState) Pause() {
	if !c.Running {
		return
	}
	c.Running = false
	c.PauseTime = GetCurrTime() - c.StartTime
}

func (c *ClockState) Resume() {
	if c.Running {
		return
	}
	c.Running = true
	c.StartTime = GetCurrTime()
}

func (c *ClockState) Reset() {
	c.StartTime = 0
	c.PauseTime = 0
	c.Running = false
}

func (c *ClockState) GetDuration() int64 {
	if !c.Running {
		return c.PauseTime
	}
	return c.PauseTime + GetCurrTime() - c.StartTime
}

// SafeClock wraps ClockState in a mutex so it can be used safely across threads
type SafeClock struct {
	Mutex sync.Mutex
	Clock ClockState
}

func NewSafeClock(clock *ClockState) *SafeClock {
	return &SafeClock{
		Mutex: sync.Mutex{},
		Clock: *clock,
	}
}
