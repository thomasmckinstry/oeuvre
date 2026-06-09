package partials

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"encoding/json"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"strconv"
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

var errorStyle lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Red)

type WorkValuesMsg []string

type WorkFormModel struct {
	headerText string
	errorMsg   string
	focused    bool
	cursor     int
	height     int
	width      int
	forms      []tea.Model
}

func (m *WorkFormModel) ClearComponents() {
	m.cursor = 0
	m.errorMsg = ""
	for _, form := range m.forms {
		switch form := form.(type) {
		case *components.TextInputModel:
			form.Clear()
		case *components.CheckboxModel:
			form.Clear()
		case *components.TagInputModel:
			form.Clear()
		}
	}
}

func InitialWorkFormModel(width int, height int) *WorkFormModel {
	title := components.InitialTextInput(width, "Title", "{ title }", nil)
	year := components.InitialTextInput(width, "Year", "{ year }", nil)
	mediums := Config.MediaOptions
	medium := components.InitialCheckbox(mediums, "Medium", width)
	statuses := Config.StatusOptions
	status := components.InitialArrow(statuses, "Status", width, 3)

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

	tags := components.InitialInput(20, "{ tags }", "Tags", width-1, false, tagSuggestions)
	forms := []tea.Model{&title, &year, &tags, &medium, &status}
	return &WorkFormModel{
		headerText: "Add Work:",
		forms:      forms,
		focused:    false,
		cursor:     0,
		height:     height,
		width:      24,
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
			case *components.TextInputModel:
				if i == 0 {
					form.Textinput.SetValue(details[Title])
				} else {
					form.Textinput.SetValue(details[Released])
				}
			case *components.TagInputModel:
				form.Tags = strings.Split(details[Tags], ", ")
			case *components.CheckboxModel:
				for _, entry := range strings.Split(details[Medium], ", ") {
					index := Medium_stoi(entry)
					form.EntryVals[index] = true
				}
			case *components.ArrowModel:
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
					valid    bool
				)

				// PREPPING DATA FOR ENTRY TO DB
				for i, form := range m.forms {
					switch form := form.(type) {
					case *components.TextInputModel:
						content = string(form.GetContents())
						if i == YearForm {
							_, err = strconv.Atoi(content)
							if err != nil {
								valid = false
								m.errorMsg = "Invalid input for Year."
								err = nil
							}
						}
					case *components.TagInputModel:
						tags = form.GetContents()
						marshaledContent, err := json.Marshal(tags)
						CheckError("Failed to marshal input data to JSON: ", err)
						content = string(marshaledContent)
					case *components.CheckboxModel:
						entries := form.GetContents()
						var convertedContents []int
						if len(entries) == 0 {
							valid = false
							m.errorMsg = "Invalid input for Medium."
							err = nil
							break
						}
						for _, entry := range entries {
							convertedContents = append(convertedContents, Medium_stoi(entry))
						}
						marshaledContent, err := json.Marshal(convertedContents)
						CheckError("Failed to marshal input data to JSON: ", err)
						content = string(marshaledContent)
					case *components.ArrowModel:
						content = string(Status_stoi(form.GetContents()))
					}
					CheckError("Failed to marshal input data to JSON: ", err)
					contents = append(contents, string(content))
				}
				if valid {
					cmds = tea.Batch(cmds, func() tea.Msg { return NewWorkMsg(contents) })
				}
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			m.focused = true
		case key.Matches(msg, DefaultWorkFormKeyMap.Unfocus):
			if !m.focused || m.cursor == EnterForm {
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
			}
		case key.Matches(msg, DefaultWorkFormKeyMap.Up):
			if m.cursor == len(m.forms) {
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
			if m.cursor != EnterForm {
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
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
			s = lipgloss.JoinVertical(lipgloss.Left, s, textinputStyle.BorderForeground(lipgloss.Color("#D17600")).Render(formView.Content))
		} else {
			s = lipgloss.JoinVertical(lipgloss.Left, s, textinputStyle.Render(formView.Content))
		}
		s += "\n"
	}
	isFocused := m.cursor == EnterForm
	enter := RenderFocused(enterStyle, "CONFIRM", isFocused)
	if m.errorMsg != "" {
		enter = lipgloss.JoinVertical(lipgloss.Center, errorStyle.Render(m.errorMsg), enter)
	}
	enter = lipgloss.PlaceVertical(m.height-lipgloss.Height(s)-2, lipgloss.Bottom, enter)
	s = lipgloss.JoinVertical(lipgloss.Center, s, enter)

	v := tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
