package core

import (
	"sync/atomic"
	"time"
)

// Escal
type Escalator struct {
	lastSucc uint64
	lastFail uint64
	lastTime time.Time
	step     float64
	maxRate  int64
}

func NewEscalator(maxRate int64) *Escalator {
	return &Escalator{
		lastTime: time.Now(),
		step:     1.15,
		maxRate:  maxRate,
	}
}

// Tick
func (e *Escalator) Tick(c *Counters) float64 {
	succ, fail, _, _, _ := c.Snapshot()
	ds := succ - e.lastSucc
	df := fail - e.lastFail
	e.lastSucc = succ
	e.lastFail = fail

	total := ds + df
	if total == 0 {
		return 1.0
	}
	errRate := float64(df) / float64(total)

	if errRate < 0.05 {
		return e.step
	}
	if errRate > 0.20 {
		return 1 / e.step
	}
	return 1.0
}

var _ = atomic.LoadUint64