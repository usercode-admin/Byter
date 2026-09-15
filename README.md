# ⚡ BYTER

<div align="center">

```text
      d8888b. db    db d888888b d88888b d8888b.
      88  `8D `8b  d8' `~~88~~' 88'     88  `8D
      88oooY'  `8bd8'     88    88ooooo 88oobY'
      88~~~b.    88       88    88~~~~~ 88`8b
      88   8D    88       88    88.     88 `88.
      Y8888P'    YP       YP    Y88888P 88   YD
```

### Multi-Layer Network Security Research Tool

**L3 • L4 • L7 • Raw Sockets • HTTP/2 • TLS Fingerprinting • Automation**

<br>

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge\&logo=go\&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-Linux-FCC624?style=for-the-badge\&logo=linux\&logoColor=black)
![License](https://img.shields.io/badge/License-MIT-6366F1?style=for-the-badge)
![Architecture](https://img.shields.io/badge/Architecture-L3%20%7C%20L4%20%7C%20L7-8B5CF6?style=for-the-badge)

</div>

---

## 📌 Overview

**BYTER** is a Linux-focused network security research and traffic-generation framework written in **Go**.

The project is structured around multiple networking layers and combines low-level packet construction with higher-level HTTP client behavior.

```text
                    ┌──────────────────────┐
                    │        BYTER         │
                    │    Interactive CLI   │
                    └──────────┬───────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
         ┌────▼────┐      ┌────▼────┐      ┌────▼────┐
         │   L3    │      │   L4    │      │   L7    │
         │ Network │      │Transport│      │   HTTP  │
         └────┬────┘      └────┬────┘      └────┬────┘
              │                │                │
              └────────────────┼────────────────┘
                               │
                     ┌─────────▼─────────┐
                     │  Core / Runtime   │
                     │ State • Counters  │
                     │ Rate • Scheduling │
                     └───────────────────┘
