package confirm

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Yes    key.Binding
	No     key.Binding
	Toggle key.Binding
	Quit   key.Binding
	Help   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Yes, k.No, k.Toggle}, // first column
		{k.Help, k.Quit},        // second column
	}
}

var keys = keyMap{
	Yes: key.NewBinding(
		key.WithKeys("y", "Y", "left"),
		key.WithHelp("←/y", "yes"),
	),
	No: key.NewBinding(
		key.WithKeys("n", "N", "right"),
		key.WithHelp("→/n", "no"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/^C", "quit"),
	),
}
