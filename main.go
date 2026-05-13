package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"golang.org/x/term"
	"net/http"
	_ "net/http/pprof"
)

var (
	width  int
	height int
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
	case utils.NewWorkMsg:
		_, cmd = m.homeModel.Update(msg)
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

		default:
			switch m.currViews[m.cursor] {
			case "home":
				_, cmd = m.homeModel.Update(msg)
			case "add":
				_, cmd = m.addModel.Update(msg)
			}
			cmds = tea.Batch(cmds, cmd)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmds
}

func (m *model) View() tea.View {
	var view tea.View
	switch m.currViews[m.cursor] {
	case "home":
		view = m.homeModel.View()
	case "add":
		view = m.addModel.View()
	}

	// Send the UI for rendering
	return view
}

func main() {
	var err error
	width, height, err = term.GetSize(1)
	utils.CheckError("Failed to get terminal size: ", err)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		utils.CheckError("Failed to set up debug logging: ", err)
		defer func() {
			err = f.Close()
			utils.CheckError("Failed to close debug.log: ", err)
		}()
	}

	if len(os.Getenv("PPROF")) > 0 {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	mainModel := initialModel()
	program := tea.NewProgram(&mainModel)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
