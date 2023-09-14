package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/wI2L/scrabbler/cmd"
)

func main() {
	log.SetFlags(0)

	cmd.Root.SilenceErrors = true

	if err := cmd.Root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
