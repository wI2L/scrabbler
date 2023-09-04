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

const (
	defaultDict    = "dictionaries/french/ods8.txt"
	defaultDistrib = "french" // 🇫🇷
)

func init() {
	Root.Flags().BoolP("debug", "v", false, "enable debug logging")

	Root.Flags().StringP("dictionary", "w", defaultDict, "dictionary file path")
	Root.Flags().StringP("distribution", "d", defaultDistrib, "tiles distribution")
}

var (
	Root = &cobra.Command{
		Use:  "scrabbler",
		Long: "scrabbler — pick your tiles, but not yourself",
		RunE: run,
	}
)

func run(cmd *cobra.Command, _ []string) error {
	dn := cmd.Flag("distribution").Value.String()
	dv, ok := distributions[dn]
	if !ok {
		return fmt.Errorf("unknown distribution: %s", dn)
	}
	df := cmd.Flag("dictionary").Value.String()
	dict, err := loadDictionary(df)
	if err != nil {
		return fmt.Errorf("failed to read dictionary: %s", err)
	}
	tw, th, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get term size: %s", err)
	}
	out := termenv.NewOutput(os.Stdout)
	out.SetWindowTitle("scrabbler")

	tui := newTUI(dv, tw, th, dict)
	prg := tea.NewProgram(tui, tea.WithAltScreen(), tea.WithOutput(out))

	envDebug := cmd.Flag("debug").Value.String()

	if ok, err := strconv.ParseBool(envDebug); err == nil && ok {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("failed to open log file: %s", err)
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
