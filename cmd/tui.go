package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/wI2L/scrabbler/internal/bubbles/confirm"
	"github.com/wI2L/scrabbler/internal/bubbles/gridmenu"
)

type state int

const (
	lang state = iota
	draw
	play
)

type tui struct {
	game     *game
	state    state
	input    textinput.Model
	confirm  confirm.Model
	menu     gridmenu.Model
	timer    timer.Model
	width    int
	height   int
	insights int
	opts     options
}

type options struct {
	dictPath      string
	showPoints    bool
	wordLength    int
	minVowels     int
	minConsonants int
	timerDuration time.Duration
	predicates    []drawPredicate
}

var _ tea.Model = &tui{}

func newTUI(distribName string, width, height int, opts options) (*tui, error) {
	tui := tui{
		state:  lang,
		width:  width,
		height: height,
		opts:   opts,
	}
	if distribName != "" {
		if err := tui.initGame(distribName); err != nil {
			return nil, err
		}
	}
	return &tui, nil
}

func (ui *tui) initGame(dn string) error {
	distrib, ok := distributions[dn]
	if !ok {
		return fmt.Errorf("unknown distribution: %s", dn)
	}
	var (
		err  error
		dict indexedDict
	)
	if ui.opts.dictPath == "" {
		dict, err = distrib.dictionary(ui.opts.wordLength)
		if err != nil {
			return fmt.Errorf("failed to load dictionary: %s", err)
		}
	} else {
		dict, err = loadDictionaryFile(ui.opts.dictPath, distrib.lang, ui.opts.wordLength)
		if err != nil {
			return fmt.Errorf("failed to read dictionary file %q: %s", ui.opts.dictPath, err)
		}
	}
	ui.game = &game{
		bag:     newBag(distrib),
		draw:    &tiles{},
		distrib: distrib,
		dict:    dict,
		wordLen: ui.opts.wordLength,
	}
	ui.game.drawTiles(
		ui.opts.minVowels,
		ui.opts.minConsonants,
		ui.opts.predicates...,
	)
	log.Printf("Starting new game with %q distribution...\n", dn)
	log.Printf("Initial draw is: %s\n", ui.game.draw)

	ui.state = draw

	return nil
}

func (ui *tui) Init() tea.Cmd {
	ui.input = textinput.New()
	{
		ui.input.CharLimit = 0
		ui.input.Prompt = "Enter tiles played: "
		ui.input.Placeholder = "word"
		ui.input.Validate = func(w string) error {
			return ui.game.playWord(w, true)
		}
	}
	ui.menu = gridmenu.New(distribChoices(), 4, 7)
	{
		ui.menu.Width = ui.width
		ui.menu.Margin(6, 2)
	}
	ui.confirm = confirm.New("Accept draw?", true)

	if ui.opts.timerDuration != 0 {
		ui.timer = timer.NewWithInterval(ui.opts.timerDuration, time.Second)
	}
	return nil
}

func (ui *tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Width == 0 && m.Height == 0 {
			return ui, nil
		}
		ui.width, ui.height = m.Width, m.Height

	case timer.TickMsg, timer.StartStopMsg:
		if ui.state != play {
			return ui, nil
		}
		var cmd tea.Cmd
		ui.timer, cmd = ui.timer.Update(msg)
		return ui, cmd

	case tea.KeyMsg:
		switch m.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return ui, tea.Quit
		case tea.KeyCtrlR:
			if ui.state == draw {
				ui.insights = 0
				ui.game.resetDraw(true)
				ui.game.drawTiles(ui.opts.minVowels, ui.opts.minConsonants, ui.opts.predicates...)
			}
			return ui, nil
		case tea.KeyEnter:
			switch ui.state {
			case lang:
				if err := ui.initGame(ui.menu.Selection()); err != nil {
					return nil, tea.Quit
				}
				ui.state = draw
				return ui, nil
			case draw:
				ok := ui.confirm.Value()
				if ok {
					log.Println("draw accepted")

					ui.input.Focus()
					ui.state = play
					if ui.opts.timerDuration != 0 {
						return ui, ui.timer.Start()
					}
				} else {
					ui.insights = 0
					ui.game.drawTiles(ui.opts.minVowels, ui.opts.minConsonants, ui.opts.predicates...)

					log.Printf("draw rejected, new draw: %s\n", ui.game.draw)
				}
				return ui, nil
			case play:
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
				ui.game.drawTiles(ui.opts.minVowels, ui.opts.minConsonants, ui.opts.predicates...)

				log.Printf("new draw: %s\n", ui.game.draw)

				ui.insights = 0
				ui.state = draw
				ui.input.Reset()

				if ui.opts.timerDuration != 0 {
					ui.timer.Stop()
					ui.timer.Timeout = ui.opts.timerDuration
				}
				return ui, nil
			}
		case tea.KeyCtrlG:
			ui.insights++
			return ui, nil
		}
	}
	var cmd tea.Cmd

	switch ui.state {
	case lang:
		ui.menu, cmd = ui.menu.Update(msg)
		return ui, cmd
	case draw:
		ui.confirm, cmd = ui.confirm.Update(msg)
		return ui, cmd
	case play:
		ui.input, cmd = ui.input.Update(msg)
		return ui, cmd
	}
	return ui, nil
}

func (ui tui) View() string {
	var s string

	if ui.state == lang {
		s += lipgloss.NewStyle().Bold(true).Render("Choose a language")
		s += strings.Repeat("\n", 3)
		s += ui.menu.View()
	} else {
		if ui.game.bag.isEmpty() && ui.game.draw.isEmpty() {
			s = boldText.Render("Game finished")
		} else {
			s = ui.runningView()
		}
	}
	return lipgloss.Place(
		ui.width, ui.height,
		lipgloss.Center, lipgloss.Center,
		s,
	)
}

func (ui tui) runningView() string {
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
						ui.wordListView(width),
					))
				}
			}
		} else {
			sb.WriteString(faintText.Render("(ctrl+g to show insight)"))
		}
		sb.WriteString(strings.Repeat("\n", 3))
	}
	switch ui.state {
	case draw:
		sb.WriteString(ui.confirm.View())
	case play:
		sb.WriteString(ui.input.View())

		if ui.state == play && ui.opts.timerDuration != 0 {
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

func (ui tui) wordListView(maxWidth int) string {
	const wordSep = " ■ "

	var (
		lines     []string
		lineWidth int
		builder   strings.Builder
	)
	for _, w := range ui.game.scrabbles {
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

func distribChoices() []gridmenu.Choice {
	c := make([]gridmenu.Choice, 0, len(distributions))

	for k, v := range distributions {
		c = append(c, gridmenu.Choice{
			Name:        k,
			Description: v.name,
		})
	}
	sort.Slice(c, func(i, j int) bool {
		return c[i].Name < c[j].Name
	})
	return c
}
