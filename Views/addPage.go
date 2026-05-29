package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"time"
)

var addStyle lipgloss.Style = lipgloss.NewStyle().
	Align(lipgloss.Center).
	PaddingLeft(1).
	PaddingRight(1).
	BorderStyle(lipgloss.DoubleBorder())

type AddModel struct {
	form          *components.WorkFormModel
	width, height int
}

func InitialAddModel(width, height int) *AddModel {
	form := components.InitialWorkFormModel(25, height)
	return &AddModel{
		width:  width,
		height: height,
		form:   form,
	}
}

func (m *AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (*AddModel, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ViewMsg:
		m.form.ClearComponents()
	case tea.WindowSizeMsg:
		m.form.Update(msg)
		m.height = msg.Height
		m.width = msg.Width
	default:
		_, cmd = m.form.Update(msg)
		var ok bool
		var workMsg []string
		if cmd != nil {
			workMsg, ok = cmd().(NewWorkMsg)
		}
		if ok {
			DebugLog("addPage got NewWorkMsg: ", msg)
			cmds = tea.Batch(cmds, func() tea.Msg { return ViewMsg(0) })
			// ADDING TO DATABASE
			db := database.GetDB()
			date := time.Now().Format(time.UnixDate)

			var id int64
			query, err := db.Prepare(`
					INSERT OR REPLACE INTO works (date_added, title, media_type, work_status, tags, year_released)
					VALUES (?, ?, ?, ?, ?, ?)
					RETURNING work_id
				`)
			CheckError("Failed to prepare insert statement: ", err)
			statusInt := int(workMsg[StatusForm][0])
			CheckError("Failed to convert string to int: ", err)
			row, err := query.Exec(date, workMsg[TitleForm], workMsg[MediumForm], statusInt, workMsg[TagsForm], workMsg[YearForm])
			id, err = row.LastInsertId()
			CheckError("Failed to insert to works table: ", err)
			err = query.Close()

			query, err = db.Prepare(`
				 INSERT OR REPLACE INTO tags_table (tag_name)
				 VALUES (?)
				`)
			CheckError("Failed to prepare insert statement: ", err)
			for _, tag := range workMsg[TagsForm] {
				_, err = query.Exec(tag)
			}
			err = query.Close()

			CheckError("Failed to close insert to works table: ", err)
			cmds = tea.Batch(cmds, func() tea.Msg { return ViewMsg(0) })
			cmds = tea.Batch(cmds, func() tea.Msg { return append(workMsg, string(id)) })
		}
		cmds = tea.Batch(cmds, cmd)
	}
	return m, cmds
}

func (m *AddModel) View() tea.View {
	var c *tea.Cursor
	var s string

	formView := m.form.View()
	c = formView.Cursor

	s = addStyle.Render(formView.Content)

	if c != nil {
		c.Y += 1
		c.X += (m.width / 2) - (lipgloss.Width(s) / 2) + 2
	}

	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, s)
	s = lipgloss.PlaceVertical(m.height, lipgloss.Center, s)

	v := tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
