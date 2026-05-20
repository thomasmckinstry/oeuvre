package views

import (
	key "charm.land/bubbles/v2/key"
	textarea "charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"database/sql"
	components "github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"log"
	"strconv"
	"time"
)

type entry struct {
	content, date string
	id            int
}

type WorkPageModel struct {
	work              *components.WorkFormModel
	currWorkId        int
	textArea          textarea.Model
	notes, reviews    []entry
	writingMode       string
	focused           bool
	writing           bool
	width, height     int
	tabCursor         int
	mainCursor        int
	entryCursor       int
	rightCursor       int
	writingCursor     int
	displayCursor     int
	tabStyle          lipgloss.Style
	tabsStyle         lipgloss.Style
	buttonStyle       lipgloss.Style
	displayStyle      lipgloss.Style
	entryContentStyle lipgloss.Style
	detailsStyle      lipgloss.Style
}

type workKeyMap struct {
	TopLevelUp    key.Binding
	TopLevelDown  key.Binding
	TopLevelLeft  key.Binding
	TopLevelRight key.Binding
	Confirm       key.Binding
	Left          key.Binding
	Right         key.Binding
	Exit          key.Binding
}

const (
	work int = iota
	tabs
	add
	mainCursorCount
)

const (
	header int = iota
	display
	rightCursorCount
)

var writeHeader = []string{"NOTE", "REVIEW"}

var (
	unfocused = lipgloss.Color("#6E3F00")
	focused   = lipgloss.Color("#D17600")
)

var defaultWorkMap = workKeyMap{
	TopLevelUp:    key.NewBinding(key.WithKeys("K")),
	TopLevelDown:  key.NewBinding(key.WithKeys("J")),
	TopLevelLeft:  key.NewBinding(key.WithKeys("H")),
	TopLevelRight: key.NewBinding(key.WithKeys("L")),
	Left:          key.NewBinding(key.WithKeys("h", "left")),
	Right:         key.NewBinding(key.WithKeys("l", "right")),
	Confirm:       key.NewBinding(key.WithKeys("enter")),
	Exit:          key.NewBinding(key.WithKeys("esc")),
}

func renderFocused(style lipgloss.Style, content string, isFocused bool) string {
	if isFocused {
		return style.BorderForeground(focused).Render(content)
	} else {
		return style.Render(content)
	}
}

func InitialWorkPage(width, height int) *WorkPageModel {
	workForm := components.InitialWorkFormModel(22, height)

	ti := textarea.New()
	ti.SetWidth(width - 1)
	ti.SetHeight(height - 5)
	ti.SetVirtualCursor(false)
	ti.SetStyles(textarea.DefaultStyles(true)) // default to dark styles.

	return &WorkPageModel{
		work:     workForm,
		textArea: ti,
		width:    width,
		height:   height,
		tabStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(unfocused),
		tabsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			MarginLeft(1).
			MarginRight(1).
			BorderForeground(unfocused),
		detailsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			MarginRight(1).
			BorderForeground(unfocused),
		buttonStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(unfocused),
		displayStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			PaddingRight(1).
			BorderTop(true).
			BorderForeground(unfocused),
		entryContentStyle: lipgloss.NewStyle().
			MarginRight(1),
	}
}

func (m *WorkPageModel) SetWork(workDetails []string) {
}

func (m *WorkPageModel) Init() tea.Cmd {
	return nil
}

