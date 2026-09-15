//go:build linux

package utils

import (
	"runtime"
	"syscall"
	"unsafe"
)

func sysSendmmsgNum() uintptr {
	switch runtime.GOARCH {
	case "amd64":
		return 307
	case "arm64":
		return 269
	case "386":
		return 345
	case "arm":
		return 374
	default:
		return 0
	}
}

type mmsghdr struct {
	Hdr syscall.Msghdr
	Len uint32
	_   [4]byte
}

func SendBatch(fd int, pkts [][]byte, addr *syscall.SockaddrInet4) int {
	n := len(pkts)
	if n == 0 {
		return 0
	}

	sysno := sysSendmmsgNum()
	if sysno == 0 {
		return sendBatchFallback(fd, pkts, addr)
	}

	hdrs := make([]mmsghdr, n)
	iovs := make([]syscall.Iovec, n)

	for i := 0; i < n; i++ {
		if len(pkts[i]) == 0 {
			continue
		}
		iovs[i].Base = &pkts[i][0]
		iovs[i].SetLen(len(pkts[i]))
		hdrs[i].Hdr.Iov = &iovs[i]
		hdrs[i].Hdr.Iovlen = 1
		hdrs[i].Hdr.Name = (*byte)(unsafe.Pointer(addr))
		hdrs[i].Hdr.Namelen = uint32(syscall.SizeofSockaddrInet4)
	}

	r, _, errno := syscall.Syscall6(
		sysno,
		uintptr(fd),
		uintptr(unsafe.Pointer(&hdrs[0])),
		uintptr(n),
		0, 0, 0,
	)
	if errno != 0 {
		return sendBatchFallback(fd, pkts, addr)
	}
	return int(r)
}

func sendBatchFallback(fd int, pkts [][]byte, addr *syscall.SockaddrInet4) int {
	n := 0
	for _, p := range pkts {
		if len(p) == 0 {
			continue
		}
		if err := syscall.Sendto(fd, p, 0, addr); err == nil {
			n++
		}
	}
	return n
}