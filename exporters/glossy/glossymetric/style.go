package glossymetric

import (
	"os"

	"charm.land/lipgloss/v2"
)

func isNoColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

var (
	styleHeader     = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	styleMetricName = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	styleLabel      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleValue      = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleUnit       = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	styleAttrKey    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	styleAttrVal    = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	styleSeparator  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func renderStyled(style lipgloss.Style, s string) string {
	if isNoColor() {
		return s
	}

	return style.Render(s)
}
