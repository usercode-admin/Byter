package l7

import (
    "io"
    "net/http"
    "net/url"
    "regexp"
    "time"
)

var (
    linkRe   = regexp.MustCompile(`(?is)<link[^>]+href=["']([^"']+)["']`)
    scriptRe = regexp.MustCompile(`(?is)<script[^>]+src=["']([^"']+)["']`)
    imgRe    = regexp.MustCompile(`(?is)<img[^>]+src=["']([^"']+)["']`)
)

func FetchFlow(client *http.Client, s *Session, targetURL string, winShrink, interval int, refererChain bool, jsSolver string) error {
    docReq, err := http.NewRequest("GET", targetURL, nil)
    if err != nil {
        return err
    }
    ApplyProfile(docReq, s.Profile)

    resp, err := client.Do(docReq)
    if err != nil {
        return err
    }
    body, err := io.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        return err
    }

    if jsSolver != "none" && jsSolver != "" {
        _, _ = SolveJS(jsSolver, string(body), targetURL, s)
    }

    if refererChain {
        base, _ := url.Parse(targetURL)
        for _, u := range extractResources(string(body), base) {
            go fetchResource(client, s, u, targetURL, winShrink, interval)
        }
    }

    slowReadBody(resp, winShrink, interval)
    return nil
}

func extractResources(html string, base *url.URL) []string {
    var out []string
    seen := map[string]bool{}
    for _, re := range []*regexp.Regexp{linkRe, scriptRe, imgRe} {
        for _, m := range re.FindAllStringSubmatch(html, -1) {
            u, err := base.Parse(m[1])
            if err != nil {
                continue
            }
            str := u.String()
            if !seen[str] {
                seen[str] = true
                out = append(out, str)
            }
        }
    }
    return out
}

func fetchResource(client *http.Client, s *Session, u, referer string, winShrink, interval int) {
    req, err := http.NewRequest("GET", u, nil)
    if err != nil {
        return
    }
    ApplySubResource(req, s.Profile, referer, destFor(u))

    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    slowReadBody(resp, winShrink, interval)
}

var (
    cssExt = regexp.MustCompile(`\.css($|\?)`)
    jsExt  = regexp.MustCompile(`\.js($|\?)`)
    imgExt = regexp.MustCompile(`\.(png|jpg|jpeg|gif|webp|avif)($|\?)`)
)

func destFor(u string) string {
    switch {
    case cssExt.MatchString(u):
        return "style"
    case jsExt.MatchString(u):
        return "script"
    case imgExt.MatchString(u):
        return "image"
    }
    return "empty"
}

func slowReadBody(resp *http.Response, winShrink, interval int) {
    if resp == nil || resp.Body == nil {
        return
    }
    size := winShrink
    if size < 1 {
        size = 1
    }
    buf := make([]byte, size)
    for {
        _, err := resp.Body.Read(buf)
        if err != nil {
            return
        }
        time.Sleep(time.Duration(interval) * time.Second)
    }
}