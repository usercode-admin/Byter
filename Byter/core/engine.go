package core

import (
	"net"
	"time"
)

type FloodMode int

const (
	FloodSteady FloodMode = iota
	FloodBurst
	FloodAdaptive
)

func (f FloodMode) String() string {
	switch f {
	case FloodBurst:
		return "burst"
	case FloodAdaptive:
		return "adaptive"
	default:
		return "steady"
	}
}

type GlobalConfig struct {
	Target  string
	Port    int
	Rate    int
	Threads int
	Flood   FloodMode
}

type L3Config struct {
	Target        string
	Rate          int
	Threads       int
	Flood         FloodMode
	OffsetOverlap bool
	FuzzID        bool
	BatchSize     int
}

type L4Config struct {
	Target    string
	Port      int
	Window    int
	Flags     []string
	SpoofSeq  bool
	Rate      int
	Threads   int
	Flood     FloodMode
	BatchSize int
}

type L7Config struct {
	URL           string
	Threads       int
	ConnPerThread int
	WinShrink     int
	Interval      int
	Browser       string
	JSSolver      string
	JA3Spoof      bool
	HTTP2         bool
	RefererChain  bool
	Mode          string
	SlowHeaderPct int
	SlowPostPct   int
	ProxyList     []string
	Escalate      bool
}

type WaveConfig struct {
	PulseOn  time.Duration
	PulseOff time.Duration
	Jitter   time.Duration
	Vectors  []string
}

type EngineState struct {
	RawFD          int
	ResolvedTarget net.IP
	PrebuiltIPv4   []byte

	Global GlobalConfig
	L3     L3Config
	L4     L4Config
	L7     L7Config
	Wave   WaveConfig

	Counters *Counters

	StopCurrent chan struct{}
	CurrentMode string
}

func NewEngineState() *EngineState {
	return &EngineState{
		Global: GlobalConfig{
			Target:  "127.0.0.1",
			Port:    80,
			Rate:    1000,
			Threads: 8,
			Flood:   FloodSteady,
		},
		L3: L3Config{
			Rate:          1000,
			Threads:       8,
			Flood:         FloodSteady,
			OffsetOverlap: true,
			FuzzID:        true,
			BatchSize:     256,
		},
		L4: L4Config{
			Port:      80,
			Window:    65535,
			Flags:     []string{"ACK", "FIN"},
			SpoofSeq:  true,
			Rate:      1000,
			Threads:   8,
			Flood:     FloodSteady,
			BatchSize: 256,
		},
		L7: L7Config{
			Threads:       8,
			ConnPerThread: 50,
			WinShrink:     10,
			Interval:      5,
			Browser:       "chrome",
			JSSolver:      "goja",
			JA3Spoof:      true,
			HTTP2:         true,
			RefererChain:  true,
			Mode:          "hybrid",
			SlowHeaderPct: 10,
			SlowPostPct:   10,
			Escalate:      true,
		},
		Wave: WaveConfig{
			PulseOn:  60 * time.Second,
			PulseOff: 5 * time.Second,
			Jitter:   3 * time.Second,
			Vectors:  []string{"l7", "l4", "l3"},
		},
		Counters: NewCounters(),
	}
}

func (e *EngineState) StopRunning() {
	if e.StopCurrent != nil {
		close(e.StopCurrent)
		e.StopCurrent = nil
		e.CurrentMode = ""
	}
}

func (e *EngineState) IsRunning() bool { return e.StopCurrent != nil }