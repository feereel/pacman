package network

import (
	"time"
)

type Netclock struct {
	FrameHole    float64
	StartTime    time.Time
	NextFrameEnd time.Time
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
		NextFrameEnd: currentTime,
		FrameTimeout: time.Duration(FrameTimeout) * time.Millisecond,

		SafeFrameTime:   time.Duration(FrameTimeout-dif-dif) * time.Millisecond,
		UnsafeFrameTime: time.Duration(dif) * time.Millisecond,

		WaitingUntill: currentTime,
	}
	go c.StartClock()

	return &c
}

func (c *Netclock) StartClock() {
	c.WaitingUntill = c.WaitingUntill.Add(c.UnsafeFrameTime)
	c.IsSafe = false
	time.Sleep(time.Until(c.WaitingUntill))
	for {
		c.NextFrameEnd = c.NextFrameEnd.Add(c.FrameTimeout)

		c.WaitingUntill = c.WaitingUntill.Add(c.SafeFrameTime)
		c.IsSafe = true
		time.Sleep(time.Until(c.WaitingUntill))

		c.WaitingUntill = c.WaitingUntill.Add(c.UnsafeFrameTime * 2)
		c.IsSafe = false
		time.Sleep(time.Until(c.WaitingUntill))
	}
}

func (c *Netclock) WaitUntilSafeFrame() {
	if !c.IsSafe {
		time.Sleep(time.Until(c.WaitingUntill))
	}
}
