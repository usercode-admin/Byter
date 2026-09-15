//go:build linux

package cmd

import (
    "byter/core"
    "byter/utils"
    "fmt"
    "net"
    "os"
    "runtime"
    "runtime/debug"
)

func Boot() {
    runtime.GOMAXPROCS(4)
    debug.SetGCPercent(200)
    debug.SetMaxThreads(10000)

    PrintBanner()
    fmt.Print("  Wait a moment...")

    state := core.NewEngineState()

    if os.Geteuid() != 0 {
        fmt.Println()
        fmt.Println("[!] Root or CAP_NET_RAW is required to open a raw socket.")
        os.Exit(1)
    }

    fd, err := utils.NewRawSocket()
    if err != nil {
        fmt.Println()
        fmt.Printf("[!] init raw socket fail: %v\n", err)
        os.Exit(1)
    }
    state.RawFD = fd

    ip, err := net.ResolveIPAddr("ip4", state.Global.Target)
    if err != nil {
        fmt.Println()
        fmt.Printf("[!] resolve default target err: %v\n", err)
        os.Exit(1)
    }
    state.ResolvedTarget = ip.IP

    state.PrebuiltIPv4 = utils.BuildIPv4(utils.IPv4Header{
        VersionIHL: 4,
        TTL:        64,
        Protocol:   6,
        Src:        utils.AddrTo4(net.IPv4(127, 0, 0, 1).To4()),
        Dst:        utils.AddrTo4(ip.IP),
    })

    fmt.Println(" done")
    fmt.Printf("[*] raw socket fd=%d | target=%s | uid=%d | gomaxprocs=%d\n",
        state.RawFD, state.ResolvedTarget.String(), os.Geteuid(), runtime.GOMAXPROCS(0))

    PrintExamples()
    RunShell(state)
}