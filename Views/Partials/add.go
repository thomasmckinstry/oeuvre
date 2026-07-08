package partials

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/ouevre/utils"
)

type AddModel struct {
	text     string
	selected bool
	width    int
}

type addKeyMap struct {
	Nav key.Binding
}

var defaultAddKey = addKeyMap{
	Nav: key.NewBinding(
		key.WithKeys("ctrl+j", "ctrl+h", "ctrl+k", "ctrl+l"),
		key.WithHelp("ctrl+h/j/k/l", "Navigate away from add"),
	),
}

func InitialAdd() AddModel {
	return AddModel{
		text:     "Add",
		selected: true,
		width:    18,
	}
}

func (m *AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultAddKey.Nav):
			m.selected = !m.selected
		}
	}
	cmd = func() tea.Msg { return utils.NavMsg(true) }
	return m, cmd
}

func (m *AddModel) View() tea.View {
	return tea.NewView(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.text))
}
