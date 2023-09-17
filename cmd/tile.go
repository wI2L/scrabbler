package cmd

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
	tr := transform.Chain(
		norm.NFD,                           // decompose
		runes.Remove(runes.In(unicode.Mn)), // remove diacritics
		norm.NFC,                           // recompose
		cases.Upper(language.Und),          // uppercase
	)
	s, _, _ := transform.String(tr, t.L)

	switch s {
	case "A", "E", "I", "O", "U", "Y":
		return kindVowel
	default:
		return kindConsonant
	}
}

func (t tile) view(withPoints bool) string {
	if withPoints {
		return t.L + subscriptPoints(t.points)
	}
	return t.L
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

func (s tiles) view(withPoints bool) string {
	strs := make([]string, 0, len(s))

	for _, t := range s {
		var style lipgloss.Style
		if t.leftover {
			style = leftoverTileStyle
		} else {
			style = tileStyle
		}
		if withPoints {
			style = style.Align(lipgloss.Right)
		}
		strs = append(strs, style.Render(t.view(withPoints)))
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

func subscriptPoints(i uint) string {
	const zero = 0x00002080 // U+2080
	switch {
	case i == 0:
		return string(rune(zero))
	case i == 10:
		return "ₓ" // X for 10
	default:
		return string(rune(zero + i))
	}
}
