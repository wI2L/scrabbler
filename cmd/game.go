package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

type drawPredicate interface {
	Take(t tile, drawCount int) bool
}

type splitTiles struct {
	vowels     tiles
	consonants tiles
}

// game represents a Scrabble game.
type game struct {
	bag       *splitTiles
	draw      *splitTiles
	distrib   distribution
	dict      indexedDict
	drawCount int
	playCount int
	wordLen   uint8
	scrabbles []string
}

// newBag returns a new full splitTiles filled with the
// tiles represented by the given distribution.
func newBag(d distribution) *splitTiles {
	b := &splitTiles{
		vowels:     make(tiles, 0),
		consonants: make(tiles, 0),
	}
	for _, v := range d.letters {
		// Create a tile that represent the letter
		// and add it as many times as its frequency.
		t := tile{
			letter: letter{
				L:         strings.ToUpper(v.L),
				frequency: v.frequency,
				points:    v.points,
			},
		}
		if t.kind() == kindVowel {
			b.vowels.fill(t, v.frequency)
		} else {
			b.consonants.fill(t, v.frequency)
		}
	}
	return b.shuffle()
}

func (s *splitTiles) isEmpty() bool {
	return len(s.vowels) == 0 && len(s.consonants) == 0
}

func (s splitTiles) String() string {
	return s.vowels.String() + " " + s.consonants.String()
}

func (s splitTiles) length() int {
	return len(s.vowels) + len(s.consonants)
}

func (s *splitTiles) tiles() tiles {
	tiles := make(tiles, 0, len(s.vowels)+len(s.consonants))
	tiles = append(tiles, s.vowels...)
	tiles = append(tiles, s.consonants...)

	return tiles
}

func (s *splitTiles) shuffle() *splitTiles {
	s.vowels.shuffle()
	s.consonants.shuffle()

	return s
}

func (s splitTiles) view() string {
	sb := strings.Builder{}

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		s.vowels.view(),
		s.consonants.view(),
	))
	sb.WriteByte('\n')

	return sb.String()
}

func (s *splitTiles) draw(kind letterKind, n uint, predicates []drawPredicate) (tiles, error) {
	var ts *tiles

	switch kind {
	case kindConsonant:
		ts = &s.consonants
	case kindVowel:
		ts = &s.vowels
	default:
		return nil, fmt.Errorf("unknown kind %v", kind)
	}
	ts.shuffle()

	n = min(n, uint(len(*ts)))
	draw := make(tiles, 0, n)

L:
	for i := uint(0); i < n; i++ {
		if len(*ts) == 0 {
			// Collection exhausted, return
			// what has been drawn so far.
			return draw, nil
		}
		idx := rand.Intn(len(*ts))

		for _, p := range predicates {
			if !p.Take((*ts)[idx], int(i)) {
				continue L
			}
		}
		draw.add(ts.pickAt(idx))
	}
	return draw, nil
}

func (g *game) drawWithRequirements(vowels int) error {
	g.resetDraw(false)
	g.drawCount++

	// Pick the required quantity of vowels,
	// minus any leftover frm previous draw.
	v, err := g.bag.draw(kindVowel, uint(vowels-len(g.draw.vowels)), nil)
	if err != nil {
		return err
	}
	g.draw.vowels.add(v...)

	// Draw enough consonants to complete the word.
	wantConsonants := int(g.wordLen) - len(g.draw.vowels) - len(g.draw.consonants)
	c, err := g.bag.draw(kindConsonant, uint(wantConsonants), nil)
	if err != nil {
		return err
	}
	g.draw.consonants.add(c...)

	// Find perfect words from draw.
	g.scrabbles = g.dict.findWords(g.draw.tiles(), g.distrib)

	return nil
}

// playWord withdraws the tiles required to play the given
// word from the slice, or return an error if the word cannot
// be played, leaving the slice untouched.
func (g *game) playWord(word string, check bool) error {
	cpy := append(g.draw.vowels, g.draw.consonants...)

	for _, r := range word {
		if idx := cpy.findTile(string(unicode.ToUpper(r))); idx != -1 {
			// Pop and drop.
			_ = cpy.pickAt(idx)
		} else {
			return fmt.Errorf("word contains unavailable letter '%c'", r)
		}
	}
	if !check {
		for i := range cpy {
			cpy[i].leftover = true
		}
		g.draw.vowels, g.draw.consonants = cpy.splitByKind()
		g.playCount++
		g.drawCount = 0
	}
	return nil
}

// resetDraw puts back all recently drawn tiles
// to the bag. If full is true, all tiles are put
// back to the bag, which imply that any tiles
// from the previous draw are not kept.
func (g *game) resetDraw(full bool) {
	for i := len(g.draw.vowels) - 1; i >= 0; i-- {
		t := g.draw.vowels[i]
		if !t.leftover || full {
			g.bag.vowels.add(g.draw.vowels.pickAt(i))
		}
	}
	for i := len(g.draw.consonants) - 1; i >= 0; i-- {
		t := g.draw.consonants[i]
		if !t.leftover || full {
			g.bag.consonants.add(g.draw.consonants.pickAt(i))
		}
	}
}
