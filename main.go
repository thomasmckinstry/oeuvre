package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/donderom/bubblon"
	"github.com/thomasmckinstry/Bubbletea-Tutorial/Views"
	"golang.org/x/term"
)

var (
	width   int
	height  int
	rows    []table.Row
	columns []table.Column
)

type model struct {
	cursor    int
	currViews []string
	homeModel views.HomeModel
}

func initialModel() model {
	return model{
		currViews: make([]string, 2),
		homeModel: views.InitialHome(width, height),
		cursor:    0,
	}
}

func (m model) Init() tea.Cmd {
	m.currViews[0] = "home"
	m.currViews[1] = "add"
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.homeModel, cmd = m.homeModel.Update(msg)
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "K":
			m.homeModel, cmd = m.homeModel.Update(msg)
			cmds = append(cmds, cmd)
			if m.cursor > 0 {
				m.cursor--
			}
		case "J":
			m.homeModel, cmd = m.homeModel.Update(msg)
			cmds = append(cmds, cmd)
			if m.cursor < len(m.currViews) {
				m.cursor++
			}
		case "j", "k", "up", "down":
			if m.currViews[m.cursor] == "list" {
				m.homeModel, cmd = m.homeModel.Update(msg)
			}
			cmds = append(cmds, cmd)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	s := ""
	s += m.homeModel.View()

	// Send the UI for rendering
	return s
}

func main() {
	var err error
	width, height, err = term.GetSize(1)

	if err != nil {
		log.Fatal(err)
		return
	}
	mainModel := initialModel()
	controller, _ := bubblon.New(mainModel)
	program := tea.NewProgram(controller, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
