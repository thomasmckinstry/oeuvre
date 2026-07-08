package views

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	components "github.com/thomasmckinstry/ouevre/Views/Components"
	partials "github.com/thomasmckinstry/ouevre/Views/Partials"
	. "github.com/thomasmckinstry/ouevre/utils"

	"log"
	"os"
)

const (
	addBtn int = iota
	filter
	sort
	numSidebarForms
)

const (
	sidebar int = iota
	list
	numPartialsHome
)

type homeKeyMap struct {
	TopLevelUp    key.Binding
	TopLevelDown  key.Binding
	TopLevelLeft  key.Binding
	TopLevelRight key.Binding
	SidebarNav    key.Binding
	Confirm       key.Binding
}

var defaultHomeKeyMap = homeKeyMap{
	TopLevelUp: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "Move up between sections"),
	),
	TopLevelDown: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "Move down between sections"),
	),
	TopLevelLeft: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "Move left between sections"),
	),
	TopLevelRight: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "Move right between sections"),
	),
	SidebarNav: key.NewBinding(
		key.WithKeys("k", "up", "j", "down"),
		key.WithHelp("k/↑", "Move up within a section"),
		key.WithHelp("j/↓", "Move down within a section"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Confirm an input or focus a component"),
	),
}

var (
	formStyle lipgloss.Style = lipgloss.NewStyle().
			MarginLeft(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6E3F00"))
	listStyle lipgloss.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("#6E3F00")).
			PaddingTop(1)
)

type HomeModel struct {
	sidebarCursor int
	mainCursor    int
	sidebarViews  []tea.Model
	listModel     tea.Model
}

func InitialHome(width int, height int) *HomeModel {
	list := partials.InitialList(width-19, height)
	add := partials.InitialAdd()
	filter := partials.InitialFilter(height - (7))
	sort := components.InitialArrow([]string{"title", "medium", "status", "tags", "release date"}, "Sort", 18, 3)

	sidebarList := []tea.Model{}
	sidebarList = append(sidebarList, &add)
	sidebarList = append(sidebarList, &filter)
	sidebarList = append(sidebarList, &sort)

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
	case DeleteWorkMsg:
		_, cmd = m.listModel.Update(msg)
	case NewWorkMsg:
		_, cmd = m.listModel.Update(msg)
	case tea.WindowSizeMsg:
		_, cmd = m.listModel.Update(msg)
		cmds = tea.Batch(cmds, cmd)
		_, cmd = m.sidebarViews[1].Update(msg)
		cmds = tea.Batch(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultHomeKeyMap.TopLevelUp):
			if m.mainCursor == sidebar && m.sidebarCursor > addBtn {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				nav, ok := cmd().(NavMsg)
				if ok && nav == true {
					m.sidebarCursor--
					_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
					cmds = tea.Batch(cmds, cmd)
				}
			}
		case key.Matches(msg, defaultHomeKeyMap.TopLevelDown):
			if m.mainCursor == sidebar && m.sidebarCursor < numSidebarForms-1 {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				nav, ok := cmd().(NavMsg)
				if ok && nav == true {
					m.sidebarCursor++
					_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
					cmds = tea.Batch(cmds, cmd)
				}
			}
		case key.Matches(msg, defaultHomeKeyMap.TopLevelLeft):
			if m.mainCursor == list {
				m.mainCursor--
				m.listModel, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case key.Matches(msg, defaultHomeKeyMap.TopLevelRight):
			if m.mainCursor == sidebar {
				m.mainCursor++
				m.listModel, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case key.Matches(msg, defaultHomeKeyMap.SidebarNav):
			if m.mainCursor == list {
				m.listModel, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		case key.Matches(msg, defaultHomeKeyMap.Confirm):
			if m.sidebarCursor == addBtn && m.mainCursor == sidebar {
				if len(os.Getenv("DEBUG")) > 0 {
					log.Println("homePage sending AddMsg")
				}
				cmds = tea.Batch(cmds, func() tea.Msg { return (ViewMsg(1)) })
			} else if m.mainCursor == sidebar {
				_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
				cmds = tea.Batch(cmds, cmd)
				if cmd != nil {
					msg, ok := cmd().(partials.FilterMsg)
					if ok {
						_, cmd = m.listModel.Update(msg)
						cmds = tea.Batch(cmds, cmd)
					}
				}
			} else {
				_, cmd = m.listModel.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			}
		default:
			_, cmd = m.sidebarViews[m.sidebarCursor].Update(msg)
			cmds = tea.Batch(cmds, cmd)
			sort, ok := m.sidebarViews[m.sidebarCursor].(*components.ArrowModel)
			if ok {
				DebugLog("Sending SortMsg: ", partials.SortMsg(sort.OptionsCursor))
				_, cmd = m.listModel.Update(partials.SortMsg(sort.OptionsCursor))
				cmds = tea.Batch(cmds, cmd)
			}
		}
	}

	return m, cmds
}

func (m *HomeModel) View() tea.View {
	var c *tea.Cursor
	s := ""
	sidebarContent := []string{}
	for i, form := range m.sidebarViews {
		formView := form.View()
		isFocused := i == m.sidebarCursor && m.mainCursor == sidebar
		sidebarContent = append(sidebarContent, RenderFocused(formStyle, formView.Content, isFocused))
		if formView.Cursor != nil {
			c = formView.Cursor
			c.Y += lipgloss.Height(s)
			c.X += 1
		}
	}
	sidebar := lipgloss.JoinVertical(lipgloss.Center, sidebarContent...)

	isFocused := m.mainCursor == list
	list := RenderFocused(listStyle, m.listModel.View().Content, isFocused)
	s = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, list)

	// Send the UI for rendering
	view := tea.NewView(s)
	view.Cursor = c
	view.AltScreen = true
	return view
}
