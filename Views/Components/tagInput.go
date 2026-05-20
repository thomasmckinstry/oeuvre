package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
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

type TagInputModel struct {
	textInput  textinput.Model
	Tags       []string
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

func (m *TagInputModel) Clear() {
	m.Tags = []string{}
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
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Confirm): // Add a tag from the current text input and empty the text input OR focus the component
			if !m.selected {
				m.selected = true
			} else if m.selected && !m.textInput.Focused() {
				m.textInput.Focus()
			} else if m.textInput.Value() != "" && len(m.Tags) < m.tagCnt {
				m.Tags = append(m.Tags, m.textInput.Value())
				m.textInput.Reset()
			}
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Down): // Nav between tags
			if m.tagsCursor < len(m.Tags)-1 && !m.textInput.Focused() && m.selected {
				m.tagsCursor++
			} else { // TODO: If I could come up with a way to avoid this duplication that would be great
				m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
				cmds = tea.Batch(cmds, cmd)
			}
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Up): // Nav between tags
			if m.tagsCursor > 0 && !m.textInput.Focused() && m.selected {
				m.tagsCursor--
			} else {
				m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
				cmds = tea.Batch(cmds, cmd)
			}
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultTagMap.Delete):
			if len(m.Tags) > 0 && !m.textInput.Focused() {
				m.Tags = append(m.Tags[:m.tagsCursor], m.Tags[m.tagsCursor+1:]...)
				if m.tagsCursor >= len(m.Tags) {
					m.tagsCursor--
				}
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg) // Default to typing in the text input
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
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
		c.X += 1 // Aligns it correctly with the text
	}

	s = lipgloss.JoinVertical(lipgloss.Left, s, m.textInput.View())

	var wrapTags lipgloss.Style
	if m.selected && !m.textInput.Focused() {
		wrapTags = m.tagsStyle.BorderForeground(lipgloss.Color("#D17600"))
	} else {
		wrapTags = m.tagsStyle
	}

	if len(m.Tags) == 0 {
		s = lipgloss.JoinVertical(lipgloss.Left, s, wrapTags.Render())
	}

	for index, tag := range m.Tags {
		tagStr := ""
		if tag == "" {
			continue
		}
		tagStr = " - " + tag
		if index == m.tagsCursor && !m.textInput.Focused() && m.selected { // Color selected field
			tagStr = m.tagStyle.Render(tagStr)
		}

		if index == 0 {
			tagStr = wrapTags.Render(tagStr)
		}

		s = lipgloss.JoinVertical(lipgloss.Left, s, tagStr)
	}
	v := tea.NewView(s)
	v.Cursor = c
	return v
}
