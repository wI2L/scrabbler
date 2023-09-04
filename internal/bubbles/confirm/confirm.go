package confirm

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	value  bool
	keys   keyMap
	help   help.Model
	prompt string
}

func New(prompt string) Model {
	m := Model{
		value:  true,
		keys:   keys,
		help:   help.New(),
		prompt: prompt,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Value() bool {
	return m.value
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can
		// gracefully truncate its view as needed.
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Yes):
			m.value = true
		case key.Matches(msg, m.keys.No):
			m.value = false
		case key.Matches(msg, m.keys.Toggle):
			m.value = !m.value
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

var (
	selectionStyle = lipgloss.NewStyle().Bold(true)
)

const arrowRight = "▸"

func (m Model) View() string {
	sb := strings.Builder{}

	sb.WriteString(m.prompt)
	sb.WriteString(" ") //spacer

	if m.value {
		sb.WriteString(arrowRight)
		sb.WriteString(selectionStyle.Copy().Foreground(lipgloss.Color("#76E083")).Render("Yes"))
		sb.WriteString(" ")
		sb.WriteString("No")
	} else {
		sb.WriteString("Yes")
		sb.WriteString(" ")
		sb.WriteString(arrowRight)
		sb.WriteString(selectionStyle.Copy().Foreground(lipgloss.Color("#f9746a")).Render("No"))
	}
	var hv string

	if !m.help.ShowAll {
		hv = m.help.ShortHelpView(m.keys.ShortHelp())
	} else {
		hs := m.help.Styles.FullSeparator.Render(m.help.FullSeparator)
		hw := lipgloss.Width(hs)
		st := lipgloss.NewStyle().MarginLeft(hw)
		hv = st.Render(m.help.FullHelpView(m.keys.FullHelp()))
	}
	sb.WriteString(strings.Repeat("\n", 2))
	sb.WriteString(hv)

	return sb.String()
}
