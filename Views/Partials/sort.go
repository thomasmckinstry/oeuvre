package partials

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SortModel struct {
	options       []string
	optionsCursor int
	selected      bool
	mainStyle     lipgloss.Style
	contentStyle  lipgloss.Style
	textStyle     lipgloss.Style
}

func (m SortModel) toggleBorder() lipgloss.Style {
	if m.selected == true {
		return m.mainStyle.BorderForeground(lipgloss.Color("#6E3F00"))
	}
	return m.mainStyle.BorderForeground(lipgloss.Color("#D17600"))
}

func (m SortModel) toggleText() lipgloss.Style {
	if m.selected == true {
		return m.textStyle.Foreground(lipgloss.Color("#6E3F00"))
	}
	return m.textStyle.Foreground(lipgloss.Color("#D17600"))
}

func InitialSort(height int) SortModel {
	return SortModel{
		options:       []string{"title", "release date", "genre", "theme", "medium"},
		optionsCursor: 0,
		selected:      false,
		mainStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderTop(true).
			Width(18).
			Height(height).
			Align(lipgloss.Center),
		contentStyle: lipgloss.NewStyle().
			MarginTop(1),
		textStyle: lipgloss.NewStyle(),
	}
}

func (m *SortModel) Init() tea.Cmd {
	return nil
}

func (m *SortModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "L", "H", "J", "K":
			m.mainStyle = m.toggleBorder()
			m.textStyle = m.toggleText()
			m.selected = !m.selected
		case "l", "right":
			m.optionsCursor = (m.optionsCursor + 1) % len(m.options)
		case "h", "left":
			m.optionsCursor = (m.optionsCursor - 1) % len(m.options)
			if m.optionsCursor < 0 {
				m.optionsCursor = len(m.options) - 1
			}
		}

	}
	return m, nil
}

func (m *SortModel) View() tea.View {
	header := lipgloss.PlaceHorizontal(18, lipgloss.Center, "Sort")
	contents := lipgloss.PlaceHorizontal(18, lipgloss.Center,
		lipgloss.PlaceVertical(3, lipgloss.Center,
			m.contentStyle.Render(
				lipgloss.JoinHorizontal(lipgloss.Center, "< ",
					lipgloss.PlaceHorizontal(14, lipgloss.Center,
						m.textStyle.Render(m.options[m.optionsCursor])), " >"),
			),
		),
	)

	return tea.NewView(m.mainStyle.Render(lipgloss.JoinVertical(lipgloss.Right, header, contents)))
}