```

BYTER is intended for **authorized security research, controlled laboratory testing, network-stack experimentation, and educational purposes**.

---

# ✨ Features

## 🌐 Layer 3

The L3 subsystem contains low-level IPv4 packet-generation functionality.

Key components include:

* Raw IPv4 packet construction
* Randomized source-address generation
* TTL variation
* IP identification handling
* Fragment-related packet fields
* Randomized payload pool
* Configurable packet batches
* Linux raw-socket transmission
* Batch transmission through `sendmmsg`

Relevant implementation:

```text
core/l3_fragment.go
utils/raw_socket.go
utils/batch_send.go
utils/tcp_craft.go
```

---

## 🔗 Layer 4

The L4 subsystem focuses on manually constructed IPv4/TCP packets and configurable TCP state fields.

Supported functionality includes:

* TCP flag combinations
* Sequence-number handling
* Source-port randomization
* TCP window manipulation
* TCP option templates
* TTL variation
* TCP checksum generation
* Batched packet transmission

The TCP flag constants include:

```text
FIN
SYN
RST
PSH
ACK
URG
ECE
CWR
```

TCP packet construction is implemented in:

```text
core/l4_state.go
utils/tcp_craft.go
```

---

## 🌍 Layer 7

The L7 subsystem provides an HTTP-oriented traffic engine with configurable client behavior.

It supports three operating modes:

```text
slow
flood
hybrid
```

The `hybrid` mode combines the two L7 execution paths concurrently.

The L7 subsystem includes:

* HTTP client management
* HTTP/1.1 support
* HTTP/2 support
* TLS configuration
* uTLS-based ClientHello selection
* Browser-style request headers
* Session cookies
* JavaScript execution options
* Resource discovery
* Referer handling
* DNS caching
* Proxy-pool support
* Configurable connection concurrency

---

# 🧬 TLS & Browser Profiles

BYTER includes browser-oriented profiles and TLS ClientHello selection.

Available application profiles in `profile.go`:

| Profile | Available |
| :------ | :-------: |
| Chrome  |     ✅     |
| Firefox |     ✅     |
| Safari  |     ✅     |

The L7 flood engine additionally contains a User-Agent pool containing desktop and mobile browser identities.

TLS ClientHello selection is handled through:

```text
github.com/refraction-networking/utls
```

The implementation selects different ClientHello presets according to the selected browser identity.

---

# 🌐 HTTP/2

HTTP/2 support is implemented using:

```text
golang.org/x/net/http2
```

The project contains separate HTTP transport handling for HTTP/2 and HTTP/1.1 operation.

Relevant files:

```text
core/l7/client.go
core/l7/flood.go
```

---

# 🤖 JavaScript Solver

The L7 subsystem provides three JavaScript modes:

```text
none
goja
chromedp
```

### Goja

Uses the embedded JavaScript runtime:

```text
github.com/dop251/goja
```

The implementation extracts inline `<script>` elements and evaluates their contents.

### Chromedp

Uses:

```text
github.com/chromedp/chromedp
```

to create a browser context, navigate to a page and retrieve the resulting document.

### None

Disables JavaScript processing.

---

# 🍪 Sessions & Cookies

The `Session` abstraction stores:

* HTTP cookie jar
* Selected browser profile
* Referer information
* Arbitrary session tokens

Session management is implemented in:

```text
core/l7/jar.go
```

The standard Go `http.CookieJar` implementation is used for cookie persistence.

---

# 🔄 Resource & Referer Handling

When resource processing is enabled, the L7 request flow can inspect returned HTML and identify:

```text
<link ...>
<script ...>
<img ...>
```

The resource parser resolves discovered URLs relative to the original page.

Resource types are classified as:

```text
style
script
image
empty
```

Relevant implementation:

```text
core/l7/request.go
```

---

# 🌐 DNS Cache

BYTER contains a small in-memory DNS cache.

```text
DNSCache
├── host → []IP
├── timestamp tracking
└── configurable expiration
```

The current implementation uses a **60-second cache lifetime**.

It also provides random selection when multiple addresses are returned.

Implementation:

```text
core/l7/dns.go
```

---

# 🔁 Proxy Pool

The L7 package includes a simple rotating proxy pool.

```text
ProxyPool
├── proxy list
├── mutex protection
└── round-robin selection
```

A selected proxy can be applied directly to an `http.Transport`.

Implementation:

```text
core/l7/proxy.go
```

---

# ⏱️ Rate Control

The core engine contains a configurable rate limiter with three modes:

| Mode       | Purpose                                              |
| :--------- | :--------------------------------------------------- |
| `steady`   | Controlled continuous pacing                         |
| `burst`    | Immediate batch processing                           |
| `adaptive` | Runtime adjustment based on observed success/failure |

The adaptive limiter observes counters and modifies its current rate within internal bounds.

Implementation:

```text
core/food.go
```

---

# 📈 Adaptive Escalation

The project also contains an `Escalator` component.

It observes recent success/failure deltas and calculates an adjustment factor.

Conceptually:

```text
             Runtime results
                    │
                    ▼
             ┌─────────────┐
             │  Escalator  │
             └──────┬──────┘
                    │
          ┌─────────┼─────────┐
          ▼         ▼         ▼
       Increase   Stable   Decrease
```

Implementation:

```text
core/escalate.go
```

---

# 🌊 Pulse Wave

`PulseWave` provides automated ON/OFF execution cycles.

The scheduler supports:

* Pulse ON duration
* Pulse OFF duration
* Timing jitter
* Multiple configured vectors
* Temporary runtime configuration changes
* Automatic restoration of L3/L4 configuration

Supported vector names include:

```text
l3
l3-frag

l4
l4-state

l7
l7-heavy
```

Architecture:

```text
       ┌───────────────┐
       │   Pulse ON    │
       └───────┬───────┘
               │
          Run vectors
               │
       ┌───────▼───────┐
       │   Pulse OFF   │
       └───────┬───────┘
               │
               ▼
            Repeat
