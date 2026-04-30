package components

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TextInputModel struct {
	textinput textinput.Model
	title     string
	width     int
}

func (m *TextInputModel) GetContents() string {
	return m.textinput.Value()
}

func (m *TextInputModel) Clear() {
	m.textinput.Reset()
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
		textinput: input,
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
		switch msg.String() {
		case "enter":
			if m.textinput.Focused() {
				m.textinput.Blur()
				break
			}
			m.textinput.Focus()
		case "esc":
			m.textinput.Blur()
		}
	}
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = tea.Batch(cmds, cmd)
	cmd = func() tea.Msg { return NavMsg(!m.textinput.Focused()) }
	cmds = tea.Batch(cmds, cmd)
	return m, cmds
}

func (m *TextInputModel) View() tea.View {
	var s string
	var c = m.textinput.Cursor()
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.title)
	if m.textinput.Focused() {
		c.Y += lipgloss.Height(s)
		c.X += 1 // Aligns it correctly with the text
	}

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.textinput.View())
	v := tea.NewView(s)
	v.Cursor = c
	return v
}
