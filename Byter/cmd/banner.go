package cmd

import "fmt"

const bannerArt = `
      d8888b. db    db d888888b d88888b d8888b.
      88  ` + "`" + `8D ` + "`" + `8b  d8' ` + "`" + `~~88~~' 88'     88  ` + "`" + `8D
      88oooY'  ` + "`" + `8bd8'     88    88ooooo 88oobY'
      88~~~b.    88       88    88~~~~~ 88` + "`" + `8b
      88   8D    88       88    88.     88 ` + "`" + `88.
      Y8888P'    YP       YP    Y88888P 88   YD
`

func PrintBanner() {
    fmt.Print(bannerArt)
    fmt.Println("      [+] hint: Use ./byte -h, --h for support")
    fmt.Println("    [1] layer3")
    fmt.Println("    [2] layer4")
    fmt.Println("    [3] layer7")
    fmt.Println("    [4] muwave")
}
