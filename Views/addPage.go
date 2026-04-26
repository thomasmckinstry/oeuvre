package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"encoding/json"
	"github.com/thomasmckinstry/Bubbletea-Tutorial/Views/Components"
	database "github.com/thomasmckinstry/Bubbletea-Tutorial/db"
	"log"
	"os"
)

type AddModel struct {
	headerText     string
	focused        bool
	cursor         int
	height         int
	width          int
	forms          []tea.Model
	status         []string
	tags           []string
	style          lipgloss.Style
	headerStyle    lipgloss.Style
	textinputStyle lipgloss.Style
	enterStyle     lipgloss.Style

	errorMsg string
}

func InitialAddModel(width int) *AddModel {
	title := components.InitialTextInput(width, "Title", "{ title }", nil)
	mediums := []string{"Movie", "Book", "Show", "Anime", "Manga", "Comic", "Show", "Animated", "Live Action"} // TODO: Query the db for this.
	medium := components.InitialCheckbox(mediums, "Medium", width)
	statuses := []string{"Pending", "Started", "Hiatus", "Completed", "Dropped"} // TODO: Query the db for this.
	status := components.InitialCheckbox(statuses, "Status", width)

	var tagSuggestions []string
	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM tags_table`)
	if err != nil {
		log.Fatal("Failed to query tags from database: ", err)
	}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if len(os.Getenv("DEBUG")) > 0 {
			log.Println("Got tag of length: ", len(tag))
		}
		if err != nil {
			log.Fatal("Failed to scan tags: ", err)
		}
		tagSuggestions = append(tagSuggestions, tag)
	}
	rows.Close()
	tagSuggestions = []string{"test"}

	tags := components.InitialInput(20, "{ tags }", "Tags", width, false, tagSuggestions)
	forms := []tea.Model{&title, &tags, &medium, &status}
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
			BorderTop(true),
		enterStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")),
	}
}

func (m *AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (*AddModel, tea.Cmd) {
	if len(os.Getenv("DEBUG")) > 0 {
		log.Println("AddPage got: ", msg)
	}
	var cmds tea.Cmd
	var cmd tea.Cmd
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Height(msg.Height - (7))
		m.height = msg.Height - 7
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.cursor == len(m.forms) {
				var contents []string
				var content []byte
				var err error
				for _, form := range m.forms {
					switch form := form.(type) {
					case *components.TextInputModel:
						content = []byte(form.GetContents())
					case *components.TagInputModel:
						content, err = json.Marshal(form.GetContents()) // TODO: Marshal this to JSON
					case *components.CheckboxModel:
						content, err = json.Marshal(form.GetContents()) // TODO: Marshal this to JSONA
					}
					if err != nil {
						log.Fatal("Failed to marshal input data to JSON: ", err)
					}
					contents = append(contents, string(content))
				}
				if len(os.Getenv("DEBUG")) > 0 {
					log.Println(contents)
				}
				db := database.GetDB() //  TODO: Move this shit (whole block) into a function
				id := database.GetWorkID(db)
				query, err := db.Prepare(`
					INSERT INTO works (date_added, title, media_type, work_status, tags, year_released, work_id)
					VALUES (?, ?, ?, ?, ?, ?, ?)
				`)
				if err != nil {
					log.Fatal("Failed to prepare insert statement: ", err)
				}
				_, err = query.Exec("test_date", contents[0], contents[2], contents[1][0], contents[3], 0, id)
				if err != nil {
					log.Fatal("Failed to insert to works table: ", err)
				}
				query.Close()
				cmd = func() tea.Msg { return ViewMsg(0) }
				cmds = tea.Batch(cmds, cmd)
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			m.focused = true
		case "esc":
			_, cmd = m.forms[m.cursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			m.focused = false
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
			c.Y += lipgloss.Height(s) + 2 // TODO: Make the + 2 not hardcoded
			c.X += 1
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
