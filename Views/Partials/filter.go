package partials

import (
	"log"
	"os"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/Bubbletea-Tutorial/Views/Components"
	database "github.com/thomasmckinstry/Bubbletea-Tutorial/db"
)

type Form interface {
	GetContents() []string
}

// TODO: I can probably sub out most of this file for a huh? component
// Don't want to do that because it's more hands off
// TODO: Need a different way to index into the components because of different Model types.
type FilterModel struct {
	headerText     string
	selected       bool // Indicates if the cursor is interacting with Filter
	focused        bool
	cursor         int
	height         int
	forms          []tea.Model // Can I get this to use pointers to the actual models? I think right now I'm copying them
	status         []string
	tags           []string
	style          lipgloss.Style
	headerStyle    lipgloss.Style
	textinputStyle lipgloss.Style
	enterStyle     lipgloss.Style

	errorMsg string
}

// TODO: This should be a utils
func (m FilterModel) toggleBorder() lipgloss.Style {
	if m.selected == true {
		return m.style.BorderForeground(lipgloss.Color("#6E3F00"))
	}
	return m.style.BorderForeground(lipgloss.Color("#D17600"))
}

func InitialFilter(height int) FilterModel {
	tagSuggestions := []string{}

	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM tags_table`)
	if err != nil {
		log.Fatal("Failed to query tags from database: ", err)
	}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Fatal("Failed to scan tags: ", err)
		}
		tagSuggestions = append(tagSuggestions, tag)
	}

	titleSuggestions := []string{}

	rows, err = db.Query(`SELECT title FROM works`)
	if err != nil {
		log.Fatal("Failed to query tags from database: ", err)
	}
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			log.Fatal("Failed to scan titles: ", err)
		}
		titleSuggestions = append(titleSuggestions, tag)
	}

	titleInput := components.InitialTextInput(14, "Title", "{ title }", titleSuggestions)
	tagsInput := components.InitialInput(5, "{ tag }", "Tag", 14, false, tagSuggestions)
	mediums := []string{"Movie", "Book", "Show", "Anime", "Manga", "Comic", "Show", "Animated", "Live Action"} // TODO: Query the db for this.
	mediumInput := components.InitialCheckbox(mediums, "Medium", 14)
	statuses := []string{"Pending", "Started", "Hiatus", "Completed", "Dropped"} // TODO: Query the db for this.
	statusInput := components.InitialCheckbox(statuses, "Status", 14)

	forms := []tea.Model{&titleInput, &tagsInput, &mediumInput, &statusInput}

	return FilterModel{
		headerText: "Filter",
		selected:   false,
		focused:    false,
		cursor:     0,
		height:     height,
		forms:      forms,
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderTop(true).
			Width(18).
			Height(height).
			Align(lipgloss.Center),
		headerStyle: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(16),
		textinputStyle: lipgloss.NewStyle().
			Width(17).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderLeft(true),
		enterStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")),
	}
}

func (m *FilterModel) Init() tea.Cmd {
	return nil
}

func (m *FilterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				var contents [][]string
				for _, form := range m.forms {
					switch form := form.(type) {
					case *components.TextInputModel:
						contents = append(contents, form.GetContents())
					case *components.TagInputModel:
						contents = append(contents, form.GetContents())
					case *components.CheckboxModel:
						contents = append(contents, form.GetContents())
					}
				}
				if len(os.Getenv("DEBUG")) > 0 {
					log.Println(contents)
				}
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			m.focused = true
		case "esc":
			_, cmd = m.forms[m.cursor].Update(msg)
			m.focused = false
		case "L", "H", "J", "K":
			m.style = m.toggleBorder()
			m.selected = !m.selected
		case "j", "down":
			if m.cursor > len(m.forms)-1 {
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			msg, ok := cmd().(components.NavMsg)
			if m.cursor < len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				_, cmd = m.forms[m.cursor].Update(msg)
			} else if m.cursor >= len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#D17600"))
			}
		case "k", "up":
			if m.cursor == len(m.forms) {
				m.enterStyle = m.enterStyle.BorderForeground(lipgloss.Color("#6E3F00"))
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			msg, ok := cmd().(components.NavMsg)
			if m.cursor > 0 && ok && bool(msg) {
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
			}
		default:
			_, cmd = m.forms[m.cursor].Update(msg)
		}
	}
	return m, cmd
}

// TODO: Add styling to make it clear that a textbox is selected.
// TODO: Iterate over m.forms instead of having a bunch of different conditional blocks
func (m *FilterModel) View() tea.View {
	var c *tea.Cursor
	//header:
	s := m.headerStyle.Render(m.headerText)

	for i, form := range m.forms {
		s += "\n"
		formView := form.View()
		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s) + 2 // TODO: Make the + 2 not hardcoded
			c.X += 1
		}
		if i == m.cursor && m.selected {
			s = lipgloss.JoinVertical(lipgloss.Left, s, m.textinputStyle.BorderForeground(lipgloss.Color("#D17600")).Render(formView.Content))
		} else {
			s = lipgloss.JoinVertical(lipgloss.Left, s, m.textinputStyle.Render(formView.Content))
		}
	}
	enter := m.enterStyle.Render(lipgloss.PlaceHorizontal(15, lipgloss.Center, "ENTER"))
	enter = lipgloss.PlaceVertical(m.height-lipgloss.Height(s)-1, lipgloss.Bottom, enter)
	s = lipgloss.JoinVertical(lipgloss.Left, s, enter)

	v := tea.NewView(m.style.Render(s))
	v.Cursor = c
	return v
}
