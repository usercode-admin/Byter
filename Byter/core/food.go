package core

import (
    "sync/atomic"
    "time"
)

// RateLimit
type RateLimiter struct {
    mode     FloodMode
    baseRate int
    current  int64
    lastFail uint64
    lastSucc uint64
    batch    int
}

func NewRateLimiter(mode FloodMode, baseRate int) *RateLimiter {
    if baseRate < 1 {
        baseRate = 1
    }
    return &RateLimiter{
        mode:     mode,
        baseRate: baseRate,
        current:  int64(baseRate),
        batch:    256,
    }
}

// Wait
func (r *RateLimiter) Wait() {
    switch r.mode {
    case FloodBurst:
        return
    default:
        rate := atomic.LoadInt64(&r.current)
        if rate <= 0 {
            rate = 1
        }
        // sleep theo batch
        us := int64(1_000_000) * int64(r.batch) / rate
        if us > 0 {
            time.Sleep(time.Duration(us) * time.Microsecond)
        }
    }
}

func (r *RateLimiter) Batch() int {
    if r.mode == FloodBurst {
        return 1
    }
    return r.batch
}

func (r *RateLimiter) Adapt(c *Counters) {
    if r.mode != FloodAdaptive {
        return
    }
    succ, fail, _, _, _ := c.Snapshot()
    deltaSucc := succ - r.lastSucc
    deltaFail := fail - r.lastFail
    r.lastSucc = succ
    r.lastFail = fail

    total := deltaSucc + deltaFail
    if total == 0 {
        return
    }
    errRate := float64(deltaFail) / float64(total)

    cur := atomic.LoadInt64(&r.current)
    if errRate > 0.1 {
        cur = int64(float64(cur) * 0.8)
    } else {
        cur = int64(float64(cur) * 1.1)
    }
    if cur < 100 {
        cur = 100
    }
    if cur > 500000 {
        cur = 500000
    }
    atomic.StoreInt64(&r.current, cur)
}

func (r *RateLimiter) Current() int64 {
    return atomic.LoadInt64(&r.current)
}