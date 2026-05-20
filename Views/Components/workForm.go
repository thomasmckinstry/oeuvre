package components

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"encoding/json"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"strings"
)

type WorkFormMap struct {
	Up      key.Binding
	Down    key.Binding
	Focus   key.Binding
	Unfocus key.Binding
}

var DefaultWorkFormKeyMap = WorkFormMap{
	Up: key.NewBinding(
		key.WithKeys("K", "k", "up"),
		key.WithHelp("K/k/↑", "Move up between sections"),
	),
	Down: key.NewBinding(
		key.WithKeys("J", "j", "down"),
		key.WithHelp("J/j/↓", "Move down between sections"),
	),
	Focus: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Confirm an input or focus a component"),
	),
	Unfocus: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Unfocus a component"),
	),
}

type WorkValuesMsg []string

type WorkFormModel struct {
	headerText     string
	focused        bool
	cursor         int
	height         int
	width          int
	forms          []tea.Model
	headerStyle    lipgloss.Style
	textinputStyle lipgloss.Style
	enterStyle     lipgloss.Style
}

func (m *WorkFormModel) ClearComponents() {
	if m.cursor == len(m.forms) {
		m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#6E3F00"))
	}
	m.cursor = 0
	for _, form := range m.forms {
		switch form := form.(type) {
		case *TextInputModel:
			form.Clear()
		case *CheckboxModel:
			form.Clear()
		case *TagInputModel:
			form.Clear()
		}
	}
}

func InitialWorkFormModel(width int, height int) *WorkFormModel {
	title := InitialTextInput(width, "Title", "{ title }", nil)
	year := InitialTextInput(width, "Year", "{ year }", nil)
	mediums := []string{"Anime", "Manga", "Movie", "Book", "Comic", "Show", "Animated", "Live Action"} // TODO: Query the db for this.
	medium := InitialCheckbox(mediums, "Medium", width)
	statuses := []string{"Pending", "Started", "Hiatus", "Completed", "Dropped"} // TODO: Query the db for this.
	status := InitialArrow(statuses, "Status", width, 3)

	var tagSuggestions []string
	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM tags_table`)
	CheckError("Failed to query tags from database: ", err)
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		CheckError("Failed to scan tags: ", err)
		tagSuggestions = append(tagSuggestions, tag)
	}
	err = rows.Close()
	CheckError("Failed to close tags query: ", err)

	tags := InitialInput(20, "{ tags }", "Tags", width-1, false, tagSuggestions)
	forms := []tea.Model{&title, &year, &tags, &medium, &status}
	return &WorkFormModel{
		headerText: "Add Work:",
		forms:      forms,
		focused:    false,
		cursor:     0,
		height:     height,
		width:      24,
		headerStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(width + 2),
		textinputStyle: lipgloss.NewStyle().
			Width(width + 3).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderLeft(true),
		enterStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")),
	}
}

func (m *WorkFormModel) Init() tea.Cmd {
	return nil
}

func (m *WorkFormModel) Update(msg tea.Msg) (*WorkFormModel, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case WorkDetails:
		details := []string(msg)
		for i, form := range m.forms {
			switch form := form.(type) {
			case *TextInputModel:
				if i == 0 {
					form.Textinput.SetValue(details[Title])
				} else {
					form.Textinput.SetValue(details[Released])
				}
			case *TagInputModel:
				form.Tags = strings.Split(details[Tags], ", ")
			case *CheckboxModel:
				for _, entry := range strings.Split(details[Medium], ", ") {
					index := Medium_stoi(entry)
					form.EntryVals[index] = true
				}
			case *ArrowModel:
				form.OptionsCursor = Status_stoi(details[Status])
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultWorkFormKeyMap.Focus):
			if m.cursor == EnterForm {
				var (
					contents []string
					content  string
					tags     []string
					err      error
				)

				// PREPPING DATA FOR ENTRY TO DB
				for _, form := range m.forms {
					switch form := form.(type) {
					case *TextInputModel:
						content = string(form.GetContents())
					case *TagInputModel:
						tags = form.GetContents()
						marshaledContent, err := json.Marshal(tags)
						CheckError("Failed to marshal input data to JSON: ", err)
						content = string(marshaledContent)
					case *CheckboxModel:
						entries := form.GetContents()
						var convertedContents []int
						for _, entry := range entries {
							convertedContents = append(convertedContents, Medium_stoi(entry))
						}
						marshaledContent, err := json.Marshal(convertedContents)
						CheckError("Failed to marshal input data to JSON: ", err)
						content = string(marshaledContent)
					case *ArrowModel:
						content = string(Status_stoi(form.GetContents()))
					}
					CheckError("Failed to marshal input data to JSON: ", err)
					contents = append(contents, string(content))
				}
				cmds = tea.Batch(cmds, func() tea.Msg { return NewWorkMsg(contents) })
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			m.focused = true
		case key.Matches(msg, DefaultWorkFormKeyMap.Unfocus):
			if !m.focused {
				cmds = tea.Batch(cmds, func() tea.Msg { return (ViewMsg(0)) })
			} else if m.cursor < len(m.forms) {
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				msg, ok := cmd().(NavMsg)
				if ok && bool(msg) {
					m.focused = false
				}
			}
		case key.Matches(msg, DefaultWorkFormKeyMap.Down):
			if m.cursor >= len(m.forms) {
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			msg, ok := cmd().(NavMsg)
			if m.cursor < len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.cursor >= len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#D17600"))
			}
		case key.Matches(msg, DefaultWorkFormKeyMap.Up):
			if m.cursor == len(m.forms) {
				m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#6E3F00"))
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			msg, ok := cmd().(NavMsg)
			if m.cursor > 0 && ok && bool(msg) {
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		default:
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
		}
	}
	return m, cmds
}

func (m *WorkFormModel) View() tea.View {
	var c *tea.Cursor
	s := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.headerText)

	for i, form := range m.forms {
		formView := form.View()
		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s)
		}
		if i == m.cursor {
			s = lipgloss.JoinVertical(lipgloss.Left, s, m.textinputStyle.BorderForeground(lipgloss.Color("#D17600")).Render(formView.Content))
		} else {
			s = lipgloss.JoinVertical(lipgloss.Left, s, m.textinputStyle.Render(formView.Content))
		}
		s += "\n"
	}
	enter := m.enterStyle.Render(lipgloss.PlaceHorizontal(15, lipgloss.Center, "CONFIRM"))
	enter = lipgloss.PlaceVertical(m.height-lipgloss.Height(s)-1, lipgloss.Bottom, enter)
	s = lipgloss.JoinVertical(lipgloss.Center, s, enter)

	v := tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
