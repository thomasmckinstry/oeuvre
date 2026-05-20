package components

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
)

type CheckboxModel struct {
	cursor    int
	entries   []string
	EntryVals []bool
	Title     string
	selected  bool
	width     int

	entryStyle lipgloss.Style
}

type checkboxKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
	Unfocus key.Binding
}

var defaultCheckboxMap = checkboxKeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "Move up a checkbox"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "Move down a checkbox"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select a checkbox"),
	),
	Unfocus: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Unfocus the checkbox navigation"),
	),
}

func (m *CheckboxModel) Clear() {
	m.EntryVals = make([]bool, len(m.entries))
}

func (m *CheckboxModel) GetContents() []string {
	var contents []string
	for i, entry := range m.EntryVals {
		if entry {
			contents = append(contents, m.entries[i])
		}
	}
	return contents
}

func InitialCheckbox(entries []string, title string, width int) CheckboxModel {
	return CheckboxModel{
		cursor:    0,
		entries:   entries,
		width:     width,
		Title:     title,
		EntryVals: make([]bool, len(entries)),
		entryStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D17600")),
	}
}

func (*CheckboxModel) Init() tea.Cmd {
	return nil
}

func (m *CheckboxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(os.Getenv("DEBUG")) > 0 {
			log.Println("Checkbox input received", msg.String())
		}
		switch {
		case key.Matches(msg, defaultCheckboxMap.Unfocus): // Unfocus the component
			m.selected = false
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
		case key.Matches(msg, defaultCheckboxMap.Confirm): // Add a tag from the current text input and empty the text input OR focus the component
			if !m.selected {
				m.selected = true
			} else {
				m.EntryVals[m.cursor] = !m.EntryVals[m.cursor]
			}
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
		case key.Matches(msg, defaultCheckboxMap.Down): // Nav between tags
			if m.cursor < len(m.entries)-1 && m.selected {
				m.cursor++
			}
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }
		case key.Matches(msg, defaultCheckboxMap.Up):
			if m.cursor > 0 && m.selected {
				m.cursor--
			}
			cmd = func() tea.Msg { return utils.NavMsg(!m.selected) }

		}
	}
	return m, cmd
}

func (m *CheckboxModel) View() tea.View {
	var s string
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.Title)
	for i, medium := range m.entries {
		var entry string
		check := " "
		if m.EntryVals[i] {
			check = "x"
		}
		entry = lipgloss.PlaceHorizontal(m.width-2, lipgloss.Center, medium)
		entry = lipgloss.JoinHorizontal(lipgloss.Center, fmt.Sprintf(" [%s] ", check), entry)
		if i == m.cursor && m.selected {
			entry = m.entryStyle.Render(entry)
		}
		s = lipgloss.JoinVertical(lipgloss.Center, s, entry)
	}
	return tea.NewView(s)
}
