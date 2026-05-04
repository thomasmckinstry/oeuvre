package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	partials "github.com/thomasmckinstry/MediaLogger-TUI/Views/Partials"

	"log"
	"os"
)

var height int

type HomeModel struct {
	sidebarCursor int
	mainCursor    int
	sidebarViews  []tea.Model
	listModel     tea.Model
}

type ViewMsg int

func InitialHome(width int, height int) *HomeModel {
	list := partials.InitialList(width-19, height)
	add := partials.InitialAdd() // height = 1 Note: I think each side of the border adds 1
	filter := partials.InitialFilter(height - (7))
	sort := partials.InitialSort(3)

	sidebarList := []tea.Model{}               //make([]tea.Model, 3)
	sidebarList = append(sidebarList, &add)    //[0] = add
	sidebarList = append(sidebarList, &filter) //[1] = filter
	sidebarList = append(sidebarList, &sort)   //[2] = sort

	return &HomeModel{
		sidebarViews:  sidebarList,
		listModel:     &list,
		sidebarCursor: 0,
		mainCursor:    0,
	}
}

func (m *HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (*HomeModel, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		_, cmd = m.listModel.Update(msg)
		cmds = tea.Batch(cmds, cmd)
		_, cmd = m.sidebarViews[1].Update(msg)
		cmds = tea.Batch(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		//TODO: See if  I can clean up the homepage nav so I can just have it in one condition
		case "K":
			if m.mainCursor == 0 && m.sidebarCursor > 0 {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				m.sidebarCursor--
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case "J":
			if m.mainCursor == 0 && m.sidebarCursor < 2 {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				m.sidebarCursor++
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case "H":
			if m.mainCursor > 0 {
				m.mainCursor--
				m.listModel, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case "L":
			if m.mainCursor < 1 {
				m.mainCursor++
				m.listModel, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case "j", "k", "up", "down":
			if m.mainCursor == 1 {
				m.listModel, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case "enter":
			if m.sidebarCursor == 0 && m.mainCursor == 0 {
				if len(os.Getenv("DEBUG")) > 0 {
					log.Println("homePage sending AddMsg")
				}
				cmds = tea.Batch(cmds, func() tea.Msg { return (ViewMsg(1)) })
			} else {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				if cmd != nil {
					msg, ok := cmd().(partials.FilterMsg)
					if ok {
						_, cmd = m.listModel.Update(msg)
						cmds = tea.Batch(cmds, cmd)
					}
				}
			}
		default:
			_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			if cmd != nil {
				msg, ok := cmd().(partials.SortMsg)
				if ok {
					_, cmd = m.listModel.Update(msg)
					cmds = tea.Batch(cmds, cmd)
				}
			}
		}
	}

	return m, cmds
}

func (m *HomeModel) View() tea.View {
	var c *tea.Cursor
	s := ""
	sidebarContent := []string{}
	for _, form := range m.sidebarViews {
		formView := form.View()
		sidebarContent = append(sidebarContent, formView.Content)

		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s)
		}
	}
	sidebar := lipgloss.JoinVertical(lipgloss.Center, sidebarContent...)

	list := m.listModel.View()
	s = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, list.Content)

	// Send the UI for rendering
	view := tea.NewView(s)
	view.Cursor = c
	view.AltScreen = true
	return view
}
