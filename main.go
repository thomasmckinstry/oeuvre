package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
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
	homeModel *views.HomeModel
	addModel  *views.AddModel
}

func initialModel() model {
	homeAddr := views.InitialHome(width, height)
	addAddr := views.InitialAddModel(22)
	return model{
		currViews: make([]string, 2),
		homeModel: homeAddr,
		addModel:  addAddr,
		cursor:    0,
	}
}

func (m *model) Init() tea.Cmd {
	m.currViews[0] = "home"
	m.currViews[1] = "add"
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case views.ViewMsg:
		if len(os.Getenv("DEBUG")) > 0 {
			log.Println("main received ViewMsg for ", int(msg))
		}
		m.cursor = int(msg)
	case tea.WindowSizeMsg:
		_, cmd = m.homeModel.Update(msg)
		cmds = tea.Batch(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "K", "L", "H", "J":
			if m.currViews[m.cursor] == "home" {
				_, cmd = m.homeModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.currViews[m.cursor] == "add" {
				_, cmd = m.addModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case "j", "k", "up", "down", "left", "right", "h", "l":
			if m.currViews[m.cursor] == "home" {
				_, cmd = m.homeModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.currViews[m.cursor] == "add" {
				_, cmd = m.addModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
			cmds = tea.Batch(cmds, cmd)
		default:
			if m.currViews[m.cursor] == "home" {
				_, cmd = m.homeModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.currViews[m.cursor] == "add" {
				_, cmd = m.addModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmds
}

func (m *model) View() tea.View {
	var view tea.View
	if len(os.Getenv("DEBUG")) > 0 {
		log.Println("main cursor at:", m.cursor)
	}
	if m.currViews[m.cursor] == "home" {
		view = m.homeModel.View()
	} else if m.currViews[m.cursor] == "add" {
		view = m.addModel.View()
	}

	// Send the UI for rendering
	return view
}

func main() {
	var err error
	width, height, err = term.GetSize(1)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if len(os.Getenv("DEBUG")) > 0 {
		log.Println("Successfully initialized connection to database")
	}
	if err != nil {
		log.Fatal(err)
		return
	}

	mainModel := initialModel()
	program := tea.NewProgram(&mainModel)
	if len(os.Getenv("DEBUG")) > 0 {
		log.Println("Successfully Initialized Program")
	}
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
