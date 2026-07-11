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

// PrintBanner renders the FALTASI wordmark plus the "control centre"
// tagline, in brand green.
func PrintBanner() {
	fmt.Println(TitleStyle.Render(rawBanner))
	fmt.Println(SubtleStyle.Render("           P O S   C O N T R O L   C E N T R E"))
	fmt.Println(SubtleStyle.Render("──────────────────────────────────────────────────"))
}
