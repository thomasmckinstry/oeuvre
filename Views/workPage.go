package views

import (
	key "charm.land/bubbles/v2/key"
	textarea "charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	components "github.com/thomasmckinstry/MediaLogger-TUI/Views/Components"
	. "github.com/thomasmckinstry/MediaLogger-TUI/utils"
)

type WorkPageModel struct {
	work          *components.WorkFormModel
	textArea      textarea.Model
	focused       bool
	writing       bool
	width, height int
	tabCursor     int
	mainCursor    int
	rightCursor   int
	writingCursor int
	tabStyle      lipgloss.Style
	tabsStyle     lipgloss.Style
	buttonStyle   lipgloss.Style
	displayStyle  lipgloss.Style
	detailsStyle  lipgloss.Style
}

type workKeyMap struct {
	TopLevelUp    key.Binding
	TopLevelDown  key.Binding
	TopLevelLeft  key.Binding
	TopLevelRight key.Binding
	Confirm       key.Binding
	Left          key.Binding
	Right         key.Binding
	Exit          key.Binding
}

const (
	work int = iota
	tabs
	add
	mainCursorCount
)

const (
	header int = iota
	display
	rightCursorCount
)

var writeHeader = []string{"NOTE", "REVIEW"}

var (
	unfocused = lipgloss.Color("#6E3F00")
	focused   = lipgloss.Color("#D17600")
)

var defaultWorkMap = workKeyMap{
	TopLevelUp:    key.NewBinding(key.WithKeys("K")),
	TopLevelDown:  key.NewBinding(key.WithKeys("J")),
	TopLevelLeft:  key.NewBinding(key.WithKeys("H")),
	TopLevelRight: key.NewBinding(key.WithKeys("L")),
	Left:          key.NewBinding(key.WithKeys("h", "left")),
	Right:         key.NewBinding(key.WithKeys("l", "right")),
	Confirm:       key.NewBinding(key.WithKeys("enter")),
	Exit:          key.NewBinding(key.WithKeys("esc")),
}

func renderFocused(style lipgloss.Style, content string, isFocused bool) string {
	if isFocused {
		return style.BorderForeground(focused).Render(content)
	} else {
		return style.Render(content)
	}
}

func InitialWorkPage(width, height int) *WorkPageModel {
	workForm := components.InitialWorkFormModel(22, height)

	ti := textarea.New()
	ti.SetVirtualCursor(false)
	ti.SetStyles(textarea.DefaultStyles(true)) // default to dark styles.

	return &WorkPageModel{
		work:     workForm,
		textArea: ti,
		width:    width,
		height:   height,
		tabStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(unfocused),
		tabsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			MarginLeft(1).
			MarginRight(1).
			BorderForeground(unfocused),
		detailsStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			MarginRight(1).
			BorderForeground(unfocused),
		buttonStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(unfocused),
		displayStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			PaddingRight(1).
			Width(width - 28).
			BorderTop(true).
			BorderForeground(unfocused),
	}
}

func (m *WorkPageModel) Init() tea.Cmd {
	return nil
}

