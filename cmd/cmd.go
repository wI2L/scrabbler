package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	debug         bool
	wordLength    uint8
	vowels        uint8
	consonants    uint8
	showPoints    bool
	timerDuration time.Duration

	Root = &cobra.Command{
		Use:  "scrabbler",
		Long: "scrabbler — pick tiles, but not yourself!",
		RunE: run,
	}
)

func init() {
	setupFlags()
}

func run(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	dn := cmd.Flag("distribution").Value.String()
	dp := cmd.Flag("dictionary").Value.String()

	tw, th, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("cannot get term size: %s", err)
	}
	if wordLength < 7 || wordLength > 8 {
		return fmt.Errorf("word length must be 7 or 8")
	}
	if vowels+consonants > wordLength {
		return fmt.Errorf("required vowels and consonants exceed word length")
	}
	out := termenv.NewOutput(os.Stdout)
	out.SetWindowTitle("scrabbler")

	tui, err := newTUI(dn, tw, th, options{
		dictPath:      dp,
		wordLength:    int(wordLength),
		minVowels:     int(vowels),
		minConsonants: int(consonants),
		showPoints:    showPoints,
		timerDuration: timerDuration,
	})
	if err != nil {
		return err
	}
	prg := tea.NewProgram(
		tui,
		tea.WithAltScreen(),
		tea.WithOutput(out),
	)
	if debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("cannot open log file: %s", err)
		}
		defer func() {
			_ = f.Close()
		}()
	} else {
		log.SetOutput(io.Discard)
	}
	// See https://github.com/charmbracelet/lipgloss/issues/73
	lipgloss.SetHasDarkBackground(termenv.HasDarkBackground())

	_, err = prg.Run()

	return err
}

func setupFlags() {
	f := Root.Flags()

	f.StringP(
		"dictionary",
		"d",
		"",
		"custom dictionary file path",
	)
	f.StringP(
		"distribution",
		"l",
		"",
		"letter distribution language",
	)
	f.Uint8VarP(
		&wordLength,
		"word-length",
		"w",
		7,
		"the number of tiles to draw",
	)
	f.Uint8Var(
		&vowels,
		"vowels",
		0,
		"number of required vowel letters",
	)
	f.Uint8Var(
		&consonants,
		"consonants",
		0,
		"number of required consonant letters",
	)
	f.BoolVarP(
		&debug,
		"debug",
		"v",
		false,
		"enable debug logging",
	)
	f.BoolVarP(
		&showPoints,
		"show-points",
		"p",
		false,
		"show letter points",
	)
	f.DurationVarP(
		&timerDuration,
		"timer",
		"t",
		0,
		"enable play timer (default 5m)",
	)
	// Set a default value for the timer duration
	// if the flag is passed without a value.
	f.Lookup("timer").NoOptDefVal = "5m"
}
