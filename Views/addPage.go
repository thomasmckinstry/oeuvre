package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"encoding/json"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"time"
)

const (
	titleForm int = iota
	yearForm
	mediumForm
	tagsForm
	statusForm
)

type AddModel struct {
	headerText     string
	focused        bool
	cursor         int
	height         int
	width          int
	forms          []tea.Model
	style          lipgloss.Style
	headerStyle    lipgloss.Style
	textinputStyle lipgloss.Style
	enterStyle     lipgloss.Style
}

func clearComponents(m *AddModel) {
	if m.cursor == len(m.forms) {
		m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#6E3F00"))
	}
	m.cursor = 0
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

func InitialAddModel(width int) *AddModel {
	title := components.InitialTextInput(width, "Title", "{ title }", nil)
	year := components.InitialTextInput(width, "Year", "{ year }", nil)
	mediums := []string{"Movie", "Book", "Show", "Anime", "Manga", "Comic", "Show", "Animated", "Live Action"} // TODO: Query the db for this.
	medium := components.InitialCheckbox(mediums, "Medium", width)
	statuses := []string{"Pending", "Started", "Hiatus", "Completed", "Dropped"} // TODO: Query the db for this.
	status := components.InitialCheckbox(statuses, "Status", width)

	var tagSuggestions []string
	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM tags_table`)
	utils.CheckError("Failed to query tags from database: ", err)
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		utils.CheckError("Failed to scan tags: ", err)
		tagSuggestions = append(tagSuggestions, tag)
	}
	err = rows.Close()
	utils.CheckError("Failed to close tags query: ", err)

	tags := components.InitialInput(20, "{ tags }", "Tags", width, false, tagSuggestions)
	forms := []tea.Model{&title, &year, &tags, &medium, &status}
	return &AddModel{
		headerText: "Add Work:",
		forms:      forms,
		focused:    false,
		cursor:     0,
		height:     height,
		width:      24,
		style: lipgloss.NewStyle().
			Height(height).
			Align(lipgloss.Center).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.DoubleBorder()),
		headerStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(20),
		textinputStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderLeft(true),
		enterStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")),
	}
}

func (m *AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (*AddModel, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Height(msg.Height - (7))
		m.height = msg.Height - 7
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.cursor == len(m.forms) {
				var contents []string
				var content string
				var err error
				for _, form := range m.forms {
					switch form := form.(type) {
					case *components.TextInputModel:
						content = string(form.GetContents())
					case *components.TagInputModel:
						marshaledContent, err := json.Marshal(form.GetContents())
						utils.CheckError("Failed to marshal input data to JSON: ", err)
						content = string(marshaledContent)
					case *components.CheckboxModel:
						inputName := form.Title
						if inputName == "Medium" {
							entries := form.GetContents()
							var convertedContents []int
							for _, entry := range entries {
								convertedContents = append(convertedContents, utils.Medium_stoi(entry))
							}
							marshaledContent, err := json.Marshal(convertedContents)
							utils.CheckError("Failed to marshal input data to JSON: ", err)
							content = string(marshaledContent)
						} else {
							content = form.GetContents()[0]
						}
					}
					utils.CheckError("Failed to marshal input data to JSON: ", err)
					contents = append(contents, string(content))
				}
				utils.DebugLog("Inserting to db: ", contents)

				date := time.Now().Format(time.UnixDate)
				db := database.GetDB()
				query, err := db.Prepare(`
					INSERT INTO works (date_added, title, media_type, work_status, tags, year_released)
					VALUES (?, ?, ?, ?, ?, ?)
				`)
				utils.CheckError("Failed to prepare insert statement: ", err)
				statusInt := int(contents[statusForm][0])
				utils.CheckError("Failed to convert string to int: ", err)
				_, err = query.Exec(date, contents[titleForm], contents[tagsForm], statusInt, contents[mediumForm], contents[yearForm])
				utils.CheckError("Failed to insert to works table: ", err)
				err = query.Close()
				utils.CheckError("Failed to close insert to works table: ", err)
				cmd = func() tea.Msg { return ViewMsg(0) }
				cmds = tea.Batch(cmds, cmd)
				clearComponents(m)
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			m.focused = true
		case "esc":
			if !m.focused {
				clearComponents(m)
				cmds = tea.Batch(cmds, func() tea.Msg { return (ViewMsg(0)) })
			} else if m.cursor < len(m.forms) {
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd) // TODO: This needs to return a navMsg to determine if I can unfocus
				msg, ok := cmd().(components.NavMsg)
				if ok && bool(msg) {
					m.focused = false
				}
			}
		case "j", "down":
			if m.cursor > len(m.forms)-1 {
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			msg, ok := cmd().(components.NavMsg)
			if m.cursor < len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.cursor >= len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#D17600"))
			}
		case "k", "up":
			if m.cursor == len(m.forms) {
				m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#6E3F00"))
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			msg, ok := cmd().(components.NavMsg)
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

func (m *AddModel) View() tea.View {
	var c *tea.Cursor
	s := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.headerText)

	for i, form := range m.forms {
		formView := form.View()
		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s) + 1 // TODO: Make the + 2 not hardcoded
			c.X += 2
		}
		if i == m.cursor {
			s = lipgloss.JoinVertical(lipgloss.Center, s, m.textinputStyle.BorderForeground(lipgloss.Color("#D17600")).Render(formView.Content))
		} else {
			s = lipgloss.JoinVertical(lipgloss.Center, s, m.textinputStyle.Render(formView.Content))
		}
		s += "\n"
	}
	enter := m.enterStyle.Render(lipgloss.PlaceHorizontal(15, lipgloss.Center, "ENTER"))
	enter = lipgloss.PlaceVertical(m.height-lipgloss.Height(s)-1, lipgloss.Bottom, enter)
	s = lipgloss.JoinVertical(lipgloss.Center, s, enter)

	v := tea.NewView(m.style.Render(s))
	v.Cursor = c
	v.AltScreen = true
	return v
}
