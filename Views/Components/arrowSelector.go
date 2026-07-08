package components

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	. "github.com/thomasmckinstry/ouevre/utils"
)

type arrowInputKeyMap struct {
	Nav   key.Binding
	Left  key.Binding
	Right key.Binding
}

var defaultArrowMap = arrowInputKeyMap{
	Nav: key.NewBinding(key.WithKeys("ctrl+h", "ctrl+j", "ctrl+k", "ctrl+l"),
		key.WithHelp("ctrl+h/j/k/l", "Navigate in/out of arrow input.")),
	Left: key.NewBinding(key.WithKeys("h", "left"),
		key.WithHelp("h/left", "Navigate left in arrow input.")),
	Right: key.NewBinding(key.WithKeys("l", "right"),
		key.WithHelp("l/right", "Navigate right in arrow input.")),
}

var (
	arrowContentStyle lipgloss.Style = lipgloss.NewStyle().
		MarginTop(1)
)

type ArrowModel struct {
	Options       []string
	OptionsCursor int
	title         string
	width         int
}

func (m *ArrowModel) GetContents() int {
	return m.OptionsCursor
}

func InitialArrow(options []string, title string, width int, height int) ArrowModel {
	return ArrowModel{
		Options:       options,
		OptionsCursor: 0,
		title:         title,
		width:         width,
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
	contents := lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		arrowContentStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Center, "< ",
				lipgloss.PlaceHorizontal(m.width-4, lipgloss.Center,
					options), " >"),
		),
	)

	return tea.NewView(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, header, contents)))
}
