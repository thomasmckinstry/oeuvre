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
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"golang.org/x/term"
	"net/http"
	_ "net/http/pprof"
)

var (
	width  int
	height int
)

const (
	home int = iota
	add
	work
	viewsCount
)

type model struct {
	cursor        int
	homeModel     *views.HomeModel
	addModel      *views.AddModel
	workPageModel *views.WorkPageModel
}

func initialModel() model {
	homeAddr := views.InitialHome(width, height)
	addAddr := views.InitialAddModel(width, height)
	workAddr := views.InitialWorkPage(width, height)
	return model{
		homeModel:     homeAddr,
		addModel:      addAddr,
		workPageModel: workAddr,
		cursor:        home,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case NewWorkMsg:
		_, cmd = m.homeModel.Update(msg)
	case WorkDetails:
		_, cmd = m.workPageModel.Update(msg)
	case ViewMsg:
		if len(os.Getenv("DEBUG")) > 0 {
			log.Println("main received ViewMsg for ", int(msg))
		}
		m.addModel, _ = m.addModel.Update(msg)
		m.cursor = int(msg)
	case tea.WindowSizeMsg:
		switch m.cursor {
		case home:
			_, cmd = m.homeModel.Update(msg)
		case work:
			_, cmd = m.workPageModel.Update(msg)
		case add:
			_, cmd = m.addModel.Update(msg)
		}
		cmds = tea.Batch(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		default:
			switch m.cursor {
			case home:
				_, cmd = m.homeModel.Update(msg)
			case add:
				_, cmd = m.addModel.Update(msg)
			case work:
				_, cmd = m.workPageModel.Update(msg)
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
	switch m.cursor {
	case home:
		view = m.homeModel.View()
	case add:
		view = m.addModel.View()
	case work:
		view = m.workPageModel.View()
	}

	// Send the UI for rendering
	return view
}

func main() {
	var err error
	width, height, err = term.GetSize(1)
	CheckError("Failed to get terminal size: ", err)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		CheckError("Failed to set up debug logging: ", err)
		defer func() {
			err = f.Close()
			CheckError("Failed to close debug.log: ", err)
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
