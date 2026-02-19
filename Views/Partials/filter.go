package partials

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FilterModel struct {
	text  string
	style lipgloss.Style
}

func (m FilterModel) selectView() {
	m.style = m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func (m FilterModel) deselectView() {
	m.style = m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func InitialFilter(height int) FilterModel {
	return FilterModel{
		text: "Filter",
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderTop(true).
			Width(18).
			Height(height).
			Align(lipgloss.Center),
	}
}

func (m FilterModel) Init() tea.Cmd {
	return nil
}

func (m FilterModel) Update(msg tea.Msg) (FilterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Height(msg.Height - (9))
	}
	return m, nil
}

func (m FilterModel) View() string {
	return m.style.Render(m.text)
}
