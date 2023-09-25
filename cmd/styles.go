package cmd

import "github.com/charmbracelet/lipgloss"

var (
	tileStyle = lipgloss.NewStyle().
			Width(3).
			Bold(true).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("#FBE7D1")).
			BorderForeground(lipgloss.Color("#DFC6A0"))

	inuseTileStyle = tileStyle.Copy().BorderForeground(lipgloss.Color("#FFFFFF"))

	boldText     = lipgloss.NewStyle().Bold(true)
	italicText   = lipgloss.NewStyle().Italic(true)
	faintText    = lipgloss.NewStyle().Faint(true)
	scrabbleList = lipgloss.NewStyle().Faint(true).Italic(true)
	alertText    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)