func (m *WorkPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd, cmds tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		_, cmd = m.work.Update(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultWorkMap.Confirm):
			if m.mainCursor == work && !m.focused && !m.writing {
				m.focused = true
			} else if m.mainCursor == add && m.rightCursor == header && !m.writing {
				m.writing = true
				return m, cmds
			} else if m.writing && m.writingCursor == 0 && !m.textArea.Focused() {
				m.textArea.Focus()
				return m, cmd
			} else {
				m.work, cmd = m.work.Update(msg)
			}
			cmds = tea.Batch(cmds, cmd)
		case key.Matches(msg, defaultWorkMap.Exit):
			if m.mainCursor == work {
				_, cmd = m.work.Update(msg)
				_, ok := cmd().(ViewMsg)
				if ok && m.focused {
					cmd = nil
					m.focused = false
				}
			} else if m.writing {
				if m.textArea.Focused() {
					m.textArea.Blur()
				} else {
					m.writing = false
					m.textArea.Reset()
				}
			} else {
				cmd = func() tea.Msg { return ViewMsg(0) }
			}
		case key.Matches(msg, defaultWorkMap.TopLevelRight):
			if m.mainCursor < mainCursorCount-1 && !m.focused {
				m.mainCursor++
			} else if m.mainCursor == work {
				_, cmd = m.work.Update(msg)
			}
		case key.Matches(msg, defaultWorkMap.TopLevelLeft):
			if m.rightCursor == display {
				m.mainCursor = 0
			} else if m.mainCursor > work && !m.writing {
				m.mainCursor--
				_, cmd = m.work.Update(msg)
				cmds = tea.Batch(cmds, cmd)
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			}
		case key.Matches(msg, defaultWorkMap.TopLevelDown):
			if m.mainCursor > work && m.rightCursor < display && !m.writing {
				m.rightCursor++
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.writing && !m.textArea.Focused() && m.writingCursor == 0 {
				m.writingCursor++
			}
		case key.Matches(msg, defaultWorkMap.TopLevelUp):
			if m.mainCursor > work && m.rightCursor > header && !m.writing {
				m.rightCursor--
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			} else if m.writing && m.writingCursor > 0 {
				m.writingCursor--
			}
		case key.Matches(msg, defaultWorkMap.Left):
			if m.mainCursor == tabs && m.rightCursor == header && !m.writing {
				m.tabCursor = 0
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			}
		case key.Matches(msg, defaultWorkMap.Right):
			if m.mainCursor == tabs && m.rightCursor == header && !m.writing {
				m.tabCursor = 1
			} else if m.focused {
				_, cmd = m.work.Update(msg)
			}
		default:
			if m.mainCursor == work && m.focused {
				_, cmd = m.work.Update(msg)
			}
		}
	default:
		m.textArea, cmd = m.textArea.Update(msg)
	}

	if m.textArea.Focused() {
		m.textArea, cmd = m.textArea.Update(msg)
	}

	cmds = tea.Batch(cmds, cmd)
	return m, cmds
}

func (m *WorkPageModel) View() tea.View {
	var (
		c *tea.Cursor
		v tea.View
		s string
	)

	if m.writing {
		writingHeader := "NEW " + writeHeader[m.tabCursor]
		textarea := renderFocused(m.tabsStyle, m.textArea.View(), m.writingCursor == 0)
		button := renderFocused(m.buttonStyle, "CONFIRM", m.writingCursor == 1)

		s = lipgloss.JoinVertical(lipgloss.Center, writingHeader, textarea, button)

		v = tea.NewView(s)
		c = m.textArea.Cursor()
		if c != nil {
			c.Y += 2
			c.X += 1
		}
		v.Cursor = c
		v.AltScreen = true
		return v
	}

	workView := m.work.View()
	c = workView.Cursor
	details := workView.Content
	isFocused := m.focused
	details = renderFocused(m.tabsStyle, details, isFocused)

	isSelected := m.mainCursor == work
	details = renderFocused(m.detailsStyle, details, isSelected)

	s = lipgloss.JoinHorizontal(lipgloss.Top, s, details)

	headerContent := "ADD"
	isFocused = m.mainCursor == add && m.rightCursor == header
	headerContent = renderFocused(m.buttonStyle, headerContent, isFocused)

	renderedTabs := []string{}
	tabsArr := []string{"NOTES", "REVIEWS"}
	for i, tab := range tabsArr {
		tab = lipgloss.PlaceHorizontal(9, lipgloss.Center, tab)
		isFocused = m.tabCursor == i
		renderedTabs = append(renderedTabs, renderFocused(m.tabStyle, tab, isFocused))
	}

	tabsContent := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs[0], " ", renderedTabs[1])
	isFocused = m.mainCursor == tabs && m.rightCursor == header
	tabsContent = renderFocused(m.tabsStyle, tabsContent, isFocused)

	headerContent = lipgloss.JoinHorizontal(lipgloss.Top, tabsContent, headerContent)

	displayContent := ""
	isFocused = m.mainCursor > work && m.rightCursor == display
	displayContent = renderFocused(m.displayStyle, displayContent, isFocused)

	rightSide := lipgloss.JoinVertical(lipgloss.Left, headerContent, displayContent)

	s = lipgloss.JoinHorizontal(lipgloss.Top, s, rightSide)

	v = tea.NewView(s)
	v.Cursor = c
	v.AltScreen = true
	return v
}