```

Implementation:

```text
core/pulse_wave.go
```

---

# 📊 Runtime Counters

The core state contains a shared counter subsystem for runtime statistics.

The shell exposes:

```text
status
```

which reports values such as:

```text
mode
uptime
success count
failure count
bytes
active connections
```

The project also provides a live terminal counter through:

```text
core/live_counter.go
```

---

# 🖥️ Interactive CLI

BYTER launches into an interactive shell after initialization.

```text
shell >
```

The CLI provides:

```text
help
show
status
set
stop
banner
run
exit
quit
q
```

Numeric shortcuts are also available:

```text
1 → L3
2 → L4
3 → L7
4 → Pulse Wave
```

The CLI implementation is located in:

```text
cmd/shell.go
```

---

# ⚙️ Configuration Model

The runtime state is represented by `EngineState`.

```text
EngineState
│
├── GlobalConfig
├── L3Config
├── L4Config
├── L7Config
├── WaveConfig
│
├── RawFD
├── ResolvedTarget
├── PrebuiltIPv4
│
├── Counters
├── StopCurrent
└── CurrentMode
```

This keeps the CLI configuration and individual network engines connected through a shared runtime state.

Implementation:

```text
core/engine.go
```

---

# 📁 Project Structure

```text
byter/
│
├── main.go
├── go.mod
├── config.txt
│
├── cmd/
│   ├── banner.go
│   ├── init.go
│   ├── menu.go
│   └── shell.go
│
├── core/
│   ├── counter.go
│   ├── engine.go
│   ├── escalate.go
│   ├── food.go
│   ├── l3_fragment.go
│   ├── l4_state.go
│   ├── live_counter.go
│   ├── l7_asymmetric.go
│   ├── pulse_wave.go
│   │
│   └── l7/
│       ├── client.go
│       ├── dns.go
│       ├── flood.go
│       ├── headers.go
│       ├── jar.go
│       ├── profile.go
│       ├── proxy.go
│       ├── request.go
│       ├── slow_header.go
│       ├── slow_post.go
│       └── solver.go
│
└── utils/
    ├── batch_send.go
    ├── raw_socket.go
    └── tcp_craft.go
```

---

# 🧩 Module Reference

## `main.go`

Application entry point.

```text
main()
 └── cmd.Boot()
```

---

## `cmd/`

Responsible for startup, CLI interaction and user-facing output.

| File        | Responsibility                            |
| :---------- | :---------------------------------------- |
| `banner.go` | BYTER ASCII banner                        |
| `init.go`   | Linux initialization and raw socket setup |
| `menu.go`   | Built-in examples                         |
| `shell.go`  | Interactive command shell                 |

---

## `core/`

Contains the main runtime engines.

| File               | Responsibility                  |
| :----------------- | :------------------------------ |
| `engine.go`        | Runtime state and configuration |
| `counter.go`       | Statistics counters             |
| `live_counter.go`  | Live statistics display         |
| `food.go`          | Rate limiter                    |
| `escalate.go`      | Rate adjustment logic           |
| `l3_fragment.go`   | L3 packet engine                |
| `l4_state.go`      | L4 packet engine                |
| `l7_asymmetric.go` | L7 orchestration                |
| `pulse_wave.go`    | Pulse scheduling                |

---

## `core/l7/`

Application-layer networking components.

| File             | Responsibility                   |
| :--------------- | :------------------------------- |
| `client.go`      | HTTP client / transport creation |
| `dns.go`         | DNS cache                        |
| `flood.go`       | L7 traffic engine                |
| `headers.go`     | Browser-style headers            |
| `jar.go`         | Session and cookies              |
| `profile.go`     | Browser profiles                 |
| `proxy.go`       | Proxy pool                       |
| `request.go`     | HTTP request/resource flow       |
| `slow_header.go` | Slow-header component            |
| `slow_post.go`   | Slow-POST component              |
| `solver.go`      | JavaScript processing            |

---

## `utils/`

Low-level networking primitives.

| File            | Responsibility                             |
| :-------------- | :----------------------------------------- |
| `raw_socket.go` | Linux raw socket creation                  |
| `batch_send.go` | Linux batch packet transmission            |
| `tcp_craft.go`  | IPv4/TCP packet construction and checksums |

---

# 📦 Dependencies

The project currently declares the following Go dependencies:

| Dependency | Purpose                       |
| :--------- | :---------------------------- |
| `chromedp` | Browser automation            |
| `goja`     | Embedded JavaScript runtime   |
| `utls`     | TLS ClientHello customization |
| `x/net`    | HTTP/2 networking             |
| `x/sys`    | System-level functionality    |

Go module:

```text
module byter

