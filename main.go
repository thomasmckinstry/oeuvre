package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	tea "charm.land/bubbletea/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views"
	"github.com/thomasmckinstry/MediaLogger-TUI/utils"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
	"golang.org/x/term"
)

var (
	width  int
	height int
)

const (
	home int = iota
	add
	work
	confirm
	viewsCount
)

type model struct {
	cursor        int
	homeModel     *views.HomeModel
	addModel      *views.AddModel
	workPageModel *views.WorkPageModel
	confirmModel  *views.ConfirmModel
}

func initialModel() model {
	homeAddr := views.InitialHome(width, height)
	addAddr := views.InitialAddModel(width, height)
	workAddr := views.InitialWorkPage(width, height)
	confirmAddr := views.InitialConfirmModel(width, height)
	return model{
		confirmModel:  confirmAddr,
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
	case ConfirmationMsg:
		DebugLog("Got confirm msg: ", msg)
		m.cursor = ConfirmPage
		_, cmd = m.confirmModel.Update(msg)
	case DeleteWorkMsg:
		_, cmd = m.homeModel.Update(msg)
	case NewWorkMsg:
		_, cmd = m.homeModel.Update(msg)
	case WorkDetails:
		_, cmd = m.workPageModel.Update(msg)
	case ViewMsg:
		DebugLog("Main got ViewMsg: ", msg)
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
			case confirm:
				_, cmd = m.confirmModel.Update(msg)
			}
			cmds = tea.Batch(cmds, cmd)
		}
	}

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
	case confirm:
		view = m.confirmModel.View()
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

	utils.ReadConfig("config.yaml")

	mainModel := initialModel()
	program := tea.NewProgram(&mainModel)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
