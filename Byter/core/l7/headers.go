package l7

import "net/http"

func ApplyProfile(req *http.Request, p Profile) {
    h := req.Header

    h.Set("User-Agent", p.UserAgent)
    h.Set("Accept", p.Accept)
    h.Set("Accept-Language", p.AcceptLanguage)
    h.Set("Accept-Encoding", p.AcceptEncoding)
    h.Set("Connection", p.Connection)
    h.Set("Upgrade-Insecure-Requests", p.UpgradeInsecure)
    h.Set("Sec-Fetch-Dest", p.SecFetchDest)
    h.Set("Sec-Fetch-Mode", p.SecFetchMode)
    h.Set("Sec-Fetch-Site", p.SecFetchSite)
    h.Set("Sec-Fetch-User", p.SecFetchUser)

    if p.SecCHUA != "" {
        h.Set("sec-ch-ua", p.SecCHUA)
        h.Set("sec-ch-ua-mobile", p.SecCHUAMobile)
        h.Set("sec-ch-ua-platform", p.SecCHUAPlatform)
    }
    if p.CacheControl != "" {
        h.Set("Cache-Control", p.CacheControl)
    }
    if p.Pragma != "" {
        h.Set("Pragma", p.Pragma)
    }
    if p.TE != "" {
        h.Set("TE", p.TE)
    }
}

func ApplySubResource(req *http.Request, p Profile, referer string, dest string) {
    ApplyProfile(req, p)
    req.Header.Set("Sec-Fetch-Dest", dest)
    req.Header.Set("Sec-Fetch-Mode", "no-cors")
    req.Header.Set("Sec-Fetch-Site", "same-origin")
    if referer != "" {
        req.Header.Set("Referer", referer)
    }
    req.Header.Del("Upgrade-Insecure-Requests")
}