go 1.22
```

---

# 🐧 Platform

BYTER is currently designed around **Linux-specific networking functionality**.

Several components explicitly use Linux build constraints:

```go
//go:build linux
```

This applies to functionality such as:

* Raw sockets
* Linux system calls
* `sendmmsg`
* Low-level packet transmission

The current implementation should therefore be treated as a **Linux-first project**.

---

# 🔐 Permissions

The initialization code checks for elevated privileges before opening the raw socket.

The application expects either:

```text
root
```

or equivalent networking capabilities such as:

```text
CAP_NET_RAW
```

The exact capability configuration depends on the Linux environment.

---

# 🚀 Build

Clone the repository:

```bash
git clone https://github.com/yourusername/byter.git
cd byter
```

Download Go dependencies:

```bash
go mod tidy
```

Build:

```bash
go build -o byte .
```

Run inside an authorized Linux laboratory environment with the privileges required by the configured networking functionality:

```bash
sudo ./byte
```

---

# 🧪 Development

Format the source:

```bash
go fmt ./...
```

Run the Go test suite:

```bash
go test ./...
```

Refresh module dependencies:

```bash
go mod tidy
```

Build without creating a named output:

```bash
go build .
```

---

# 🗺️ Project Design

BYTER separates the application into four major layers:

```text
┌──────────────────────────────────────────────┐
│                  CLI / CMD                   │
│        Commands • Shell • Initialization     │
└──────────────────────┬───────────────────────┘
                       │
┌──────────────────────▼───────────────────────┐
│                   CORE                       │
│   State • Counters • Rate • Scheduling       │
└───────────────┬───────────────┬───────────────┘
                │               │
        ┌───────▼───────┐ ┌────▼────────────┐
        │   L3 / L4     │ │       L7        │
        │ Packet Engine │ │ HTTP Subsystem  │
        └───────┬───────┘ └──────┬─────────┘
                │                 │
        ┌───────▼─────────────────▼────────┐
        │             UTILS                │
        │ Raw Socket • Packet Crafting     │
        │ Batch Transmission               │
        └──────────────────────────────────┘
```

This modular structure keeps the command interface separate from the underlying networking engines.

---

# 🛣️ Roadmap

Potential areas for future development:

* [ ] More comprehensive automated tests
* [ ] Improved configuration validation
* [ ] Structured logging
* [ ] Better runtime telemetry
* [ ] Configuration profiles
* [ ] Cleaner error reporting
* [ ] Improved cross-architecture handling
* [ ] Additional HTTP client profiles
* [ ] More robust resource parsing
* [ ] Benchmarking and profiling utilities
* [ ] Safer laboratory/test-environment controls

---

# ⚠️ Responsible Use

BYTER is a **security research and network experimentation tool**.

Use it only against systems and networks where you have explicit authorization.

Recommended environments include:

* Local virtual machines
* Isolated test networks
* Personal infrastructure
* Authorized penetration-testing environments
* CTF/laboratory environments
* Systems specifically designated for security testing

Do **not** use the software to disrupt, degrade, overwhelm, or interfere with third-party infrastructure.

The author and contributors are not responsible for damage, service disruption, or other consequences resulting from unauthorized use.

---

# 📜 License

BYTER is released under the **MIT License**.

See [`LICENSE`](LICENSE) for the complete license text.

```text
MIT License

Copyright (c) 2024 Byter

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

<div align="center">

## ⚡ BYTER

**Network Security Research • Linux • Go**

```text
L3  ───  L4  ───  L7
 │       │       │
 └───────┴───────┘
        CORE
         │
       BYTER
```

*Built for controlled security research and network experimentation.*

</div>
