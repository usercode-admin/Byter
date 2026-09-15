package l7

import (
    "context"
    "fmt"
    "regexp"
    "strings"
    "time"

    "github.com/chromedp/chromedp"
    "github.com/dop251/goja"
)

func SolveJS(engine, html, pageURL string, s *Session) (string, error) {
    switch engine {
    case "none", "":
        return "", nil
    case "goja":
        return solveGoja(html)
    case "chromedp":
        return solveChromedp(pageURL, s)
    default:
        return "", fmt.Errorf("js-solver is not supported: %s", engine)
    }
}

var scriptTagRe = regexp.MustCompile(`(?is)<script[^>]*>(.*?)</script>`)

func solveGoja(html string) (string, error) {
    matches := scriptTagRe.FindAllStringSubmatch(html, -1)
    if len(matches) == 0 {
        return "", nil
    }
    vm := goja.New()
    var last goja.Value
    for _, m := range matches {
        code := strings.TrimSpace(m[1])
        if code == "" || strings.HasPrefix(code, "src=") {
            continue
        }
        v, err := vm.RunString(code)
        if err == nil {
            last = v
        }
    }
    if last == nil {
        return "", nil
    }
    return last.String(), nil
}

func solveChromedp(pageURL string, s *Session) (string, error) {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    var finalHTML string

    err := chromedp.Run(ctx,
        chromedp.Navigate(pageURL),
        chromedp.Sleep(3*time.Second),
        chromedp.OuterHTML("html", &finalHTML),
    )
    if err != nil {
        return "", err
    }
    return finalHTML, nil
}