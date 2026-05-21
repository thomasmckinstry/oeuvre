package components

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
)

type arrowInputKeyMap struct {
	Nav   key.Binding
	Left  key.Binding
	Right key.Binding
}

var defaultArrowMap = arrowInputKeyMap{
	Nav:   key.NewBinding(key.WithKeys("H", "J", "K", "L")),
	Left:  key.NewBinding(key.WithKeys("h", "left")),
	Right: key.NewBinding(key.WithKeys("l", "right")),
}

type ArrowModel struct {
	Options       []string
	OptionsCursor int
	selected      bool
	title         string
	width         int
	mainStyle     lipgloss.Style
	contentStyle  lipgloss.Style
	textStyle     lipgloss.Style
}

func (m *ArrowModel) GetContents() string {
	return m.Options[m.OptionsCursor]
}

func InitialArrow(options []string, title string, width int, height int) ArrowModel {
	return ArrowModel{
		Options:       options,
		OptionsCursor: 0,
		selected:      false,
		mainStyle: lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center),
		contentStyle: lipgloss.NewStyle().
			MarginTop(1),
		textStyle: lipgloss.NewStyle().
			Foreground(Unfocused),
		title: title,
	}
}

func (m *ArrowModel) Init() tea.Cmd {
	return nil
}

func (m *ArrowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultArrowMap.Nav):
			m.selected = !m.selected
		case key.Matches(msg, defaultArrowMap.Right):
			m.OptionsCursor = (m.OptionsCursor + 1) % len(m.Options)
		case key.Matches(msg, defaultArrowMap.Left):
			m.OptionsCursor = (m.OptionsCursor - 1) % len(m.Options)
			if m.OptionsCursor < 0 {
				m.OptionsCursor = len(m.Options) - 1
			}
		}
	}
	cmd = func() tea.Msg { return NavMsg(true) }
	cmds = tea.Batch(cmds, cmd)
	return m, cmds
}

func (m *ArrowModel) View() tea.View {
	header := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.title)
	options := m.Options[m.OptionsCursor]
	if m.selected {
		options = m.textStyle.Foreground(lipgloss.Color("#D17600")).Render(options)
	} else {
		options = m.textStyle.Render(options)
	}
	contents := lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		m.contentStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Center, "< ",
				lipgloss.PlaceHorizontal(m.width-4, lipgloss.Center,
					options), " >"),
		),
	)

	return tea.NewView(m.mainStyle.Render(lipgloss.JoinVertical(lipgloss.Center, header, contents)))
}
