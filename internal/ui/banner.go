package ui

import "fmt"

const rawBanner = `
 ███████╗ █████╗ ██╗  ████████╗ █████╗ ███████╗██╗
 ██╔════╝██╔══██╗██║  ╚══██╔══╝██╔══██╗██╔════╝██║
 █████╗  ███████║██║     ██║   ███████║███████╗██║
 ██╔══╝  ██╔══██║██║     ██║   ██╔══██║╚════██║██║
 ██║     ██║  ██║███████╗██║   ██║  ██║███████║██║
 ╚═╝     ╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝
`

func PrintBanner() {
	fmt.Println(TitleStyle.Render(rawBanner))
	fmt.Println(SubtleStyle.Render("           P O S   C O N T R O L   C E N T R E"))
	fmt.Println(SubtleStyle.Render("──────────────────────────────────────────────────"))
}