func (m *WorkPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd, cmds tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.work, cmd = m.work.Update(msg)
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(msg.Width - 1)
		m.textArea.SetHeight(msg.Height - 5)
	case WorkDetails:
		m.notes = []entry{}
		m.reviews = []entry{}
		_, cmd = m.work.Update(msg)
		DebugLog("WorkDetails: ", msg)
		db := database.GetDB()
		id, err := strconv.Atoi(msg[len(msg)-1])
		m.currWorkId = id
		row, err := db.Query(`SELECT date_added, review_text, review_id FROM reviews WHERE work_id = ?;`, id)
		for row.Next() {
			var (
				content, date string
				id            int
			)
			err = row.Scan(&date, &content, &id)
			if err != nil {
				log.Fatal("Failed to scan works row: ", err)
			}
			m.reviews = append(m.reviews, entry{content: content, date: date, id: id})
		}
		row, err = db.Query(`SELECT date_added, note_text, note_id FROM notes WHERE work_id = ?;`, id)
		DebugLog("Querying notes", id)
		for row.Next() {
			var (
				content, date string
				id            int
			)
			err = row.Scan(&date, &content, &id)
			DebugLog("Scanned note: ", date)
			if err != nil {
				log.Fatal("Failed to scan works row: ", err)
			}
			m.notes = append(m.notes, entry{content: content, date: date, id: id})
		}
		err = row.Close()
		CheckError("Failed to close works query: ", err)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultWorkMap.Confirm):
			if m.mainCursor == work && !m.focused && !m.writing {
				m.focused = true
			} else if m.mainCursor == add && m.rightCursor == header && !m.writing {
				m.writing = true
				m.writingMode = "ADD"
				return m, cmds
			} else if m.writing && m.writingCursor == 0 && !m.textArea.Focused() {
				m.textArea.Focus()
				return m, cmd
			} else if m.writing && m.writingCursor == 1 {
				content := m.textArea.Value()
				db := database.GetDB()
				m.writing = false
				var (
					query *sql.Stmt
					err   error
				)
				date := time.Now().Format(time.DateOnly)
				if m.writingMode == "ADD" {
					if m.tabCursor == 0 {
						query, err = db.Prepare(`INSERT INTO notes (date_added, note_text, work_id) VALUES (?, ?, ?)`)
						m.notes = append(m.notes, entry{content: content, date: date})
					} else {
						query, err = db.Prepare(`INSERT INTO reviews (date_added, review_text, work_id) VALUES  (?, ?, ?)`)
						m.reviews = append(m.reviews, entry{content: content, date: date})
					}
				} else {
					if m.tabCursor == 0 {
						query, err = db.Prepare(`REPLACE INTO notes (date_added, note_text, work_id, review_id) VALUES (?, ?, ?, ?)`)
						m.notes[m.entryCursor] = entry{content: content, date: date, id: m.notes[m.entryCursor].id}
					} else {
						query, err = db.Prepare(`REPLACE INTO reviews (date_added, review_text, work_id, review_id) VALUES (?, ?, ?, ?)`)
						m.reviews[m.entryCursor] = entry{content: content, date: date, id: m.reviews[m.entryCursor].id}
					}

				}
				CheckError("Failed to prepare statement to insert to DB: ", err)
				if m.writingMode == "ADD" {
					query.Exec(date, content, m.currWorkId)
				} else {
					query.Exec(date, content, m.currWorkId, m.reviews[m.entryCursor].id)
				}
				query.Close()
				m.textArea.Reset()
				m.writingCursor = 0

			} else if m.rightCursor == 1 && m.mainCursor != work {
				m.writing = true
				m.writingMode = "EDIT"
				if m.tabCursor == 0 {
					m.textArea.SetValue(m.notes[m.entryCursor].content)
				} else {
					m.textArea.SetValue(m.reviews[m.entryCursor].content)
				}
			} else if m.mainCursor == work {
				m.work, cmd = m.work.Update(msg)
				var ok bool
				var workMsg []string
				if cmd != nil {
					workMsg, ok = cmd().(NewWorkMsg)
				}
				if ok {
					DebugLog("addPage got NewWorkMsg: ", msg)
					// ADDING TO DATABASE
					db := database.GetDB()
					date := time.Now().Format(time.UnixDate)

					query, err := db.Prepare(`
					INSERT OR REPLACE INTO works (date_added, title, media_type, work_status, tags, year_released, work_id)
					VALUES (?, ?, ?, ?, ?, ?, ?)
				`)
					CheckError("Failed to prepare insert statement: ", err)
					statusInt := int(workMsg[StatusForm][0])
					CheckError("Failed to convert string to int: ", err)
					_, err = query.Exec(date, workMsg[TitleForm], workMsg[MediumForm], statusInt, workMsg[TagsForm], workMsg[YearForm], m.currWorkId)
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
				}
				cmd = nil
			}
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultWorkMap.Exit):
			if m.mainCursor == work {
				_, cmd = m.work.Update(msg)
				_, ok := cmd().(ViewMsg)
				if ok && m.focused {
					cmd = nil
					m.focused = false
				}
			} else if m.writing {
				if m.textArea.Focused() {
					m.textArea.Blur()
				} else {
					m.writing = false
					m.writingCursor = 0
					m.textArea.Reset()
				}
			} else {
				cmd = func() tea.Msg { return ViewMsg(0) }
			}
		case key.Matches(msg, defaultWorkMap.TopLevelRight):
			if m.mainCursor < mainCursorCount-1 && !m.focused {
				m.mainCursor++
			} else if m.mainCursor == work {
				_, cmd = m.work.Update(msg)
			}
		case key.Matches(msg, defaultWorkMap.TopLevelLeft):
			if m.rightCursor == display {
				m.mainCursor = 0
			} else if m.mainCursor > work && !m.writing {
				m.mainCursor--
				_, cmd = m.work.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			}
		case key.Matches(msg, defaultWorkMap.TopLevelDown):
			if m.mainCursor > work && m.rightCursor < display && !m.writing {
				m.rightCursor++
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.writing && !m.textArea.Focused() && m.writingCursor == 0 {
				m.writingCursor++
			}
		case key.Matches(msg, defaultWorkMap.TopLevelUp):
			if m.mainCursor > work && m.rightCursor > header && !m.writing {
				m.rightCursor--
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.writing && m.writingCursor > 0 {
				m.writingCursor--
			}
		case key.Matches(msg, defaultWorkMap.Left):
			if m.mainCursor == tabs && m.rightCursor == header && !m.writing {
				m.tabCursor = 0
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			}
		case key.Matches(msg, defaultWorkMap.Right):
			if m.mainCursor == tabs && m.rightCursor == header && !m.writing {
				m.tabCursor = 1
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			}
		default:
			if m.mainCursor == work && m.focused {
				_, cmd = m.work.Update(msg)
			}
		}
	default:
		m.textArea, cmd = m.textArea.Update(msg)
	}

	if m.textArea.Focused() {
		m.textArea, cmd = m.textArea.Update(msg)
	}

	cmds = tea.Batch(cmds, cmd)
	return m, cmds
}

