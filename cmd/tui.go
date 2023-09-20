package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
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
	timer    timer.Model
	width    int
	height   int
	insights int
	opts     options
}

type options struct {
	showPoints    bool
	wordLength    int
	minVowels     int
	minConsonants int
	timerDuration time.Duration
}

var _ tea.Model = &tui{}

func newTUI(d distribution, dict indexedDict, width, height int, opts options) *tui {
	return &tui{
		game: &game{
			bag:     newBag(d),
			draw:    &tiles{},
			distrib: d,
			dict:    dict,
			wordLen: opts.wordLength,
		},
		width:  width,
		height: height,
		state:  drawing,
		opts:   opts,
	}
}

func (ui *tui) Init() tea.Cmd {
	if err := ui.game.drawTiles(ui.opts.minVowels, ui.opts.minConsonants); err != nil {
		return tea.Quit
	}
	log.Printf("Initial draw is: %s\n", ui.game.draw)

	ui.input = textinput.New()
	ui.input.Prompt = "Enter tiles played: "
	ui.input.Placeholder = "word"
	ui.input.CharLimit = 0
	ui.input.Validate = func(w string) error {
		return ui.game.playWord(w, true)
	}
	ui.confirm = confirm.New("Accept draw?")

	if ui.opts.timerDuration != 0 {
		ui.timer = timer.NewWithInterval(ui.opts.timerDuration, time.Second)
	}
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
		ui.width, ui.height = m.Width, m.Height
	case timer.TickMsg, timer.StartStopMsg:
		if ui.state == drawing {
			return ui, nil
		}
		var cmd tea.Cmd
		ui.timer, cmd = ui.timer.Update(msg)
		return ui, cmd
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
					if ui.opts.timerDuration != 0 {
						return ui, ui.timer.Start()
					}
				} else {
					ui.insights = 0
					if err := ui.game.drawTiles(ui.opts.minVowels, ui.opts.minConsonants); err != nil {
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
				if err := ui.game.drawTiles(ui.opts.minVowels, ui.opts.minConsonants); err != nil {
					return nil, tea.Quit
				}
				log.Printf("new draw: %s\n", ui.game.draw)

				ui.insights = 0
				ui.state = drawing
				ui.input.Reset()
				if ui.opts.timerDuration != 0 {
					ui.timer.Stop()
					ui.timer.Timeout = ui.opts.timerDuration
				}
				return ui, nil
			}
		case tea.KeyCtrlG:
			ui.insights++
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
	var s string

	if ui.game.bag.isEmpty() && ui.game.draw.isEmpty() {
		s = boldText.Render("Game finished")
	} else {
		s = ui.mainView()
	}
	return lipgloss.Place(
		ui.width, ui.height,
		lipgloss.Center, lipgloss.Center,
		s,
	)
}

func (ui *tui) mainView() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf(boldText.Render("Draw %d.%d"),
		ui.game.playCount,
		ui.game.drawCount),
	)
	sb.WriteString(strings.Repeat("\n", 2))

	// Render the tiles of the draw.
	sb.WriteString(ui.game.draw.view(ui.opts.showPoints))
	sb.WriteByte('\n')

	if ui.game.dict != nil {
		for i, s := range ui.game.scrabbles {
			ui.game.scrabbles[i] = strings.ToLower(s)
		}
		if ui.insights >= 1 {
			if len(ui.game.scrabbles) == 0 {
				sb.WriteString(italicText.Render("no scrabble found"))
			} else {
				var plural string
				if len(ui.game.scrabbles) > 1 {
					plural = "s"
				}
				sb.WriteString(
					fmt.Sprintf("found %d scrabble%s", len(ui.game.scrabbles), plural),
				)
				if ui.insights >= 2 {
					width := ui.width / 3

					sb.WriteByte('\n')
					sb.WriteString(lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(
						scrabbleListView(ui.game.scrabbles, width),
					))
				}
			}
		} else {
			sb.WriteString(faintText.Render("(ctrl+g to show insight)"))
		}
		sb.WriteString(strings.Repeat("\n", 3))
	}
	switch ui.state {
	case drawing:
		sb.WriteString(ui.confirm.View())
	case playing:
		sb.WriteString(ui.input.View())

		if ui.state == playing && ui.opts.timerDuration != 0 {
			sb.WriteString(strings.Repeat("\n", 2))

			ts := formatDuration(ui.timer.Timeout)
			if ui.timer.Timedout() {
				ts = alertText.Render("Time elapsed")
			}
			sb.WriteString(ts)
		}
	}
	return sb.String()
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func scrabbleListView(words []string, maxWidth int) string {
	const wordSep = " ■ "

	var (
		lines     []string
		lineWidth int
		builder   strings.Builder
	)
	for _, w := range words {
		width := 0

		// Compute the rendered width of the word
		// plus separator if needed.
		if lineWidth != 0 {
			width += lipgloss.Width(wordSep)
		}
		width += lipgloss.Width(w)

		// If the length plus the current line width
		// exceed the maximum width, wrap to a new line.
		if maxWidth > 0 && lineWidth+width > maxWidth {
			lines = append(lines, scrabbleList.Render(strings.Clone(builder.String())))
			builder.Reset()
			lineWidth = lipgloss.Width(w)
		} else {
			lineWidth += width
		}
		// After a line wrap, the buffer is empty and
		// a new line shouldn't start with a separator.
		if builder.Len() != 0 {
			builder.WriteString(wordSep)
		}
		builder.WriteString(w)
	}
	// Flush the remaining buffer as the last line.
	if builder.Len() > 0 {
		lines = append(lines, scrabbleList.Render(strings.Clone(builder.String())))
	}
	return lipgloss.JoinVertical(lipgloss.Top, lines...)
}
