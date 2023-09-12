package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var wordLength uint8

func init() {
	setupFlags()
}

//nolint:gofumpt
var (
	Root = &cobra.Command{
		Use:  "scrabbler",
		Long: "scrabbler — pick your tiles, but not yourself",
		RunE: run,
	}
)

func run(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	dn := cmd.Flag("distribution").Value.String()
	dv, ok := distributions[dn]
	if !ok {
		return fmt.Errorf("unknown distribution: %s", dn)
	}
	var (
		err  error
		dict indexedDict
	)
	dp := cmd.Flag("dictionary").Value.String()
	if dp == "" {
		dict, err = dv.dictionary()
		if err != nil {
			return fmt.Errorf("failed to load dictionary: %s", err)
		}
	} else {
		dict, err = loadDictionaryFile(dp)
		if err != nil {
			return fmt.Errorf("failed to read dictionary file %q: %s", dp, err)
		}
	}
	tw, th, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("cannot get term size: %s", err)
	}
	if wordLength < 7 || wordLength > 8 {
		return fmt.Errorf("word length must be 7 or 8")
	}
	out := termenv.NewOutput(os.Stdout)
	out.SetWindowTitle("scrabbler")

	tui := newTUI(dv, tw, th, dict, wordLength)
	prg := tea.NewProgram(
		tui,
		tea.WithAltScreen(),
		tea.WithOutput(out),
	)
	envDebug := cmd.Flag("debug").Value.String()

	if ok, err := strconv.ParseBool(envDebug); err == nil && ok {
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
	log.Printf("Starting new game with %q distribution...\n", dn)

	_, err = prg.Run()
	return err
}

func setupFlags() {
	f := Root.Flags()

	f.BoolP(
		"debug",
		"v",
		false,
		"enable debug logging to a file",
	)
	f.Uint8VarP(
		&wordLength,
		"word-length",
		"w",
		7,
		"the number of tiles to draw",
	)
	f.StringP(
		"dictionary",
		"d",
		"",
		"custom dictionary file path",
	)
	f.StringP(
		"distribution",
		"l",
		defaultDistrib,
		"tiles distribution language",
	)
}
