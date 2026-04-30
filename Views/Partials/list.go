package partials

import (
	"database/sql"
	"encoding/json"
	"log"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	database "github.com/thomasmckinstry/Bubbletea-Tutorial/db"
	"github.com/thomasmckinstry/Bubbletea-Tutorial/utils"
)

var db *sql.DB

type ListModel struct {
	style lipgloss.Style
	table table.Model
}

func (m ListModel) selectView() lipgloss.Style {
	return m.style.BorderForeground(lipgloss.Color("#D17600"))
}

func (m ListModel) deselectView() lipgloss.Style {
	return m.style.BorderForeground(lipgloss.Color("#6E3F00"))
}

func InitialList(width int, height int) ListModel {
	db = database.GetDB()
	row, err := db.Query(`SELECT title, media_type, work_status, tags, year_released FROM works;`)
	defer row.Close()
	if err != nil {
		log.Fatal("Failed to query works table for list: ", err)
	}

	var rows []table.Row
	for row.Next() {
		var (
			intStatus int
			title     string
			medium    string
			tags      string
			year      string
		)
		err = row.Scan(&title, &medium, &intStatus, &tags, &year)
		if err != nil {
			log.Fatal("Failed to scan works row: ", err)
		}
		var mediumsArr []int
		var tagsArr []string
		err := json.Unmarshal([]byte(medium), &mediumsArr)
		if err != nil {
			log.Fatal("Failed to Unmarshal medium: ", err)
		}
		err = json.Unmarshal([]byte(tags), &tagsArr)
		if err != nil {
			log.Fatal("Failed to Unmarshal medium: ", err)
		}
		mediumsStr := utils.ConvertMedium(mediumsArr)
		rows = append(rows, table.Row{title, mediumsStr, utils.Status_itos(intStatus), utils.GetTagsString(tagsArr), year})
	}

	var columns = []table.Column{
		{Title: "Title", Width: width / 4},
		{Title: "Medium", Width: width / 8},
		{Title: "Status", Width: width / 8},
		{Title: "Tags", Width: width / 3},
		{Title: "Released", Width: width / 6},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(width),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#6E3F00")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("#D17600")).
		Bold(false)
	t.SetStyles(s)

	return ListModel{
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("#6E3F00")).
			PaddingTop(1).
			Width(width).
			Height(height),
		table: t,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Height(msg.Height).Width(msg.Width - 18)
		width := msg.Width - 29
		m.table.SetColumns([]table.Column{
			{Title: "Title", Width: width / 4},
			{Title: "Medium", Width: width / 8},
			{Title: "Status", Width: width / 8},
			{Title: "Tags", Width: width / 3},
			{Title: "Released", Width: width / 6},
		})
		m.table.Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "L":
			m.style = m.selectView()
			m.table.Focus()
		case "H":
			m.style = m.deselectView()
			m.table.Blur()
		case "j", "k", "up", "down":
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m ListModel) View() tea.View {
	return tea.NewView(m.style.Render(m.table.View()))
}
