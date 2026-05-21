package utils

import (
	"charm.land/lipgloss/v2"
)

var (
	Unfocused = lipgloss.Color("#6E3F00")
	Focused   = lipgloss.Color("#D17600")
)

func RenderFocused(style lipgloss.Style, content string, isFocused bool) string {
	if isFocused {
		return style.BorderForeground(Focused).Render(content)
	} else {
		return style.Render(content)
	}
}
