//go:build linux

package utils

import (
    "fmt"
    "syscall"
)

func NewRawSocket() (int, error) {
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
    if err != nil {
        return -1, fmt.Errorf("socket: %w", err)
    }
    if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
        syscall.Close(fd)
        return -1, fmt.Errorf("setsockopt: %w", err)
    }
    return fd, nil
}
