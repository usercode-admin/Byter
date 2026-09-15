package l7

type Profile struct {
    Name            string
    UserAgent       string
    Accept          string
    AcceptLanguage  string
    AcceptEncoding  string
    Connection      string
    UpgradeInsecure string
    SecFetchDest    string
    SecFetchMode    string
    SecFetchSite    string
    SecFetchUser    string
    SecCHUA         string
    SecCHUAMobile   string
    SecCHUAPlatform string
    CacheControl    string
    Pragma          string
    TE              string
    JA3             string
}

var Profiles = map[string]Profile{
    "chrome": {
        Name:            "chrome",
        UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        Accept:          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
        AcceptLanguage:  "en-US,en;q=0.9",
        AcceptEncoding:  "gzip, deflate, br, zstd",
        Connection:      "keep-alive",
        UpgradeInsecure: "1",
        SecFetchDest:    "document",
        SecFetchMode:    "navigate",
        SecFetchSite:    "none",
        SecFetchUser:    "?1",
        SecCHUA:         `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
        SecCHUAMobile:   "?0",
        SecCHUAPlatform: `"Windows"`,
        CacheControl:    "max-age=0",
        Pragma:          "",
        TE:              "",
        JA3:             "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513,29-23-24,0",
    },
    "firefox": {
        Name:            "firefox",
        UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
        Accept:          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
        AcceptLanguage:  "en-US,en;q=0.5",
        AcceptEncoding:  "gzip, deflate, br",
        Connection:      "keep-alive",
        UpgradeInsecure: "1",
        SecFetchDest:    "document",
        SecFetchMode:    "navigate",
        SecFetchSite:    "none",
        SecFetchUser:    "?1",
        SecCHUA:         "",
        SecCHUAMobile:   "",
        SecCHUAPlatform: "",
        CacheControl:    "max-age=0",
        Pragma:          "",
        TE:              "trailers",
        JA3:             "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-34-51-43-13-45-28-65037,29-23-24-25-256-257,0",
    },
    "safari": {
        Name:            "safari",
        UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
        Accept:          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
        AcceptLanguage:  "en-US,en;q=0.9",
        AcceptEncoding:  "gzip, deflate, br",
        Connection:      "keep-alive",
        UpgradeInsecure: "1",
        SecFetchDest:    "document",
        SecFetchMode:    "navigate",
        SecFetchSite:    "none",
        SecFetchUser:    "?1",
        SecCHUA:         "",
        SecCHUAMobile:   "",
        SecCHUAPlatform: "",
        CacheControl:    "max-age=0",
        Pragma:          "",
        TE:              "",
        JA3:             "771,4865-4866-4867-49196-49195-52393-49200-49199-52392-49188-49187-49162-49161-49172-49171-157-156-61-60-53-47-49160-49170-10,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24-25,0",
    },
}

func GetProfile(name string) Profile {
    if p, ok := Profiles[name]; ok {
        return p
    }
	return Profiles["chrome"]
}