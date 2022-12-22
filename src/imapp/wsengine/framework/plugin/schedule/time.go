package schedule

import (
	"runtime"
	"sync/atomic"
	"time"
)

type CancelFunc func()

const (
	stateInit = iota
	stateReady
	stateDone
)

func Interval(delay, interval time.Duration, fn func()) CancelFunc {
	var t *time.Timer
	var state int32
	t = time.AfterFunc(delay, func() {
		for atomic.LoadInt32(&state) == stateInit {
			runtime.Gosched()
		}

		if state == stateDone {
			return
		}

		fn()
		t.Reset(interval)
	})

	// ensures t != nil and is required to avoid data race in
	// AfterFunc calling t.Reset
	atomic.StoreInt32(&state, stateReady)

	return func() {
		if atomic.SwapInt32(&state, stateDone) != stateDone {
			t.Stop()
		}
	}
}

func Once(delay time.Duration, fn func()) CancelFunc {
	t := time.AfterFunc(delay, fn)

	return func() { t.Stop() }
}
