package cmd

import (
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type letterKind uint

const (
	kindVowel letterKind = iota
	kindConsonant
)

// tile represents a Scrabble tile.
type tile struct {
	letter
	leftover bool
}

// kind returns the kind of the tile's letter.
func (t tile) kind() letterKind {
	switch strings.ToUpper(t.L) {
	case "A", "E", "I", "O", "U", "Y":
		return kindVowel
	default:
		return kindConsonant
	}
}

type tiles []tile

// add appends the given tiles to the slice.
func (s *tiles) add(tiles ...tile) {
	for _, t := range tiles {
		*s = append(*s, t)
	}
}

// fill adds blankCount times the same tile to the slice.
func (s *tiles) fill(t tile, n uint) {
	for i := uint(0); i < n; i++ {
		*s = append(*s, t)
	}
}

// shuffle randomizes the order of the tiles.
func (s *tiles) shuffle() {
	rand.Shuffle(len(*s), func(i, j int) {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	})
}

// pickAt pops and return the tile at the given position.
// The index must not exceed the length of the slice.
func (s *tiles) pickAt(idx int) tile {
	t := (*s)[idx]
	*s = append((*s)[:idx], (*s)[idx+1:]...)
	return t
}

// findTiles returns the index of the first tile
// with the given letter,
func (s tiles) findTile(letter string) int {
	for i, t := range s {
		if t.L == letter {
			return i
		}
	}
	return -1
}

// splitByKind splits the tiles of the slice by their letter's kind.
func (s tiles) splitByKind() (vowels, consonants tiles) {
	for _, t := range s {
		switch t.kind() {
		case kindVowel:
			vowels = append(vowels, t)
		case kindConsonant:
			consonants = append(consonants, t)
		}
	}
	return
}

func (s tiles) String() string {
	r := make([]string, 0, len(s))
	for _, t := range s {
		r = append(r, t.L)
	}
	return strings.Join(r, " ")
}

func (s tiles) view() string {
	var strs []string
	for _, t := range s {
		if t.leftover {
			strs = append(strs, leftoverTileStyle.Render(t.L))
		} else {
			strs = append(strs, tileStyle.Render(t.L))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

func (k letterKind) String() string {
	switch k {
	case kindVowel:
		return "vowel"
	case kindConsonant:
		return "consonant"
	default:
		return "<unknown>"
	}
}
