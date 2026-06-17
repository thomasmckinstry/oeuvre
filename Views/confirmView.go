package views

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	. "github.com/thomasmckinstry/ouevre/utils"
)

const (
	cancel int = iota
	confirm
)

type ConfirmModel struct {
	cursor       int
	width        int
	height       int
	confirmation ConfirmationMsg
}

type confirmKeyMap struct {
	Left    key.Binding
	Right   key.Binding
	Confirm key.Binding
}

var defaultConfirmMap confirmKeyMap = confirmKeyMap{
	Left:    key.NewBinding(key.WithKeys("H", "h", "left")),
	Right:   key.NewBinding(key.WithKeys("L", "l", "right")),
	Confirm: key.NewBinding(key.WithKeys("enter")),
}

func InitialConfirmModel(width int, height int) *ConfirmModel {
	return &ConfirmModel{
		width:  width,
		height: height,
	}
}

func (m *ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m *ConfirmModel) Update(msg tea.Msg) (*ConfirmModel, tea.Cmd) {
	var cmd, cmds tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case ConfirmationMsg:
		m.confirmation = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultConfirmMap.Confirm):
			DebugLog("Confirming choice on confirmView: ", m.cursor)
			if m.cursor == confirm {
				cmd = m.confirmation.Function()
				cmds = tea.Batch(cmds, cmd)
				cmds = tea.Batch(cmds, func() tea.Msg { return ViewMsg(HomePage) })
			} else {
				cmd = func() tea.Msg { return ViewMsg(HomePage) }
			}
		case key.Matches(msg, defaultConfirmMap.Left):
			m.cursor = cancel
		case key.Matches(msg, defaultConfirmMap.Right):
			m.cursor = confirm
		}
	}
	return m, tea.Batch(cmds, cmd)
}

func (m *ConfirmModel) View() tea.View {
	var (
		isFocused bool
		v         tea.View
	)
	isFocused = m.cursor == cancel
	cancelBtn := RenderFocused(ButtonStyle, "CANCEL", isFocused)
	isFocused = m.cursor == confirm
	confirmBtn := RenderFocused(ButtonStyle, "CONFIRM", isFocused)
	s := lipgloss.JoinHorizontal(lipgloss.Top, cancelBtn, confirmBtn)
	s = lipgloss.JoinVertical(lipgloss.Center, m.confirmation.Msg, s)
	s = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
	v = tea.NewView(s)
	v.AltScreen = true
	return v
}
