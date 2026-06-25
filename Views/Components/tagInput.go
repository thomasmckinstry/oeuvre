package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/ouevre/utils"
	. "github.com/thomasmckinstry/ouevre/utils"
)

type tagKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
	Unfocus key.Binding
	Delete  key.Binding
}

var defaultTagMap = tagKeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "Navigate up a tag"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "Navigate down a tag"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Focus the input, or confirm the text input"),
	),
	Unfocus: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Unfocus the input"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("delete", "Remove a tag"),
	),
}

var (
	tagStyle lipgloss.Style = lipgloss.NewStyle().
			Foreground(utils.Focused)
	tagsStyle lipgloss.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Unfocused).
			BorderTop(true)
)

type TagInputModel struct {
	textInput  textinput.Model
	Tags       []string
	tagsCursor int
	title      string
	selected   bool // Top level focus, can navigate tags, no text entry
	tagCnt     int
	tagStart   int
	tagEnd     int
	width      int

	errorMsg string
}

func (m *TagInputModel) Clear() {
	m.Tags = []string{}
	m.tagsCursor = 0
	m.tagStart = 0
	m.tagEnd = 0
	m.textInput.Reset()
}

func (m *TagInputModel) GetContents() []string {
	return m.Tags
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
		Tags:       tags,
		textInput:  input,
		tagsCursor: 0,
		title:      title,
		tagCnt:     tagCnt,
		selected:   selected,
		width:      width,
	}
}

func (m *TagInputModel) Init() tea.Cmd {
	return nil
}

func (m *TagInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.errorMsg += msg.String()
		switch {
		case key.Matches(msg, defaultTagMap.Unfocus): // Unfocus the component
			if m.textInput.Focused() {
				m.textInput.Blur()
			} else if m.selected {
				m.selected = false
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Confirm): // Add a tag from the current text input and empty the text input OR focus the component
			if !m.selected {
				m.selected = true
			} else if m.selected && !m.textInput.Focused() {
				m.textInput.Focus()
			} else if m.textInput.Value() != "" {
				m.Tags = append(m.Tags, m.textInput.Value())
				if m.tagEnd-m.tagStart < m.tagCnt {
					m.tagEnd++
				}
				m.textInput.Reset()
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Down): // Nav between tags
			if m.tagsCursor < len(m.Tags)-1 && !m.textInput.Focused() && m.selected {
				m.tagsCursor++
				if m.tagsCursor == m.tagEnd {
					m.tagEnd++
					m.tagStart++
				}
			} else { // TODO: If I could come up with a way to avoid this duplication that would be great
				m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
				cmds = tea.Batch(cmds, cmd)
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Up): // Nav between tags
			if m.tagsCursor > 0 && !m.textInput.Focused() && m.selected {
				m.tagsCursor--
				if m.tagsCursor == m.tagStart-1 {
					m.tagEnd--
					m.tagStart--
				}
			} else {
				m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
				cmds = tea.Batch(cmds, cmd)
			}
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Delete):
			if len(m.Tags) > 0 && !m.textInput.Focused() {
				m.Tags = append(m.Tags[:m.tagsCursor], m.Tags[m.tagsCursor+1:]...)
				if m.tagsCursor >= len(m.Tags) {
					m.tagsCursor--
					m.tagEnd--
					if m.tagStart > 0 {
						m.tagStart--
					}
				}
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
			cmd = func() tea.Msg { return NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		}
	}
	return m, cmds
}

func (m *TagInputModel) View() tea.View {
	var s string
	var c = m.textInput.Cursor()
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.title)
	if m.textInput.Focused() {
		c.Y += lipgloss.Height(s)
	}

	clipped := lipgloss.NewStyle().MaxWidth(m.width + 1).Render(m.textInput.View())
	s = lipgloss.JoinVertical(lipgloss.Left, s, clipped)

	isFocused := m.selected && !m.textInput.Focused()

	if len(m.Tags) == 0 {
		s = lipgloss.JoinVertical(lipgloss.Left, s, RenderFocused(tagsStyle, lipgloss.PlaceHorizontal(m.width+2, lipgloss.Center, ""), isFocused))
	}

	for index, tag := range m.Tags[m.tagStart:m.tagEnd] {
		tagStr := ""
		if tag == "" {
			continue
		}
		tagStr = lipgloss.PlaceHorizontal(m.width+2, lipgloss.Left, TruncateString(" - "+tag, m.width+2))
		isFocused = index+m.tagStart == m.tagsCursor && !m.textInput.Focused() && m.selected
		if isFocused && m.selected {
			tagStr = tagStyle.Render(tagStr)
		}

		if index == 0 {
			tagStr = tagsStyle.BorderForeground(Focused).Render(tagStr)
		}

		s = lipgloss.JoinVertical(lipgloss.Left, s, tagStr)
	}
	v := tea.NewView(s)
	v.Cursor = c
	return v
}
