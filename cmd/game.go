package cmd

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const maxPredicateRetries = 50

// game represents a Scrabble game.
type game struct {
	bag       *tiles
	draw      *tiles
	distrib   distribution
	dict      indexedDict
	drawCount int
	playCount int
	wordLen   int
	scrabbles []string
}

type drawPredicate interface {
	Take(t tile, drawCount int) bool
}

// newBag returns a new full splitTiles filled with the
// tiles represented by the given distribution.
func newBag(d distribution) *tiles {
	bag := &tiles{
		vowels:     make(rack, 0),
		consonants: make(rack, 0),
	}
	caser := cases.Upper(d.lang)

	for _, v := range d.letters {
		// Create a tile that represent the letter
		// and add it as many times as its frequency.
		t := tile{
			letter: letter{
				L:         caser.String(v.L),
				frequency: v.frequency,
				points:    v.points,
			},
		}
		if t.kind() == kindVowel {
			bag.vowels.fill(t, v.frequency)
		} else {
			bag.consonants.fill(t, v.frequency)
		}
	}
	return bag.shuffle()
}

func (g *game) drawTiles(minVowels, minConsonants int) error {
	g.resetDraw(false)
	g.drawCount++

	defer func() {
		g.scrabbles = g.dict.findWords(g.draw.tiles(), g.distrib)
	}()

	// Pick first the desired quantity of vowels and
	// consonants minus any unplayed tiles from the
	// previous draw, and eventually complete with
	// random tiles.
	if minVowels > 0 {
		v := g.bag.drawByKind(kindVowel, minVowels-len(g.draw.vowels), nil)
		g.draw.vowels.add(v...)
	}
	if minConsonants > 0 {
		c := g.bag.drawByKind(kindConsonant, minConsonants-len(g.draw.consonants), nil)
		g.draw.consonants.add(c...)
	}
	if g.draw.length() == g.wordLen {
		return nil
	}
	r := g.bag.drawRandom(g.wordLen-g.draw.length(), nil)
	v, c := r.splitByKind()

	g.draw.vowels.add(v...)
	g.draw.consonants.add(c...)

	return nil
}

// playWord withdraws the tiles required to play the given
// word from the slice, or return an error if the word cannot
// be played, leaving the slice untouched.
func (g *game) playWord(word string, check bool) error {
	rack := mergeRacks(g.draw.vowels, g.draw.consonants)

	// Normalize word with NFC to combine base
	// characters and modifiers into single runes,
	// and uppercase the result.
	tr := transform.Chain(
		norm.NFC,
		cases.Upper(g.distrib.lang),
	)
	nw, _, _ := transform.String(tr, word)

	for _, r := range nw {
		if idx := rack.findTile(string(r)); idx != -1 {
			_ = rack.pickAt(idx)
		} else {
			return fmt.Errorf("word contains unavailable letter '%c'", r)
		}
	}
	if !check {
		for i := range rack {
			rack[i].leftover = true
		}
		g.draw.vowels, g.draw.consonants = rack.splitByKind()
		g.playCount++
		g.drawCount = 0
	}
	return nil
}

// resetDraw puts back all recently drawn tiles to the bag.
// If full is true, all tiles are put back to the bag, which
// imply that any tiles from the previous draw are not kept.
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

// repetitiveVowelsPredicate is a predicate that
// prevent repetitive vowels during a draw.
type repetitiveVowelsPredicate struct {
	draw      []tile
	threshold int
}

func (p *repetitiveVowelsPredicate) Take(t tile, _ int) bool {
	if t.kind() == kindVowel {
		n := 0
		for _, v := range p.draw {
			if v.L == t.L {
				n++
			}
		}
		if n > p.threshold {
			return false
		}
	}
	p.draw = append(p.draw, t)

	return true
}
