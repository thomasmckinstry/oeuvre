package partials

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddModel struct {
	text  string
	style lipgloss.Style
}

func (m AddModel) selectView() {
	m.style = m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func (m AddModel) deselectView() {
	m.style = m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func InitialAdd() AddModel {
	return AddModel{
		text: "Add",
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(lipgloss.Color("#6E3F00")).
			Width(18).
			Height(1).
			//Background(lipgloss.Color("#F3EAC7")).
			//BorderBackground(lipgloss.Color("#F3EAC7")).
			Align(lipgloss.Center),
	}
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (AddModel, tea.Cmd) {
	return m, nil
}

func (m AddModel) View() string {
	return m.style.Render(m.text)
}
