package partials

import (
	"log"
	"os"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	database "github.com/thomasmckinstry/MediaLogger-TUI/db"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
)

type Form interface {
	GetContents() []string
}

type filterKeyMap struct {
	Nav     key.Binding
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
	Unfocus key.Binding
}

const (
	title int = iota
	tags
	medium
	status
	enter
)

var (
	filterStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	headerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	textinputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00")).
			BorderLeft(true)
	enterStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#6E3F00"))
)

var filterDefaultMap = filterKeyMap{
	Nav: key.NewBinding(
		key.WithKeys("H", "K", "J", "L"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
	),
	Unfocus: key.NewBinding(
		key.WithKeys("esc"),
	),
}

type FilterMsg [][]string

// TODO: I can probably sub out most of this file for a huh? component
// Don't want to do that because it's more hands off
// TODO: Need a different way to index into the components because of different Model types.
type FilterModel struct {
	headerText string
	focused    bool
	cursor     int
	height     int
	width      int
	forms      []tea.Model
}

func InitialFilter(height int) FilterModel {
	width := 14
	tagSuggestions := []string{}

	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM tags_table`)
	CheckError("Failed to query tags from db: ", err)
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		CheckError("Failed to scan tags: ", err)
		tagSuggestions = append(tagSuggestions, tag)
	}
	err = rows.Close()
	CheckError("Failed to close tags query: ", err)

	titleSuggestions := []string{}

	rows, err = db.Query(`SELECT title FROM works`)
	CheckError("Failed to query tags from database: ", err)
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		CheckError("Failed to scan titles: ", err)
		titleSuggestions = append(titleSuggestions, tag)
	}
	err = rows.Close()
	CheckError("Failed to close title query: ", err)

	titleInput := components.InitialTextInput(width, "Title", "{ title }", titleSuggestions)
	tagsInput := components.InitialInput(5, "{ tag }", "Tag", width, false, tagSuggestions)
	mediums := utils.Config.MediaOptions // TODO: Query the db for this.
	mediumInput := components.InitialCheckbox(mediums, "Medium", width)
	statuses := utils.Config.StatusOptions // TODO: Query the db for this.
	statusInput := components.InitialCheckbox(statuses, "Status", width)

	forms := []tea.Model{&titleInput, &tagsInput, &mediumInput, &statusInput}

	return FilterModel{
		headerText: "Filter",
		focused:    false,
		cursor:     0,
		height:     height,
		width:      width,
		forms:      forms,
	}
}

func (m *FilterModel) Init() tea.Cmd {
	return nil
}

func (m *FilterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		filterStyle = filterStyle.Height(msg.Height - (7))
		m.height = msg.Height - 7
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, filterDefaultMap.Confirm):
			if m.cursor == len(m.forms) { // Cursor on enter button
				var contents [][]string
				var content []string
				for _, form := range m.forms {
					switch form := form.(type) {
					case *components.TextInputModel:
						content = []string{form.GetContents()}
					case *components.TagInputModel:
						content = form.GetContents()
					case *components.CheckboxModel:
						content = form.GetContents()
					}
					contents = append(contents, content)
				}
				if len(os.Getenv("DEBUG")) > 0 {
					log.Println("Filtering for: ", contents)
				}
				cmd = func() tea.Msg { return FilterMsg(contents) }
				return m, cmd
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			m.focused = true
		case key.Matches(msg, filterDefaultMap.Unfocus):
			_, cmd = m.forms[m.cursor].Update(msg)
			m.focused = false
		case key.Matches(msg, filterDefaultMap.Down):
			if m.cursor > len(m.forms)-1 {
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			msg, ok := cmd().(NavMsg)
			if m.cursor < len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
				_, cmd = m.forms[m.cursor].Update(msg)
			} else if m.cursor >= len(m.forms)-1 && ok && bool(msg) {
				m.cursor++
			}
		case key.Matches(msg, filterDefaultMap.Up):
			if m.cursor == len(m.forms) {
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
				break
			}
			_, cmd = m.forms[m.cursor].Update(msg)
			msg, ok := cmd().(NavMsg)
			if m.cursor > 0 && ok && bool(msg) {
				m.cursor--
				_, cmd = m.forms[m.cursor].Update(msg)
			}
		case key.Matches(msg, filterDefaultMap.Nav):
			cmd = func() tea.Msg { return NavMsg(!m.focused) }
			cmds = tea.Batch(cmds, cmd)
			fallthrough
		default:
			_, cmd = m.forms[m.cursor].Update(msg)
		}
	}
	return m, cmds
}

func (m *FilterModel) View() tea.View {
	var c *tea.Cursor
	//header:
	s := headerStyle.Render(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.headerText))

	for i, form := range m.forms {
		s += "\n"
		formView := form.View()
		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s) + 2
			c.X += 1
		}
		isFocused := m.cursor == i
		s = lipgloss.JoinVertical(lipgloss.Center, s, RenderFocused(textinputStyle, formView.Content, isFocused))
	}
	isFocused := m.cursor == enter
	enter := RenderFocused(enterStyle, lipgloss.PlaceHorizontal(15, lipgloss.Center, "ENTER"), isFocused)
	s = lipgloss.JoinVertical(lipgloss.Left, s, enter)

	v := tea.NewView(filterStyle.Render(s))
	v.Cursor = c
	return v
}