func (m *WorkPageModel) View() tea.View {
	var (
		c *tea.Cursor
		v tea.View
		s string
	)

	if m.writing {
		writingHeader := m.writingMode + " " + writeHeader[m.tabCursor]
		textarea := renderFocused(m.tabsStyle, m.textArea.View(), m.writingCursor == 0)
		button := renderFocused(m.buttonStyle, "CONFIRM", m.writingCursor == 1)

		s = lipgloss.JoinVertical(lipgloss.Center, writingHeader, textarea, button)

		v = tea.NewView(s)
		c = m.textArea.Cursor()
		if c != nil {
			c.Y += 2
			c.X += 1
		}
		v.Cursor = c
		v.AltScreen = true
		return v
	}

	workView := m.work.View()
	c = workView.Cursor
	if c != nil {
		c.Y += 1
		c.X += 1
	}
	details := workView.Content
	isFocused := m.focused
	details = renderFocused(m.tabsStyle, details, isFocused)

	isSelected := m.mainCursor == work
	details = renderFocused(m.detailsStyle.Height(m.height), details, isSelected)

	s = lipgloss.JoinHorizontal(lipgloss.Top, s, details)

	headerContent := "ADD"
	isFocused = m.mainCursor == add && m.rightCursor == header
	headerContent = renderFocused(m.buttonStyle, headerContent, isFocused)

	renderedTabs := []string{}
	tabsArr := []string{"NOTES", "REVIEWS"}
	for i, tab := range tabsArr {
		tab = lipgloss.PlaceHorizontal(9, lipgloss.Center, tab)
		isFocused = m.tabCursor == i
		renderedTabs = append(renderedTabs, renderFocused(m.tabStyle, tab, isFocused))
	}

	tabsContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs[0], " ", renderedTabs[1])
	isFocused = m.mainCursor == tabs && m.rightCursor == header
	tabsContent = renderFocused(m.tabsStyle, tabsContent, isFocused)

	headerContent = lipgloss.JoinHorizontal(lipgloss.Top, tabsContent, headerContent)

	var displayContent, position string
	var currEntry entry
	var hasEntry bool
	if m.tabCursor == 0 && len(m.notes) > 0 {
		currEntry = m.notes[m.entryCursor]
		hasEntry = true
		position = lipgloss.PlaceVertical(m.height-11, lipgloss.Bottom, strconv.Itoa(m.entryCursor+1)+"/"+strconv.Itoa(len(m.notes)))
	} else if m.tabCursor == 1 && len(m.reviews) > 0 {
		currEntry = m.reviews[m.entryCursor]
		hasEntry = true
		position = lipgloss.PlaceVertical(m.height-11, lipgloss.Bottom, strconv.Itoa(m.entryCursor+1)+"/"+strconv.Itoa(len(m.reviews)))
	}

	if hasEntry {
		var date string
		if currEntry.date != "" {
			date = currEntry.date[:10]
		}
		displayContent = lipgloss.JoinVertical(lipgloss.Right, date, m.buttonStyle.Render("EDIT"), m.buttonStyle.Render("DELETE"), position)
		content := lipgloss.PlaceHorizontal(m.width-(lipgloss.Width(displayContent)+31), lipgloss.Left, m.entryContentStyle.Width(m.width-(lipgloss.Width(displayContent)+31)).Render(currEntry.content))
		displayContent = lipgloss.JoinHorizontal(lipgloss.Top, content, displayContent)
	} else {
		displayContent = lipgloss.PlaceHorizontal(m.width-31, lipgloss.Center, "NO "+tabsArr[m.tabCursor])
		displayContent = lipgloss.PlaceVertical(m.height-4, lipgloss.Center, displayContent)
	}
	isFocused = m.mainCursor > work && m.rightCursor == display
	displayContent = renderFocused(m.displayStyle, displayContent, isFocused)

	rightSide := lipgloss.JoinVertical(lipgloss.Left, headerContent, displayContent)

	s = lipgloss.JoinHorizontal(lipgloss.Top, s, rightSide)

	v = tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
