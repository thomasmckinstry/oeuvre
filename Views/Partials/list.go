package partials

import (
	"database/sql"
	"encoding/json"
	"log"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var db *sql.DB

type listKeyMap struct {
	Nav   key.Binding
	Focus key.Binding
}

type SortMsg int

var defaultListMap = listKeyMap{
	Nav:   key.NewBinding(key.WithKeys("H", "L")),
	Focus: key.NewBinding(key.WithKeys("esc")),
}

type ListModel struct {
	style lipgloss.Style
	table table.Model
	rows  []table.Row
}

func (m ListModel) toggleSelected() lipgloss.Style {
	if !m.table.Focused() {
		return m.style.BorderForeground(lipgloss.Color("#D17600"))
	} else {
		return m.style.BorderForeground(lipgloss.Color("#6E3F00"))
	}
}

func InitialList(width int, height int) ListModel {
	db = database.GetDB()
	row, err := db.Query(`SELECT title, media_type, work_status, tags, year_released FROM works;`)
	defer func() {
		err = row.Close()
		utils.CheckError("Failed to close works query: ", err)
	}()

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
		utils.DebugLog("Scanned row: ", []string{title, medium, tags, year, string(intStatus)})
		var mediumsArr []int
		var tagsArr []string
		err := json.Unmarshal([]byte(medium), &mediumsArr)
		if err != nil {
			log.Fatal("Failed to Unmarshal medium: ", err)
		}
		err = json.Unmarshal([]byte(tags), &tagsArr)
		if err != nil {
			log.Fatal("Failed to Unmarshal tags: ", err)
		}
		mediumsStr := utils.ConvertMedium(mediumsArr)
		rows = append(rows, table.Row{title, utils.GetTagsString(tagsArr), mediumsStr, utils.Status_itos(intStatus), year})
	}

	var columns = []table.Column{
		{Title: "Title", Width: width / 4},
		{Title: "Tags", Width: width / 3},
		{Title: "Medium", Width: width / 8},
		{Title: "Status", Width: width / 8},
		{Title: "Released", Width: width / 6},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(width),
	)
	t.Blur()

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
		rows:  rows,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case utils.NewWorkMsg:
		utils.DebugLog("List got NewWorkMsg: ", msg)
		newRow := []string(msg)
		rows := m.table.Rows()
		var mediumsArr []int
		var tagsArr []string
		err := json.Unmarshal([]byte(newRow[utils.MediumForm]), &mediumsArr)
		if err != nil {
			utils.DebugLog("Medium: ", newRow[utils.MediumForm])
			log.Fatal("Failed to Unmarshal medium: ", err)
		}
		err = json.Unmarshal([]byte(newRow[utils.TagsForm]), &tagsArr)
		if err != nil {
			log.Fatal("Failed to Unmarshal tags: ", err)
		}
		intStatus, _ := strconv.Atoi(newRow[utils.StatusForm])
		mediumsStr := utils.ConvertMedium(mediumsArr)
		rows = append(rows, table.Row{newRow[utils.TitleForm], utils.GetTagsString(tagsArr), mediumsStr, utils.Status_itos(intStatus), newRow[utils.YearForm]})
		m.table.SetRows(rows)
	case tea.WindowSizeMsg:
		m.style = m.style.Height(msg.Height).Width(msg.Width - 18)
		width := msg.Width - 29
		m.table.SetColumns([]table.Column{
			{Title: "Title", Width: width / 4},
			{Title: "Tags", Width: width / 3},
			{Title: "Medium", Width: width / 8},
			{Title: "Status", Width: width / 8},
			{Title: "Released", Width: width / 6},
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
			m.style = m.toggleSelected()
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
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
			if row[utils.Title] != filter[utils.Title][0] && filter[utils.Title][0] != "" {
				continue
			} else if len(filter[utils.Status]) > 0 && row[utils.Status] != filter[utils.Status][0] {
				continue
			}
			tags := strings.Split(row[utils.Tags], ", ")
			for _, tag := range filter[utils.Tags] {
				if !slices.Contains(tags, tag) {
					include = false
					break
				}
			}
			if !include {
				continue
			}
			mediums := strings.Split(row[utils.Medium], ", ")
			for _, medium := range filter[utils.Medium] {
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
	return tea.NewView(m.style.Render(m.table.View()))
}
