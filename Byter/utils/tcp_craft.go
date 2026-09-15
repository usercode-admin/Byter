package utils

import (
	"encoding/binary"
	"math/rand"
	"net"
)

type IPv4Header struct {
	VersionIHL uint8
	TOS        uint8
	TotalLen   uint16
	ID         uint16
	FlagsFrag  uint16
	TTL        uint8
	Protocol   uint8
	Checksum   uint16
	Src        [4]byte
	Dst        [4]byte
}

type TCPHeader struct {
	SrcPort  uint16
	DstPort  uint16
	SeqNum   uint32
	AckNum   uint32
	DataOff  uint8
	Flags    uint8
	Window   uint16
	Checksum uint16
	Urgent   uint16
}

const (
	TCP_FIN = 0x01
	TCP_SYN = 0x02
	TCP_RST = 0x04
	TCP_PSH = 0x08
	TCP_ACK = 0x10
	TCP_URG = 0x20
	TCP_ECE = 0x40
	TCP_CWR = 0x80
)

func FlagsToByte(flags []string) uint8 {
	var b uint8
	for _, f := range flags {
		switch f {
		case "FIN":
			b |= TCP_FIN
		case "SYN":
			b |= TCP_SYN
		case "RST":
			b |= TCP_RST
		case "PSH":
			b |= TCP_PSH
		case "ACK":
			b |= TCP_ACK
		case "URG":
			b |= TCP_URG
		case "ECE":
			b |= TCP_ECE
		case "CWR":
			b |= TCP_CWR
		}
	}
	return b
}

func RandomIPID() uint16 { return uint16(rand.Uint32()) }
func SpoofSeq() uint32   { return rand.Uint32() }

func RandomSrcPort() uint16 {
	return uint16(1024 + rand.Intn(65535-1024))
}

func IPChecksum(b []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(b); i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	if len(b)%2 == 1 {
		sum += uint32(b[len(b)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func BuildIPv4(h IPv4Header) []byte {
	b := make([]byte, 20)
	b[0] = (h.VersionIHL << 4) | 5
	b[1] = h.TOS
	binary.BigEndian.PutUint16(b[2:], h.TotalLen)
	binary.BigEndian.PutUint16(b[4:], h.ID)
	binary.BigEndian.PutUint16(b[6:], h.FlagsFrag)
	b[8] = h.TTL
	b[9] = h.Protocol
	copy(b[12:16], h.Src[:])
	copy(b[16:20], h.Dst[:])
	binary.BigEndian.PutUint16(b[10:], checksum(b))
	return b
}

func BuildTCPWithChecksum(h TCPHeader, src, dst [4]byte) []byte {
	b := make([]byte, 20)
	binary.BigEndian.PutUint16(b[0:], h.SrcPort)
	binary.BigEndian.PutUint16(b[2:], h.DstPort)
	binary.BigEndian.PutUint32(b[4:], h.SeqNum)
	binary.BigEndian.PutUint32(b[8:], h.AckNum)
	b[12] = (h.DataOff << 4) | 0
	b[13] = h.Flags
	binary.BigEndian.PutUint16(b[14:], h.Window)
	binary.BigEndian.PutUint16(b[18:], h.Urgent)

	pseudo := make([]byte, 12)
	copy(pseudo[0:4], src[:])
	copy(pseudo[4:8], dst[:])
	pseudo[8] = 0
	pseudo[9] = 6
	binary.BigEndian.PutUint16(pseudo[10:], uint16(len(b)))
	all := append(pseudo, b...)
	binary.BigEndian.PutUint16(b[16:], checksum(all))
	return b
}

func checksum(b []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(b); i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	if len(b)%2 == 1 {
		sum += uint32(b[len(b)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func AddrTo4(ip net.IP) [4]byte {
	var out [4]byte
	if v4 := ip.To4(); v4 != nil {
		copy(out[:], v4)
	}
	return out
}