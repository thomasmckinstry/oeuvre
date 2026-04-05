package partials

import (
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/Bubbletea-Tutorial/Views/Components"
)

// TODO: I can probably sub out most of this file for a huh? component
// Don't want to do that because it's more hands off
// TODO: Need a different way to index into the components because of different Model types.
type FilterModel struct {
	headerText     string
	titleInput     tea.Model
	genreInput     tea.Model
	themeInput     tea.Model // TODO: These need additional displays for previous entries
	selected       bool      // Indicates if the cursor is interacting with Filter
	focused        bool
	cursor         int
	forms          []tea.Model // Can I get this to use pointers to the actual models? I think right now I'm copying them
	status         []string
	genres         []string
	themes         []string
	style          lipgloss.Style
	headerStyle    lipgloss.Style
	textinputStyle lipgloss.Style

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
	titleInput := components.InitialInput(3, "", "Title", 14, true)
	genreInput := components.InitialInput(3, "", "Genre", 14, false)
	themeInput := components.InitialInput(3, "", "Theme", 14, false)

	//status := []string{"Completed", "In Progress", "Started", "Pending", "Dropped"}
	forms := []tea.Model{titleInput, genreInput, themeInput} // TODO: Figure out how to have null pointers to each form
	// forms is an array of all the forms that make up the filter box.
	// This is so I can index into each one as I navigate with the keyboard

	return FilterModel{
		headerText: "Filter",
		titleInput: titleInput,
		genreInput: genreInput,
		selected:   false,
		focused:    false,
		cursor:     0,
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
			MarginTop(1).
			Width(16),
	}
}

func (m FilterModel) Init() tea.Cmd {
	return nil
}

func (m FilterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.errorMsg = ""

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.style = m.style.Height(msg.Height - (7))
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			m.focused = true
		case "esc":
			m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			m.focused = false
		case "L", "H", "J", "K":
			m.style = m.toggleBorder()
			m.selected = !m.selected
		case "j", "down": // TODO: Make these check for focused inputs before moving the cursor
			m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			msg, ok := cmd().(components.NavMsg)
			if m.cursor < len(m.forms)-1 && ok && bool(msg) { //!m.focused {
				m.cursor++
				m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			}
		case "k", "up":
			m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			msg, ok := cmd().(components.NavMsg)
			if m.cursor > 0 && ok && bool(msg) {
				m.cursor--
				m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			}
		default:
			/*if field, ok := m.forms[m.cursor].(textinput.Model); ok {
				if field.Focused() {
					m.forms[m.cursor], cmd = field.Update(msg)
				}
			}*/
			//if m.focused {
			m.forms[m.cursor], cmd = m.forms[m.cursor].Update(msg)
			//}
		}
	}
	return m, cmd
}

// TODO: Add styling to make it clear that a textbox is selected.
// TODO: Iterate over m.forms instead of having a bunch of different conditional blocks
func (m FilterModel) View() tea.View {
	var c *tea.Cursor
	//header:
	s := m.headerStyle.Render(m.headerText)

	for _, form := range m.forms {
		formView := form.View()
		s = lipgloss.JoinVertical(lipgloss.Left, s, m.textinputStyle.Render(formView.Content))
		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s)
		}
	}

	/*if m.focused {
		s += "\nfocused"
	}
	if m.selected {
		s += "\nselected"
	}*/

	//s = lipgloss.JoinVertical(lipgloss.Left, s, m.errorMsg)
	v := tea.NewView(m.style.Render(s))
	v.Cursor = c
	return v
}
