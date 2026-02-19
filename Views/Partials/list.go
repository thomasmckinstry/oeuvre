package partials

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	Style lipgloss.Style // TODO: This probably shouldn't be public
	table table.Model
}

func InitialList(width int, height int, cols []table.Column, rows []table.Row) ListModel {
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return ListModel{
		Style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("#6E3F00")).
			PaddingTop(1).
			Width(width).
			Height(height),
		table: t,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Style = m.Style.Height(msg.Height).Width(msg.Width - (20 + 2))
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "J":
			m.Style = m.Style.BorderForeground(lipgloss.Color("#D17600"))
		case "K":
			m.Style = m.Style.BorderForeground(lipgloss.Color("#6E3F00"))
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return m.Style.Render(m.table.View())
}
