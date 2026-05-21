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
	notesCursor       int
	reviewsCursor     int
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
	arrowStyle        lipgloss.Style
}

type workKeyMap struct {
	TopLevelUp    key.Binding
	TopLevelDown  key.Binding
	TopLevelLeft  key.Binding
	TopLevelRight key.Binding
	Confirm       key.Binding
	Up            key.Binding
	Down          key.Binding
	Left          key.Binding
	Right         key.Binding
	Exit          key.Binding
}

const (
	work int = iota
	tabs
	add
	del
	mainCursorCount
)

const (
	note int = iota
	review
)

const (
	header int = iota
	display
	rightCursorCount
)

var writeHeader = []string{"NOTE", "REVIEW"}

var defaultWorkMap = workKeyMap{
	TopLevelUp:    key.NewBinding(key.WithKeys("K")),
	TopLevelDown:  key.NewBinding(key.WithKeys("J")),
	TopLevelLeft:  key.NewBinding(key.WithKeys("H")),
	TopLevelRight: key.NewBinding(key.WithKeys("L")),
	Left:          key.NewBinding(key.WithKeys("h", "left")),
	Right:         key.NewBinding(key.WithKeys("l", "right")),
	Up:            key.NewBinding(key.WithKeys("k", "up")),
	Down:          key.NewBinding(key.WithKeys("j", "down")),
	Confirm:       key.NewBinding(key.WithKeys("enter")),
	Exit:          key.NewBinding(key.WithKeys("esc")),
}

