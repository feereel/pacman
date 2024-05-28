package network

import (
	"time"
)

type Netclock struct {
	FrameHole    float64
	StartTime    time.Time
	FrameTimeout time.Duration

	SafeFrameTime   time.Duration
	UnsafeFrameTime time.Duration

	IsSafe        bool
	WaitingUntill time.Time
}

func NewNetclock(FrameTimeout int64, FrameHole float64) *Netclock {
	var currentTime = time.Now()
	var dif = int64(float64(FrameTimeout) * FrameHole)

	var c = Netclock{
		FrameHole:    FrameHole,
		StartTime:    currentTime,
		FrameTimeout: time.Duration(FrameTimeout) * time.Millisecond,

		SafeFrameTime:   time.Duration(FrameTimeout-dif) * time.Millisecond,
		UnsafeFrameTime: time.Duration(dif) * time.Millisecond,

		WaitingUntill: currentTime,
	}
	go c.StartClock()

	return &c
}

func (c *Netclock) StartClock() {
	for {
		c.WaitingUntill = c.WaitingUntill.Add(c.SafeFrameTime)
		c.IsSafe = true
		time.Sleep(time.Until(c.WaitingUntill))

		c.WaitingUntill = c.WaitingUntill.Add(c.UnsafeFrameTime)
		c.IsSafe = false
		time.Sleep(time.Until(c.WaitingUntill))
	}
}

func (c *Netclock) WaitUntilSafeFrame() {
	if !c.IsSafe {
		time.Sleep(time.Until(c.WaitingUntill))
	}
}
