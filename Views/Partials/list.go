package partials

import (
	"database/sql"
	"encoding/json"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	database "github.com/thomasmckinstry/ouevre/db"
	. "github.com/thomasmckinstry/ouevre/utils"
	"slices"
	"sort"
	"strings"
)

var db *sql.DB

type listKeyMap struct {
	Nav     key.Binding
	Focus   key.Binding
	Confirm key.Binding
}

type SortMsg int

var defaultListMap = listKeyMap{
	Nav:     key.NewBinding(key.WithKeys("H", "L")),
	Focus:   key.NewBinding(key.WithKeys("esc")),
	Confirm: key.NewBinding(key.WithKeys("enter")),
}

type ListModel struct {
	table table.Model
	rows  []table.Row
}

func (m *ListModel) refreshList() {
	db = database.GetDB()

	row, err := db.Query(`SELECT title, media_type, work_status, tags, year_released, work_id FROM works;`)
	defer func() {
		err = row.Close()
		CheckError("Failed to close works query: ", err)
	}()

	var rows []table.Row
	for row.Next() {
		var (
			intStatus int
			title     string
			medium    string
			tags      string
			year      string
			id        string
		)
		err = row.Scan(&title, &medium, &intStatus, &tags, &year, &id)
		CheckError("Failed to scan works row: ", err)
		var mediumsArr []int
		var tagsArr []string
		err := json.Unmarshal([]byte(medium), &mediumsArr)
		CheckError("Failed to Unmarshal medium: ", err)
		err = json.Unmarshal([]byte(tags), &tagsArr)
		CheckError("Failed to Unmarshal tags: ", err)
		mediumsStr := ConvertMedium(mediumsArr)
		rows = append(rows, table.Row{title, GetTagsString(tagsArr), mediumsStr, Status_itos(intStatus), year, id})
	}

	m.table.SetRows(rows)
	m.rows = rows
}

func InitialList(width int, height int) ListModel {
	var columns = []table.Column{
		{Title: "Title", Width: width / 4},
		{Title: "Tags", Width: width / 3},
		{Title: "Medium", Width: width / 8},
		{Title: "Status", Width: width / 8},
		{Title: "Released", Width: width / 6},
		{Title: "Id", Width: 0},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(height-2),
		table.WithWidth(width),
	)
	t.Blur()

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Unfocused).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(Focused).
		Bold(false)
	t.SetStyles(s)

	model := ListModel{
		table: t,
	}

	model.refreshList()

	return model
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case DeleteWorkMsg, NewWorkMsg:
		m.refreshList()
	case tea.WindowSizeMsg:
		width := msg.Width - 29
		m.table.SetWidth(width)
		m.table.SetHeight(msg.Height - 2)
		m.table.SetColumns([]table.Column{
			{Title: "Title", Width: width / 4},
			{Title: "Tags", Width: width / 3},
			{Title: "Medium", Width: width / 8},
			{Title: "Status", Width: width / 8},
			{Title: "Released", Width: width / 6},
			{Title: "Id", Width: 0},
		})
		m.table.Update(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultListMap.Focus):
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case key.Matches(msg, defaultListMap.Nav):
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case key.Matches(msg, defaultListMap.Confirm):
			currWork := WorkDetails(m.table.SelectedRow())
			if len(currWork) > 0 {
				cmd = func() tea.Msg { return ViewMsg(2) }
				cmd = tea.Batch(cmd, func() tea.Msg { return WorkDetails(m.table.SelectedRow()) })
			}
		default:
			m.table, cmd = m.table.Update(msg)
		}
	case SortMsg:
		rows := m.table.Rows()
		sort.Slice(rows, func(i, j int) bool {
			return rows[i][int(msg)] < rows[j][int(msg)]
		})
		m.table.SetRows(rows)
	case FilterMsg:
		var rows []table.Row
		filter := [][]string(msg)
		for _, row := range m.rows {
			include := true
			if row[Title] != filter[Title][0] && filter[Title][0] != "" {
				continue
			} else if len(filter[Status]) > 0 && row[Status] != filter[Status][0] {
				continue
			}
			tags := strings.Split(row[Tags], ", ")
			for _, tag := range filter[Tags] {
				if !slices.Contains(tags, tag) {
					include = false
					break
				}
			}
			if !include {
				continue
			}
			mediums := strings.Split(row[Medium], ", ")
			for _, medium := range filter[Medium] {
				if !slices.Contains(mediums, medium) {
					include = false
					break
				}
			}
			if !include {
				continue
			}
			if include {
				rows = append(rows, row)
			}
		}
		m.table.SetRows(rows)
	}
	return m, cmd
}

func (m ListModel) View() tea.View {
	return tea.NewView(m.table.View())
}
