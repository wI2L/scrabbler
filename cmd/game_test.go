package cmd

import "testing"

func Test_newBag(t *testing.T) {
	for _, tt := range []struct {
		lang    string
		distrib distribution
	}{
		{
			"french",
			french,
		},
		{
			"english",
			english,
		},
	} {
		t.Run(tt.lang, func(t *testing.T) {
			var (
				bag   = newBag(tt.distrib)
				total = len(bag.vowels) + len(bag.consonants)
				freqs = make(map[string]uint)
			)
			if total != tt.distrib.tileCount {
				t.Errorf("expected %d tiles, got %d", tt.distrib.tileCount, total)
			}
			for _, t := range bag.vowels {
				freqs[t.L]++
			}
			for _, t := range bag.consonants {
				freqs[t.L]++
			}
			for _, v := range tt.distrib.letters {
				f := freqs[v.L]
				if f != v.frequency {
					t.Errorf("expected frequency of %d for letter %q, got %d", v.frequency, v.L, f)
				}
			}
		})
	}
}
