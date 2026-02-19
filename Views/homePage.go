package views

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	partials "github.com/thomasmckinstry/Bubbletea-Tutorial/Views/Partials"
)

var (
	width  int
	height int

// rows    []table.Row
// columns []table.Column
)

type HomeModel struct {
	cursor       int
	currPartials []string
	listModel    partials.ListModel
	addModel     partials.AddModel
	filterModel  partials.FilterModel
	sortModel    partials.SortModel
}

func InitialHome(width int, height int) HomeModel {
	var columns = []table.Column{ // TODO: Remove this
		{Title: "Title", Width: 30},
		{Title: "Medium", Width: 20},
		{Title: "Status", Width: 15},
		{Title: "Genre", Width: 14},
	}

	var rows = []table.Row{ // TODO: Remove this
		{"I am Your Beast", "Game", "Completed", "Action"},
		{"One Battle After Another", "Movie, Live Action", "Pending", "Action"},
	}

	return HomeModel{
		currPartials: make([]string, 2),
		listModel:    partials.InitialList(width, height, columns, rows),
		addModel:     partials.InitialAdd(), // height = 1 Note: I think each side of the border adds ~1.5
		filterModel:  partials.InitialFilter(height - (10)),
		sortModel:    partials.InitialSort(3),
		cursor:       0,
	}
}

func (m HomeModel) Init() tea.Cmd {
	m.currPartials[0] = "add"
	m.currPartials[1] = "filter"
	m.currPartials[2] = "sort"
	m.currPartials[3] = "list"
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.listModel, cmd = m.listModel.Update(msg)
		m.filterModel, cmd = m.filterModel.Update(msg)
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "K":
			cmds = append(cmds, cmd)
			m.listModel, cmd = m.listModel.Update(msg)
			cmds = append(cmds, cmd)
			if m.cursor > 0 {
				m.cursor--
			}
		case "J":
			cmds = append(cmds, cmd)
			m.listModel, cmd = m.listModel.Update(msg)
			cmds = append(cmds, cmd)
			if m.cursor < len(m.currPartials) {
				m.cursor++
			}
		case "j", "k", "up", "down":
			if m.currPartials[m.cursor] == "list" {
				m.listModel, cmd = m.listModel.Update(msg)
			}
			cmds = append(cmds, cmd)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m HomeModel) View() string {
	s := ""
	sidebar := lipgloss.JoinVertical(lipgloss.Center, m.addModel.View(), m.filterModel.View(), m.sortModel.View())
	list := m.listModel.View()
	s = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, list)

	// Send the UI for rendering
	return s
}
