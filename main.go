package main

import (
	"log"

	"github.com/wI2L/scrabbler/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		log.Fatal(err)
	}
}
