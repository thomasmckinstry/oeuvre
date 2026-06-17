package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/ouevre/utils"
)

type textKeyMap struct {
	Confirm key.Binding
	Unfocus key.Binding
}

var defaultTextMap = textKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Focus the text input"),
	),
	Unfocus: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Unfocus the text input"),
	),
}

type TextInputModel struct {
	Textinput textinput.Model
	title     string
	width     int
}

func (m *TextInputModel) GetContents() string {
	return m.Textinput.Value()
}

func (m *TextInputModel) Clear() {
	m.Textinput.Reset()
}

func InitialTextInput(width int, title string, placeholder string, suggestions []string) TextInputModel {
	input := textinput.New()
	input.Placeholder = lipgloss.PlaceHorizontal(width, lipgloss.Center, placeholder)
	input.SetSuggestions(suggestions)
	input.ShowSuggestions = true
	input.SetVirtualCursor(false) // Keeps the placeholders styling consistent
	input.Blur()
	input.CharLimit = 64
	input.SetWidth(width)

	return TextInputModel{
		Textinput: input,
		title:     title,
		width:     width,
	}
}

func (m *TextInputModel) Init() tea.Cmd {
	return nil
}

func (m *TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultTextMap.Confirm):
			if m.Textinput.Focused() {
				m.Textinput.Blur()
				break
			}
			m.Textinput.Focus()
		case key.Matches(msg, defaultTextMap.Unfocus):
			m.Textinput.Blur()
		}
	}
	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = tea.Batch(cmds, cmd)
	cmd = func() tea.Msg { return utils.NavMsg(!m.Textinput.Focused()) }
	cmds = tea.Batch(cmds, cmd)
	return m, cmds
}

func (m *TextInputModel) View() tea.View {
	var s string
	var c = m.Textinput.Cursor()
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.title)
	if m.Textinput.Focused() {
		c.Y += lipgloss.Height(s)
		//c.X += 1 // Aligns it correctly with the text
	}

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.Textinput.View())
	v := tea.NewView(s)
	v.Cursor = c
	return v
}
