package l7

import (
    "context"
    "crypto/tls"
    "net"
    "net/http"
    "time"

    utls "github.com/refraction-networking/utls"
    "golang.org/x/net/http2"
)

func NewClient(s *Session, ja3Spoof, http2On bool, timeout time.Duration, hostname string) *http.Client {
    dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}

    tr := &http.Transport{
        DialContext:         dialer.DialContext,
        MaxIdleConns:        2,
        MaxIdleConnsPerHost: 2,
        IdleConnTimeout:     90 * time.Second,
        DisableKeepAlives:   false,
        ForceAttemptHTTP2:   http2On,
    }

    if hostname == "" {
        hostname = "localhost"
    }

    if ja3Spoof {
        tr.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
            raw, err := dialer.DialContext(ctx, network, addr)
            if err != nil {
                return nil, err
            }
            cfg := &utls.Config{ServerName: hostname, InsecureSkipVerify: true}
            uconn := utls.UClient(raw, cfg, utls.HelloChrome_Auto)
            if err := uconn.HandshakeContext(ctx); err != nil {
                raw.Close()
                return nil, err
            }
            return uconn, nil
        }
    } else {
        tr.TLSClientConfig = &tls.Config{
            ServerName:         hostname,
            InsecureSkipVerify: true,
        }
    }

    if http2On {
        _ = http2.ConfigureTransport(tr)
    }

    return &http.Client{
        Transport: tr,
        Jar:       s.Jar,
        Timeout:   timeout,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return nil
        },
    }
}
