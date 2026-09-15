package core

import (
    "sync/atomic"
    "time"
)

type Counters struct {
    Success     uint64
    Fail        uint64
    BytesSent   uint64
    ActiveConns int64
    StartTime   time.Time
}

func NewCounters() *Counters {
    return &Counters{StartTime: time.Now()}
}

func (c *Counters) IncSuccess(n int) {
    atomic.AddUint64(&c.Success, 1)
    atomic.AddUint64(&c.BytesSent, uint64(n))
}

func (c *Counters) IncFail() {
    atomic.AddUint64(&c.Fail, 1)
}

func (c *Counters) IncPacket(n int) { c.IncSuccess(n) }
func (c *Counters) IncError()       { c.IncFail() }
func (c *Counters) ConnAdd(d int64) { atomic.AddInt64(&c.ActiveConns, d) }

func (c *Counters) Reset() {
    atomic.StoreUint64(&c.Success, 0)
    atomic.StoreUint64(&c.Fail, 0)
    atomic.StoreUint64(&c.BytesSent, 0)
    atomic.StoreInt64(&c.ActiveConns, 0)
    c.StartTime = time.Now()
}

func (c *Counters) Snapshot() (succ, fail, bytes uint64, conns int64, uptime time.Duration) {
    return atomic.LoadUint64(&c.Success),
        atomic.LoadUint64(&c.Fail),
        atomic.LoadUint64(&c.BytesSent),
        atomic.LoadInt64(&c.ActiveConns),
        time.Since(c.StartTime)
}