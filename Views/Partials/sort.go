package partials

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SortModel struct {
	text  string
	style lipgloss.Style
}

func (m SortModel) selectView() {
	m.style = m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func (m SortModel) deselectView() {
	m.style = m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func InitialSort(height int) SortModel {
	return SortModel{
		text: "Sort",
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderTop(true).
			Width(18).
			Height(height).
			Align(lipgloss.Center),
	}
}

func (m SortModel) Init() tea.Cmd {
	return nil
}

func (m SortModel) Update(msg tea.Msg) (SortModel, tea.Cmd) {
	return m, nil
}

func (m SortModel) View() string {
	return m.style.Render(m.text)
}
