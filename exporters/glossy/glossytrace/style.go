package glossytrace

import (
	"os"
	"time"

	"charm.land/lipgloss/v2"
)

func isNoColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

var (
	styleHeader    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleSpanName  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	styleConnector = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleAttrKey   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	styleAttrVal   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	styleEventName = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	styleTimestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleFlameBar  = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))

	styleLinkPrefix = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	styleLinkRef    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	styleDurationFast = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleDurationWarn = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleDurationSlow = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func durationStyle(d, warn, critical time.Duration) lipgloss.Style {
	switch {
	case d >= critical:
		return styleDurationSlow
	case d >= warn:
		return styleDurationWarn
	default:
		return styleDurationFast
	}
}

func renderStyled(style lipgloss.Style, s string) string {
	if isNoColor() {
		return s
	}

	return style.Render(s)
}
