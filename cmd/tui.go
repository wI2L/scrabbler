package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/wI2L/scrabbler/internal/bubbles/confirm"
)

type state int

const (
	drawing state = iota
	playing
)

type tui struct {
	game     *game
	state    state
	input    textinput.Model
	confirm  confirm.Model
	width    int
	height   int
	insights bool
}

var _ tea.Model = &tui{}

func newTUI(d distribution, width, height int, dict indexedDict) *tui {
	return &tui{
		game: &game{
			bag:     newBag(french),
			draw:    &splitTiles{},
			distrib: d,
			dict:    dict,
		},
		width:  width,
		height: height,
		state:  drawing,
	}
}

func (ui *tui) Init() tea.Cmd {
	if err := ui.game.drawTiles(3); err != nil {
		return tea.Quit
	}
	log.Printf("Initial draw is: %s\n", ui.game.draw)

	ui.input = textinput.New()
	ui.input.Prompt = "Enter tiles played: "
	ui.input.Placeholder = "word"
	ui.input.CharLimit = 7
	ui.input.Validate = func(w string) error {
		return ui.game.playWord(w, true)
	}
	ui.confirm = confirm.New("Accept draw?")

	return tea.Batch(
		ui.confirm.Init(),
	)
}

func (ui *tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Width == 0 && m.Height == 0 {
			return ui, nil
		}
		ui.width = m.Width
		ui.height = m.Height
	case tea.KeyMsg:
		switch m.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return ui, tea.Quit
		case tea.KeyEnter:
			switch ui.state {
			case drawing:
				ok := ui.confirm.Value()
				if ok {
					log.Println("draw accepted")

					ui.input.Focus()
					ui.state = playing
				} else {
					ui.insights = false
					if err := ui.game.drawTiles(3); err != nil {
						return nil, tea.Quit
					}
					log.Printf("draw rejected, new draw: %s\n", ui.game.draw)
				}
			case playing:
				word := ui.input.Value()
				if len(word) == 0 {
					break
				}
				if err := ui.game.playWord(word, false); err != nil {
					return nil, tea.Quit
				}
				log.Printf("word played: %s\n", word)
				log.Printf("%d tiles left in the bag, %d remaining tiles from previous draw\n",
					ui.game.bag.length(),
					ui.game.draw.length(),
				)
				// Draw new tiles.
				if err := ui.game.drawTiles(3); err != nil {
					return nil, tea.Quit
				}
				log.Printf("new draw: %s\n", ui.game.draw)

				ui.insights = false
				ui.state = drawing
				ui.input.Reset()
			}
		case tea.KeyCtrlG:
			ui.insights = !ui.insights
		}
	}
	var cmd tea.Cmd
	switch ui.state {
	case drawing:
		ui.confirm, cmd = ui.confirm.Update(msg)
		return ui, cmd
	case playing:
		ui.input, cmd = ui.input.Update(msg)
		return ui, cmd
	}
	return ui, nil
}

func (ui *tui) View() string {
	termWidth, termHeight := ui.width, ui.height
	sb := strings.Builder{}

	if ui.game.bag.isEmpty() && ui.game.draw.isEmpty() {
		sb.WriteString(boldText.Render("Game finished"))
		goto render
	}
	sb.WriteString(fmt.Sprintf(boldText.Render("Draw %d.%d"),
		ui.game.playCount,
		ui.game.drawCount),
	)
	sb.WriteString(strings.Repeat("\n", 2))

	// Render the tiles of the draw.
	sb.WriteString(ui.game.draw.view())
	sb.WriteByte('\n')

	if ui.game.dict != nil {
		scrabbles := ui.game.dict.findWords(ui.game.draw.tiles())

		for i, s := range scrabbles {
			scrabbles[i] = strings.ToLower(s)
		}
		if len(scrabbles) == 0 {
			sb.WriteString(italicText.Render("no scrabble found"))
		} else {
			var plural string
			if len(scrabbles) > 1 {
				plural = "s"
			}
			sb.WriteString(italicText.Render(
				fmt.Sprintf("found %d scrabble%s", len(scrabbles), plural),
			))
			if ui.insights {
				sb.WriteByte('\n')
				sb.WriteString(faintText.Render(
					strings.Join(scrabbles, " • "),
				))
			}
		}
		sb.WriteString(strings.Repeat("\n", 3))
	}
	switch ui.state {
	case drawing:
		sb.WriteString(ui.confirm.View())
	case playing:
		sb.WriteString(ui.input.View())
	}

render:
	return lipgloss.Place(
		termWidth, termHeight,
		lipgloss.Center, lipgloss.Center,
		sb.String(),
	)
}
