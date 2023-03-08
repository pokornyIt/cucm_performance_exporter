package main

import (
	"sync"
	"time"
)

const (
	RateStandardDelay     = time.Second + time.Millisecond*200
	RateStandardTestDelay = RateStandardDelay + time.Millisecond*300
	RateBaseWaitTime      = time.Minute + time.Millisecond*200
	RateRequestLimit      = 50
)

type RateControl struct {
	start    time.Time
	requests int
	mutex    sync.Mutex
}

func (r *RateControl) add() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.requests++
}

func (r *RateControl) reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.start = time.Now()
	r.requests = 0
}

func (r *RateControl) delay() time.Duration {
	defer r.add()
	timePeriod := time.Now().Sub(r.start)
	waitTime := RateBaseWaitTime - timePeriod
	if r.requests < RateRequestLimit {
		if timePeriod > RateStandardTestDelay*time.Duration(r.requests) {
			r.reset()
			return time.Millisecond
		}
		return RateStandardDelay
	}
	if timePeriod > time.Minute {
		waitTime = time.Millisecond
	}
	r.reset()
	return waitTime
}
