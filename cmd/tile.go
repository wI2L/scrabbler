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

type rack []tile

type tiles struct {
	vowels     rack
	consonants rack
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

func (r rack) String() string {
	s := make([]string, 0, len(r))
	for _, t := range r {
		s = append(s, t.L)
	}
	return strings.Join(s, " ")
}

func (t tile) view(withPoints bool) string {
	if withPoints {
		return t.L + subscriptPoints(t.points)
	}
	return t.L
}

// add appends the given tiles to the slice.
func (r *rack) add(tiles ...tile) {
	for _, t := range tiles {
		*r = append(*r, t)
	}
}

// fill adds blankCount times the same tile to the slice.
func (r *rack) fill(t tile, n uint) {
	for i := uint(0); i < n; i++ {
		*r = append(*r, t)
	}
}

// shuffle randomizes the order of the tiles.
func (r *rack) shuffle() {
	rand.Shuffle(len(*r), func(i, j int) {
		(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
	})
}

// pickAt pops and return the tile at the given position.
// The index must not exceed the length of the slice.
func (r *rack) pickAt(idx int) tile {
	t := (*r)[idx]
	*r = append((*r)[:idx], (*r)[idx+1:]...)
	return t
}

// findTiles returns the index of the first tile
// with the given letter,
func (r rack) findTile(letter string) int {
	for i, t := range r {
		if t.L == letter {
			return i
		}
	}
	return -1
}

// splitByKind splits the tiles of the slice by their letter's kind.
func (r rack) splitByKind() (vowels, consonants rack) {
	for _, t := range r {
		switch t.kind() {
		case kindVowel:
			vowels = append(vowels, t)
		case kindConsonant:
			consonants = append(consonants, t)
		}
	}
	return
}

func (r rack) view(withPoints bool) string {
	strs := make([]string, 0, len(r))

	for _, t := range r {
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

func mergeRacks(r1, r2 rack) rack {
	r := make(rack, 0, len(r1)+len(r2))
	r = append(r, r1...)
	r = append(r, r2...)

	return r
}

func (s tiles) String() string {
	return s.vowels.String() + " " + s.consonants.String()
}

func (s tiles) length() int {
	return len(s.vowels) + len(s.consonants)
}

func (s *tiles) isEmpty() bool {
	return len(s.vowels) == 0 && len(s.consonants) == 0
}

func (s *tiles) tiles() rack {
	r := make(rack, 0, len(s.vowels)+len(s.consonants))
	r = append(r, s.vowels...)
	r = append(r, s.consonants...)

	return r
}

func (s *tiles) shuffle() *tiles {
	s.vowels.shuffle()
	s.consonants.shuffle()

	return s
}

func (s tiles) view(withPoints bool) string {
	sb := strings.Builder{}

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		s.vowels.view(withPoints),
		s.consonants.view(withPoints),
	))
	sb.WriteByte('\n')

	return sb.String()
}

func (s *tiles) drawRandom(n int, predicates []drawPredicate) rack {
	n = min(n, s.length())
	if n == 0 {
		return nil
	}
	draw := make(rack, 0, n)

L:
	for i, j := 0, 0; i < n; j++ {
		if s.length() == 0 {
			return draw
		}
		s.vowels.shuffle()
		s.consonants.shuffle()

		// Pick a random tile spanning both slices.
		idx := rand.Intn(s.length())

		var ts *rack
		if idx < len(s.vowels) {
			ts = &s.vowels
		} else {
			ts = &s.consonants
			// Offset index for consonants slices.
			idx -= len(s.vowels)
		}
		if j < maxPredicateRetries {
			for _, p := range predicates {
				if !p.Take((*ts)[idx], i) {
					continue L
				}
			}
		}
		i++
		j = 0
		draw.add(ts.pickAt(idx))
	}
	return draw
}

func (s *tiles) drawByKind(kind letterKind, n int, predicates []drawPredicate) rack {
	var ts *rack

	switch kind {
	case kindConsonant:
		ts = &s.consonants
	case kindVowel:
		ts = &s.vowels
	}
	n = min(n, len(*ts))
	if n == 0 {
		return nil
	}
	draw := make(rack, 0, n)

L:
	for i, j := 0, 0; i < n; j++ {
		if len(*ts) == 0 {
			return draw
		}
		ts.shuffle()

		idx := rand.Intn(len(*ts))

		if j < maxPredicateRetries {
			for _, p := range predicates {
				if !p.Take((*ts)[idx], i) {
					continue L
				}
			}
		}
		i++
		j = 0
		draw.add(ts.pickAt(idx))
	}
	return draw
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
