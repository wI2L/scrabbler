package cmd

import (
	"sort"
	"testing"
	"unicode/utf8"

	"golang.org/x/exp/maps"
)

func Test_distribution_tiles(t *testing.T) {
	keys := maps.Keys(distributions)
	sort.Strings(keys)

	for _, k := range keys {
		d := distributions[k]
		if d.tileCount == 0 {
			t.Errorf("expected distribution %q to have a non-zero tile count", k)
		}
		b := newBag(d)
		if bl := b.length(); bl != d.tileCount {
			t.Errorf("%s: expected bag to contain %d tiles, got %d", k, d.tileCount, bl)
		}
	}
}

func Test_distribution_letters(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// This test ensure that the letters of a distribution are
	// only represented by a single Unicode code point instead
	// of a base character with combining diacritical marks/modifiers.
	for k, d := range distributions {
		for _, l := range d.letters {
			if n := utf8.RuneCountInString(l.L); n > 1 {
				t.Errorf("%s: letter %q has a rune length of %d", k, l.L, n)
			}
		}
	}
}
