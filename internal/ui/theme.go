package ui

import "github.com/charmbracelet/lipgloss"

// Brand colors — matches the FALTASI/BALCE green used across the desktop app.
var (
	ColorPrimary = lipgloss.Color("#7bc83a")
	ColorMuted   = lipgloss.Color("241")
	ColorError   = lipgloss.Color("#e5484d")
	ColorSuccess = lipgloss.Color("#7bc83a")
)

var (
	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)

	SubtleStyle = lipgloss.NewStyle().Foreground(ColorMuted)

	ErrorStyle = lipgloss.NewStyle().Foreground(ColorError).Bold(true)

	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 3)

	BadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(ColorPrimary).
			Bold(true).
			Padding(0, 1)
)
