package cmd

import (
	"sort"
	"testing"

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
