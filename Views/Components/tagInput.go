package components

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TagInputModel struct {
	textInput  textinput.Model
	tags       []string
	tagsCursor int
	title      string
	selected   bool

	tagStyle lipgloss.Style

	errorMsg string
}

func InitialInput(tagCnt int, placeholder string, title string, width int, selected bool) TagInputModel {
	tags := []string{}

	input := textinput.New()
	input.Placeholder = placeholder
	input.SetVirtualCursor(false)
	input.Blur()
	input.CharLimit = 64
	input.SetWidth(width)

	return TagInputModel{
		tags:       tags,
		textInput:  input,
		tagsCursor: 0,
		title:      title,
		selected:   selected,
		tagStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D17600")),
	}
}

func (m TagInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TagInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.errorMsg = msg.String()
		switch msg.String() {
		case "esc": // Unfocus the component
			m.textInput.Blur()
		case "j", "down": // Nav between tags
			if m.selected && !m.textInput.Focused() {
				m.selected = false
			} else if m.tagsCursor < len(m.tags)-1 && m.selected {
				m.tagsCursor++
			} else {
				m.selected = true
			}
		case "k", "up":
			if m.selected && !m.textInput.Focused() {
				m.selected = false
			} else if m.tagsCursor > 0 && m.selected {
				m.tagsCursor--
			} else {
				m.selected = true
			}
		case "delete": // Delete current tag
			if len(m.tags) > 0 {
				m.tags = append(m.tags[:m.tagsCursor], m.tags[m.tagsCursor+1:]...)
				if m.tagsCursor > len(m.tags) {
					m.tagsCursor--
				}
			}
		case "enter": // Add a tag from the current text input and empty the text input OR focus the component
			if !m.textInput.Focused() {
				m.textInput.Focus()
			} else if m.textInput.Value() != "" && len(m.tags) < 3 {
				m.tags = append(m.tags, m.textInput.Value())
				m.textInput.Reset()
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
	return m, cmd
}

func (m TagInputModel) View() tea.View {
	var s string
	s = m.title
	if m.selected {
		s = m.tagStyle.Render(s) // Get an independent "selected" style for showing color
	}

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.textInput.View())

	for index, tag := range m.tags {
		tagStr := ""
		if tag == "" {
			continue
		}
		tagStr = " - " + tag
		if index == m.tagsCursor && m.textInput.Focused() { // Color selected field
			tagStr = m.tagStyle.Render(tagStr)
		}
		s = lipgloss.JoinVertical(lipgloss.Left, s, tagStr)
	}

	return tea.NewView(s)
}
