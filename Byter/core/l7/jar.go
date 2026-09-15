package l7

import (
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "sync"
)

type Session struct {
    Jar     http.CookieJar
    Profile Profile
    Referer string
    Tokens  map[string]string
    mu      sync.Mutex
}

func NewSession(profileName string) (*Session, error) {
    jar, err := cookiejar.New(nil)
    if err != nil {
        return nil, err
    }
    return &Session{
        Jar:     jar,
        Profile: GetProfile(profileName),
        Tokens:  make(map[string]string),
    }, nil
}

func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
    s.Jar.SetCookies(u, cookies)
}

func (s *Session) Cookies(u *url.URL) []*http.Cookie {
    return s.Jar.Cookies(u)
}

func (s *Session) SaveToken(k, v string) {
    s.mu.Lock()
    s.Tokens[k] = v
    s.mu.Unlock()
}

func (s *Session) GetToken(k string) string {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.Tokens[k]
}