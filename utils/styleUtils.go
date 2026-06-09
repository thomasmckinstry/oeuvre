package utils

import (
	"charm.land/lipgloss/v2"
)

var (
	Unfocused = lipgloss.Color("#6E3F00")
	Focused   = lipgloss.Color("#D17F00")
)

var (
	ButtonStyle lipgloss.Style = lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(Unfocused)
)

func RenderFocused(style lipgloss.Style, content string, isFocused bool) string {
	if isFocused {
		return style.BorderForeground(Focused).Render(content)
	} else {
		return style.BorderForeground(Unfocused).Render(content)
	}
}

func SetTheme(focus, unfocus string) {
	Unfocused = lipgloss.Color(unfocus)
	Focused = lipgloss.Color(focus)
}
