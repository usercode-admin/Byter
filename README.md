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

### Multi-Layer Network Security Research Framework

**L3 • L4 • L7 • Traffic Analysis • Automation • Performance**

<br>

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge\&logo=go\&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-Required-FCC624?style=for-the-badge\&logo=linux\&logoColor=black)
![License](https://img.shields.io/badge/License-MIT-6366F1?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Research-22C55E?style=for-the-badge)

</div>

---

## 📖 About

**BYTER** is a high-performance, multi-layer network security research and traffic-analysis framework written in **Go**.

The project is organized around three networking layers:

* **L3 — Network Layer**
* **L4 — Transport Layer**
* **L7 — Application Layer**

It combines low-level packet handling, TCP state experimentation, HTTP traffic simulation, browser-profile emulation, configurable rate control, and automated traffic scheduling into a single interactive CLI.

> **BYTER is intended for controlled security research, laboratory environments, authorized penetration testing, and educational experimentation.**

---

## ✨ Highlights

| Component                 | Description                                                                                   |
| :------------------------ | :-------------------------------------------------------------------------------------------- |
| 🌐 **L3 Engine**          | IPv4 packet construction, fragmentation research, TTL experimentation and randomized payloads |
| 🔗 **L4 Engine**          | TCP flag/state experimentation, sequence handling and TCP option configuration                |
| 🌍 **L7 Engine**          | HTTP traffic orchestration with configurable client behavior                                  |
| 🧩 **JA3 Support**        | TLS fingerprint configuration through uTLS                                                    |
| 🌐 **HTTP/2**             | HTTP/2-capable application-layer traffic                                                      |
| 🤖 **JS Solver**          | Optional Goja / Chromedp-based JavaScript challenge handling                                  |
| 🖥️ **Browser Profiles**  | Chrome, Firefox, Safari and Edge-style request profiles                                       |
| 🍪 **Session Management** | Cookie and session handling                                                                   |
| 🔄 **Proxy Pool**         | Configurable proxy management                                                                 |
| 📈 **Rate Control**       | Steady, Burst and Adaptive traffic-control modes                                              |
| 🌊 **Pulse Wave**         | Automated traffic scheduling and vector coordination                                          |
| ⚡ **Linux Performance**   | Linux-specific socket and batch-send optimizations                                            |
| 📊 **Live Statistics**    | Real-time counters and runtime state information                                              |

---

# 🏗️ Architecture

```text
                         ┌─────────────────────┐
                         │       BYTER         │
                         │   Interactive CLI   │
                         └──────────┬──────────┘
                                    │
                         ┌──────────▼──────────┐
                         │    Core Engine      │
                         │ Configuration/State │
                         └──────────┬──────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              │                     │                     │
       ┌──────▼──────┐       ┌──────▼──────┐       ┌──────▼──────┐
       │     L3      │       │     L4      │       │     L7      │
       │   Network   │       │  Transport  │       │ Application  │
       └──────┬──────┘       └──────┬──────┘       └──────┬──────┘
              │                     │                     │
              │                     │              ┌──────▼──────┐
              │                     │              │ HTTP Engine  │
              │                     │              ├──────────────┤
              │                     │              │ Headers      │
              │                     │              │ Profiles     │
              │                     │              │ Cookies      │
              │                     │              │ Proxy Pool   │
              │                     │              │ JS Solver    │
              │                     │              └──────────────┘
              │                     │
              └─────────────────────┼─────────────────────┐
                                    │                     │
                            ┌───────▼────────┐    ┌───────▼────────┐
                            │ Rate Limiter   │    │ Pulse Wave     │
                            │ Steady/Burst/  │    │ Automation     │
                            │ Adaptive       │    │ Scheduler      │
                            └────────────────┘    └────────────────┘
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
│   ├── l7_asymmetric.go
│   ├── live_counter.go
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

# 🔍 Module Overview

### `cmd/`

Responsible for the command-line interface and interactive shell.

```text
banner.go    → Startup banner and visual interface
init.go      → Initialization and environment checks
menu.go      → Help and command information
shell.go     → Interactive command processing
```

### `core/`

Contains the primary runtime engine.

```text
engine.go          → Global configuration and runtime state
counter.go         → Traffic statistics
live_counter.go    → Real-time terminal statistics
escalate.go        → Adaptive rate management
food.go            → Rate-control implementation
pulse_wave.go      → Automated scheduling
```

### `core/l7/`

Contains application-layer networking components.

```text
client.go          → HTTP client
dns.go             → DNS caching
headers.go         → Header generation
jar.go             → Cookie/session handling
profile.go         → Browser profiles
proxy.go           → Proxy management
request.go         → Request processing
solver.go          → JavaScript challenge handling
```

### `utils/`

Low-level Linux networking utilities.

```text
raw_socket.go      → Raw socket utilities
tcp_craft.go       → IPv4/TCP packet construction
batch_send.go      → Linux batch transmission
```

---

# ⚙️ Requirements

| Requirement          | Version / Notes                                     |
| :------------------- | :-------------------------------------------------- |
| **Operating System** | Linux                                               |
| **Go**               | 1.22+                                               |
| **Privileges**       | Appropriate networking capabilities may be required |
| **Architecture**     | amd64 / arm64 depending on dependencies             |
| **Dependencies**     | Managed through Go modules                          |

Install the project dependencies with:

```bash
go mod tidy
```

---

# 🚀 Installation

## 1. Clone

```bash
git clone https://github.com/yourusername/byter.git
cd byter
```

## 2. Install dependencies

```bash
go mod tidy
```

## 3. Build

```bash
go build -o byte .
```

## 4. Start

```bash
sudo ./byte
```

> Depending on the configured networking features, BYTER may require elevated privileges or specific Linux capabilities.

---

# 🎮 Interactive Shell

After starting BYTER, the application provides an interactive shell:

```text
shell >
```

### Global Commands

| Command            | Description                   |
| :----------------- | :---------------------------- |
| `set target <ip>`  | Configure the default target  |
| `set port <n>`     | Configure the target port     |
| `set rate <n>`     | Configure the traffic rate    |
| `set threads <n>`  | Configure worker count        |
| `set flood <mode>` | Select traffic-control mode   |
| `show`             | Display current configuration |
| `status`           | Display runtime statistics    |
| `stop`             | Stop the active operation     |
| `banner`           | Display the BYTER banner      |
| `help`             | Display command information   |
| `exit`             | Exit BYTER                    |
| `quit`             | Exit BYTER                    |
| `q`                | Exit BYTER                    |

---

# 🌐 Layer Modules

## L3 — Network Layer

The L3 module provides controlled IPv4 packet-generation and fragmentation research capabilities.

```text
set l3 target <ip>
set l3 rate <n>
set l3 threads <n>
set l3 flood <mode>
set l3 overlap <true|false>
set l3 fuzz-id <true|false>
```

---

## 🔗 L4 — Transport Layer

The L4 engine focuses on TCP state and packet-construction experimentation.

```text
set l4 target <ip>
set l4 port <n>
set l4 window <n>
set l4 flags <flags>
set l4 spoof-seq <true|false>
set l4 threads <n>
```

---

## 🌍 L7 — Application Layer

The L7 subsystem provides configurable HTTP client behavior and application-layer traffic research.

```text
set l7 url <url>
set l7 mode <mode>
set l7 threads <n>
set l7 conns <n>
set l7 interval <n>
set l7 browser <profile>
set l7 js-solver <solver>
set l7 ja3 <true|false>
set l7 http2 <true|false>
set l7 referer-chain <true|false>
```

Supported browser profiles:

```text
chrome
firefox
safari
edge
```

Available JavaScript solver modes:

```text
none
goja
chromedp
```

---

# 🌊 Pulse Wave

**Pulse Wave**, internally referred to as `MuWave`, provides automated scheduling for configured traffic modules.

Conceptually:

```text
        ON                         OFF
 ┌───────────────┐          ┌───────────────┐
 │   RUN CYCLE   │          │   REST CYCLE  │
 │               │          │               │
 │ configured    │          │ wait / reset  │
 │ vectors       │          │               │
 └───────┬───────┘          └───────┬───────┘
         │                          │
         └──────────────┬───────────┘
                        │
                    NEXT CYCLE
```

Configuration:

```text
set wave on <duration>
set wave off <duration>
set wave jitter <duration>
set wave vectors <vectors>
```

---

# 📊 Traffic Control

BYTER supports multiple rate-control strategies.

### Steady

Maintains a relatively consistent configured rate.

```text
steady
```

### Burst

Allows short traffic bursts within the configured control parameters.

```text
burst
```

### Adaptive

Adjusts behavior according to runtime conditions and observed errors.

```text
adaptive
```

The rate-control subsystem is implemented primarily through:

```text
core/food.go
core/escalate.go
```

---

# 🖥️ Runtime Monitoring

BYTER includes a live terminal counter for observing runtime activity.

Example conceptual output:

```text
┌─────────────────────────────────────────┐
│              BYTER STATUS               │
├─────────────────────────────────────────┤
│ State       : RUNNING                   │
│ Layer       : L7                        │
│ Threads     : 32                        │
│ Requests    : 12,481                    │
│ Success     : 12,102                    │
│ Errors      : 379                       │
│ Runtime     : 00:02:41                  │
└─────────────────────────────────────────┘
```

---

# 🧪 Research Workflow

A typical authorized laboratory workflow can be organized as:

```text
        ┌─────────────┐
        │ Define Scope│
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Prepare Lab │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Configure   │
        │   BYTER     │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Run Test    │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Monitor     │
        │ Telemetry   │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Analyze     │
        │ Results     │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │ Document    │
        └─────────────┘
```

---

# ⚡ Performance

BYTER is designed with Linux-oriented performance considerations:

* Efficient socket handling
* Batch transmission support
* Configurable worker concurrency
* Runtime counters
* Adaptive rate management
* DNS caching
* Reusable HTTP sessions
* Modular traffic engines
* Reduced unnecessary allocations where practical

The implementation keeps performance-sensitive functionality separated from the CLI layer.

---

# 🛠️ Development

Build normally with:

```bash
go build .
```

Run tests:

```bash
go test ./...
```

Format the source:

```bash
go fmt ./...
```

Check dependencies:

```bash
go mod tidy
```

---

# 🗺️ Roadmap

Possible future improvements:

* [ ] Improved terminal UI
* [ ] More detailed runtime telemetry
* [ ] Configuration profiles
* [ ] Better error reporting
* [ ] Additional HTTP client profiles
* [ ] Extended Linux performance optimizations
* [ ] Structured logging
* [ ] Exportable test reports
* [ ] Improved laboratory/test-environment safeguards
* [ ] Automated benchmark mode

---

# ⚠️ Responsible Use

BYTER is intended **only** for:

* Security research
* Authorized penetration testing
* Controlled laboratory environments
* Network-stack experimentation
* Educational purposes
* Testing systems for which you have explicit authorization

**Do not use BYTER against systems, networks, servers, or services without explicit permission.**

Unauthorized denial-of-service activity can cause outages, data loss, service disruption, and legal consequences.

The author and contributors are not responsible for misuse of this software.

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

### ⚡ BYTER

**Built for networking research.
Designed for controlled environments.**

<br>

`Go` • `Linux` • `L3` • `L4` • `L7`

</div>
