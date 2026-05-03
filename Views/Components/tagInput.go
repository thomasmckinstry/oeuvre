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
	selected   bool // Top level focus, can navigate tags, no text entry
	tagCnt     int

	titleStyle lipgloss.Style
	tagStyle   lipgloss.Style
	tagsStyle  lipgloss.Style

	width int

	errorMsg string
}

type NavMsg bool

func (m *TagInputModel) Clear() {
	m.tags = []string{}
	m.textInput.Reset()
}

// TODO: This should be a utils
func (m TagInputModel) toggleBorder() lipgloss.Style {
	if m.selected {
		return m.tagsStyle.BorderForeground(lipgloss.Color("#6E3F00"))
	}
	return m.tagsStyle.BorderForeground(lipgloss.Color("#D17600"))
}

func (m *TagInputModel) GetContents() []string {
	return m.tags
}

func InitialInput(tagCnt int, placeholder string, title string, width int, selected bool, suggestions []string) TagInputModel {
	tags := []string{}

	input := textinput.New()
	input.Placeholder = lipgloss.PlaceHorizontal(width, lipgloss.Center, placeholder)
	input.SetSuggestions(suggestions)
	input.ShowSuggestions = true
	input.SetVirtualCursor(false) // Keeps the placeholders styling consistent
	input.Blur()
	input.CharLimit = 64
	input.SetWidth(width)

	return TagInputModel{
		tags:       tags,
		textInput:  input,
		tagsCursor: 0,
		title:      title,
		tagCnt:     tagCnt,
		selected:   selected,
		width:      width,
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D17600")),
		tagStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D17600")),
		tagsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			Width(width + 3).
			BorderTop(true),
	}
}

func (m *TagInputModel) Init() tea.Cmd {
	return nil
}

func (m *TagInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.errorMsg += msg.String()
		switch msg.String() {
		case "esc": // Unfocus the component
			if m.textInput.Focused() {
				m.textInput.Blur()
			} else if m.selected {
				m.tagsStyle = m.toggleBorder()
				m.selected = false
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
		case "enter": // Add a tag from the current text input and empty the text input OR focus the component
			if !m.selected {
				m.tagsStyle = m.toggleBorder()
				m.selected = true
			} else if m.selected && !m.textInput.Focused() {
				m.textInput.Focus()
			} else if m.textInput.Value() != "" && len(m.tags) < m.tagCnt {
				m.tags = append(m.tags, m.textInput.Value())
				m.textInput.Reset()
			}
		case "j", "down": // Nav between tags
			if m.tagsCursor < len(m.tags)-1 && !m.textInput.Focused() && m.selected {
				m.tagsCursor++
			} else { // TODO: If I could come up with a way to avoid this duplication that would be great
				m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
		case "k", "up":
			if m.tagsCursor > 0 && !m.textInput.Focused() && m.selected {
				m.tagsCursor--
			} else {
				m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
		case "delete": // Delete current tag
			if len(m.tags) > 0 && !m.textInput.Focused() {
				m.tags = append(m.tags[:m.tagsCursor], m.tags[m.tagsCursor+1:]...)
				if m.tagsCursor >= len(m.tags) {
					m.tagsCursor--
				}
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
		}
	}
	return m, cmd
}

func (m *TagInputModel) View() tea.View {
	var s string
	var c = m.textInput.Cursor()
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.title)
	if m.textInput.Focused() {
		c.Y += lipgloss.Height(s)
		c.X += 1 // Aligns it correctly with the text
	}

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.textInput.View())

	for index, tag := range m.tags {
		tagStr := ""
		if tag == "" {
			continue
		}
		tagStr = " - " + tag
		if index == m.tagsCursor && !m.textInput.Focused() && m.selected { // Color selected field
			tagStr = m.tagStyle.Render(tagStr)
		}
		if index == 0 {
			tagStr = m.tagsStyle.Render(tagStr)
		}
		s = lipgloss.JoinVertical(lipgloss.Left, s, tagStr)
	}
	if len(m.tags) == 0 {
		s = lipgloss.JoinVertical(lipgloss.Left, s, m.tagsStyle.Render())
	}
	v := tea.NewView(s)
	v.Cursor = c
	return v
}

func (m TagInputModel) getInfo() (string, []string) {
	return m.title, m.tags
}
