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
	prompt string
	keys   keyMap
	help   help.Model
	styles Styles
}

type Styles struct {
	Accept lipgloss.Style
	Reject lipgloss.Style
}

func New(prompt string, defaultValue bool) Model {
	h := help.New()
	h.FullSeparator = strings.Repeat(" ", 3)

	baseStyle := lipgloss.NewStyle().Bold(true)

	return Model{
		value:  defaultValue,
		prompt: prompt,
		keys:   keys,
		help:   h,
		styles: Styles{
			Accept: baseStyle.Copy().Foreground(lipgloss.Color("#76E083")),
			Reject: baseStyle.Copy().Foreground(lipgloss.Color("#f9746a")),
		},
	}
}

func (m Model) Value() bool {
	return m.value
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Set a width on the help menu so that it can
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

func (m Model) View() string {
	sb := strings.Builder{}

	sb.WriteString(m.prompt)
	sb.WriteString(" ") // spacer

	const arrow = "▸"

	if m.value {
		sb.WriteString(arrow)
		sb.WriteString(m.styles.Accept.Render("Yes"))
		sb.WriteString(" ")
		sb.WriteString("No")
	} else {
		sb.WriteString("Yes")
		sb.WriteString(" ")
		sb.WriteString(arrow)
		sb.WriteString(m.styles.Reject.Render("No"))
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
	sb.WriteString(strings.Repeat("\n", 3))
	sb.WriteString(hv)

	return sb.String()
}