func (m *WorkPageModel) resetCursors() {
	m.tabCursor = 0
	m.mainCursor = 0
	m.entryCursor = 0
	m.rightCursor = 0
	m.writingCursor = 0
	m.displayCursor = 0

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
			BorderForeground(Unfocused),
		tabsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			MarginLeft(1).
			MarginRight(1).
			BorderForeground(Unfocused),
		detailsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			MarginRight(1).
			BorderForeground(Unfocused),
		buttonStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(Unfocused),
		displayStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderTop(true).
			BorderForeground(Unfocused),
		entryContentStyle: lipgloss.NewStyle().
			MarginRight(1),
		arrowStyle: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(Unfocused),
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
						query, err = db.Prepare(`REPLACE INTO notes (date_added, note_text, work_id, note_id) VALUES (?, ?, ?, ?)`)
						m.notes[m.notesCursor] = entry{content: content, date: date, id: m.notes[m.notesCursor].id}
					} else {
						query, err = db.Prepare(`REPLACE INTO reviews (date_added, review_text, work_id, review_id) VALUES (?, ?, ?, ?)`)
						m.reviews[m.reviewsCursor] = entry{content: content, date: date, id: m.reviews[m.reviewsCursor].id}
					}

				}
				CheckError("Failed to prepare statement to insert to DB: ", err)
				if m.writingMode == "ADD" {
					query.Exec(date, content, m.currWorkId)
				} else {
					switch m.tabCursor {
					case note:
						query.Exec(date, content, m.currWorkId, m.notes[m.notesCursor].id)
					case review:
						query.Exec(date, content, m.currWorkId, m.reviews[m.reviewsCursor].id)
					}
				}
				query.Close()
				m.textArea.Reset()
				m.writingCursor = 0

			} else if m.rightCursor == display && m.mainCursor != work && m.entryCursor == 1 {
				db := database.GetDB()
				var (
					query *sql.Stmt
					index int
					err   error
				)
				switch m.tabCursor {
				case note:
					query, err = db.Prepare("DELETE FROM notes WHERE note_id = ?")
					index = m.notesCursor
					m.notes = append(m.notes[:index], m.notes[:index+1]...)
					if m.notesCursor > 0 {
						m.notesCursor--
					}
				case review:
					query, err = db.Prepare("DELETE FROM reviews WHERE review_id = ?")
					index = m.reviewsCursor
					m.reviews = append(m.reviews[:index], m.reviews[:index+1]...)
					if m.reviewsCursor > 0 {
						m.reviewsCursor--
					}
				}
				CheckError("Failed to prep delete entry query: ", err)
				_, err = query.Exec(index)
			} else if m.rightCursor == display && m.mainCursor != work && m.entryCursor == 0 {
				m.writing = true
				m.writingMode = "EDIT"
				if m.tabCursor == 0 {
					m.textArea.SetValue(m.notes[m.entryCursor].content)
				} else {
					m.textArea.SetValue(m.reviews[m.entryCursor].content)
				}
			} else if m.rightCursor != display && m.mainCursor == del {
				DebugLog("Deleting work: ", m.currWorkId)
				db := database.GetDB()
				query, err := db.Prepare(`DELETE FROM works WHERE work_id = ?`)
				CheckError("Failed to prepare delete work query: ", err)
				_, err = query.Exec(m.currWorkId)
				CheckError("Failed to delete work from db: ", err)
				m.resetCursors()
				cmds = tea.Batch(cmds, func() tea.Msg { return ViewMsg(0) })
				cmds = tea.Batch(cmds, func() tea.Msg { return DeleteWorkMsg(m.currWorkId) })
			} else if m.mainCursor == work {
				m.work, cmd = m.work.Update(msg)
				var ok bool
				var workMsg []string
				if cmd != nil {
					workMsg, ok = cmd().(NewWorkMsg)
					cmds = tea.Batch(cmds, cmd)
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
				cmd = func() tea.Msg { return DeleteWorkMsg(m.currWorkId) }
				DebugLog("Created DeleteWorkMsg: ", m.currWorkId)
			}
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultWorkMap.Exit):
			if m.mainCursor == work {
				_, cmd = m.work.Update(msg)
				if cmd != nil {
					_, ok := cmd().(ViewMsg)
					if ok && m.focused {
						cmd = nil
						m.focused = false
					} else {
						m.work.ClearComponents()
					}
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
				m.resetCursors()
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
			} else if m.rightCursor == display && m.mainCursor != work {
				switch m.tabCursor {
				case note:
					if m.notesCursor > 0 {
						m.notesCursor--
					}
				case review:
					if m.reviewsCursor > 0 {
						m.reviewsCursor--
					}
				}
			}
		case key.Matches(msg, defaultWorkMap.Right):
			if m.mainCursor == tabs && m.rightCursor == header && !m.writing {
				m.tabCursor = 1
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.rightCursor == display && m.mainCursor != work {
				switch m.tabCursor {
				case note:
					if m.notesCursor < len(m.notes)-1 {
						m.notesCursor++
					}
				case review:
					if m.reviewsCursor < len(m.reviews)-1 {
						m.reviewsCursor++
					}
				}
			}
		case key.Matches(msg, defaultWorkMap.Up):
			if m.mainCursor == work && m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.rightCursor == 1 && m.entryCursor == 1 {
				m.entryCursor--
			}
		case key.Matches(msg, defaultWorkMap.Down):
			if m.mainCursor == work && m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.rightCursor == display && m.entryCursor == 0 {
				m.entryCursor++
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
		textarea := RenderFocused(m.tabsStyle, m.textArea.View(), m.writingCursor == 0)
		button := RenderFocused(m.buttonStyle, "CONFIRM", m.writingCursor == 1)

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
	details = RenderFocused(m.tabsStyle, details, isFocused)

	isSelected := m.mainCursor == work
	details = RenderFocused(m.detailsStyle.Height(m.height), details, isSelected)

	s = lipgloss.JoinHorizontal(lipgloss.Top, s, details)

	headerContent := "DELETE"
	isFocused = m.mainCursor == del && m.rightCursor == header
	headerContent = RenderFocused(m.buttonStyle, headerContent, isFocused)
	isFocused = m.mainCursor == add && m.rightCursor == header
	headerContent = lipgloss.JoinHorizontal(lipgloss.Top, RenderFocused(m.buttonStyle, "ADD", isFocused), headerContent)

	renderedTabs := []string{}
	tabsArr := []string{"NOTES", "REVIEWS"}
	for i, tab := range tabsArr {
		tab = lipgloss.PlaceHorizontal(9, lipgloss.Center, tab)
		isFocused = m.tabCursor == i
		renderedTabs = append(renderedTabs, RenderFocused(m.tabStyle, tab, isFocused))
	}

	tabsContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs[0], " ", renderedTabs[1])
	isFocused = m.mainCursor == tabs && m.rightCursor == header
	tabsContent = RenderFocused(m.tabsStyle, tabsContent, isFocused)

	headerContent = lipgloss.JoinHorizontal(lipgloss.Top, tabsContent, headerContent)

	var displayContent, position string
	var currEntry entry
	var hasEntry bool
	if m.tabCursor == 0 && len(m.notes) > 0 {
		currEntry = m.notes[m.notesCursor]
		hasEntry = true
		position = lipgloss.PlaceVertical(m.height-11, lipgloss.Bottom, strconv.Itoa(m.notesCursor+1)+"/"+strconv.Itoa(len(m.notes)))
	} else if m.tabCursor == 1 && len(m.reviews) > 0 {
		currEntry = m.reviews[m.reviewsCursor]
		hasEntry = true
		position = lipgloss.PlaceVertical(m.height-11, lipgloss.Bottom, strconv.Itoa(m.reviewsCursor+1)+"/"+strconv.Itoa(len(m.reviews)))
	}

	if hasEntry {
		var date string
		if currEntry.date != "" {
			date = currEntry.date[:10]
		}
		isFocused = m.rightCursor == display && m.entryCursor == 0
		editBtn := RenderFocused(m.buttonStyle, "EDIT", isFocused)
		isFocused = m.rightCursor == display && m.entryCursor == 1
		deleteBtn := RenderFocused(m.buttonStyle, "DELETE", isFocused)
		displayContent = lipgloss.JoinVertical(lipgloss.Right, date, editBtn, deleteBtn, position)
		content := m.entryContentStyle.Width(m.width - (lipgloss.Width(displayContent) + 37)).Render(currEntry.content)
		displayContent = lipgloss.JoinHorizontal(lipgloss.Top, content, displayContent)
		leftArrow := m.arrowStyle.Render(lipgloss.PlaceVertical(lipgloss.Height(displayContent), lipgloss.Center, "<"))
		rightArrow := m.arrowStyle.Render(lipgloss.PlaceVertical(lipgloss.Height(displayContent), lipgloss.Center, ">"))
		displayContent = lipgloss.JoinHorizontal(lipgloss.Center, leftArrow, displayContent, rightArrow)
	} else {
		displayContent = lipgloss.PlaceHorizontal(m.width-31, lipgloss.Center, "NO "+tabsArr[m.tabCursor])
		displayContent = lipgloss.PlaceVertical(m.height-4, lipgloss.Center, displayContent)
	}
	isFocused = m.mainCursor > work && m.rightCursor == display
	displayContent = RenderFocused(m.displayStyle, displayContent, isFocused)

	rightSide := lipgloss.JoinVertical(lipgloss.Left, headerContent, displayContent)

	s = lipgloss.JoinHorizontal(lipgloss.Top, s, rightSide)

	v = tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
