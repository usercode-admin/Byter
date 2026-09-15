package core

import (
    "fmt"
    "os"
    "time"
)

func LiveCounter(st *EngineState, stop <-chan struct{}) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    var lastSucc, lastFail uint64
    var lastTime = time.Now()

    for {
        select {
        case <-stop:
            succ, fail, _, conns, _ := st.Counters.Snapshot()
            fmt.Fprintf(os.Stderr, "\r[succ: %d / fail: %d]  conns: %d\n",
                succ, fail, conns)
            return
        case <-ticker.C:
            succ, fail, _, conns, _ := st.Counters.Snapshot()
            now := time.Now()
            dt := now.Sub(lastTime).Seconds()
            if dt <= 0 {
                dt = 1
            }
            pps := float64(succ-lastSucc) / dt
            fps := float64(fail-lastFail) / dt
            lastSucc, lastFail = succ, fail
            lastTime = now

            fmt.Fprintf(os.Stderr,
                "\r[succ: %d / fail: %d]  pps: %.0f  fps: %.0f  conns: %d",
                succ, fail, pps, fps, conns)
        }
    }
}