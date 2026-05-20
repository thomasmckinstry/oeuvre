package views

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
)

type AddModel struct {
	form          *components.WorkFormModel
	style         lipgloss.Style
	width, height int
}

func InitialAddModel(width, height int) *AddModel {
	form := components.InitialWorkFormModel(25, height)
	return &AddModel{
		width:  width,
		height: height,
		form:   form,
		style: lipgloss.NewStyle().
			Align(lipgloss.Center).
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.DoubleBorder()),
	}
}

func (m *AddModel) Init() tea.Cmd {
	return nil
}

func (m *AddModel) Update(msg tea.Msg) (*AddModel, tea.Cmd) {
	var cmds tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ViewMsg:
		components.ClearComponents(m.form)
	case tea.WindowSizeMsg:
		m.form.Update(msg)
		m.height = msg.Height
		m.width = msg.Width
	default:
		_, cmd = m.form.Update(msg)
		cmds = tea.Batch(cmds, cmd)
	}
	return m, cmds
}

func (m *AddModel) View() tea.View {
	var c *tea.Cursor
	var s string

	formView := m.form.View()
	c = formView.Cursor
	if c != nil {
		c.Y += 1
		c.X += 2
	}

	s = m.style.Render(formView.Content)
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, s)
	s = lipgloss.PlaceVertical(m.height, lipgloss.Center, s)
	v := tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
