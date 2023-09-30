package models

import "time"

type ClockState struct {
	StartTime int64 `json:"start_time" bson:"start_time"`
	Prev      int64 `json:"prev" bson:"prev"`
	Running   bool  `json:"running" bson:"running"`
}

func NewClockState() *ClockState {
	return &ClockState{
		StartTime: 0,
		Prev:      0,
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
	c.Prev += GetCurrTime() - c.StartTime
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
	c.Prev = 0
	c.Running = false
}

func (c *ClockState) GetDuration() int64 {
	if !c.Running {
		return c.Prev
	}
	return c.Prev + GetCurrTime() - c.StartTime
}
