package cmd

import "testing"

func Test_newBag(t *testing.T) {
	for _, tt := range []struct {
		lang    string
		distrib distribution
		count   int
	}{
		{
			"french",
			french,
			102,
		},
		{
			"english",
			english,
			100,
		},
	} {
		t.Run(tt.lang, func(t *testing.T) {
			var (
				bag   = newBag(tt.distrib)
				total = len(bag.vowels) + len(bag.consonants)
				freqs = make(map[string]uint)
			)
			if total != tt.count {
				t.Errorf("expected %d tiles, got %d", tt.count, total)
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

func Test_splitTiles_draw_full(t *testing.T) {
	b := newBag(french)

	for !b.isEmpty() {
		vowels, err := b.draw(kindVowel, 3, nil)
		if err != nil {
			t.Fatal(err)
		}
		consonants, err := b.draw(kindConsonant, 4, nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(vowels, consonants)
	}
	if len(b.vowels) != 0 {
		t.Errorf("expected bag to contain no vowels")
	}
	if len(b.consonants) != 0 {
		t.Errorf("expected bag to contain no consonants")
	}
}

func Test_splitTiles_draw_byKind(t *testing.T) {
	const dc = 5
	b := newBag(french)

	for _, tt := range []struct {
		name     string
		kind     letterKind
		tiles    *tiles
		tilesLen int
	}{
		{
			"vowels",
			kindVowel,
			&b.vowels,
			len(b.vowels),
		},
		{
			"consonants",
			kindConsonant,
			&b.consonants,
			len(b.consonants),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			draw, err := b.draw(tt.kind, dc, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(draw) > dc {
				t.Errorf("expected draw to contain %d less or less %ss", dc, tt.kind)
			}
			if len(*tt.tiles) != tt.tilesLen-dc {
				t.Errorf("expected split tiles to contain %d less %ss", dc, tt.kind)
			}
			t.Log(draw)

			for i, tile := range draw {
				if tile.kind() != tt.kind {
					t.Errorf("expected tile %d to be a %s", i, tt.kind)
				}
			}
		})
	}
}

func Test_splitTiles_shuffle(t *testing.T) {
	b := newBag(french)

	// A new splitTiles is automatically shuffled at creation.
	// Snapshot the tiles order and compare it with another
	// shuffle, several times, to ensure proper shuffling.
	var (
		n    = 100
		prev = b.tiles().String()
		next string
	)
	if testing.Short() {
		n = 10
	}
	for i := 0; i < n; i++ {
		b.shuffle()
		next = b.tiles().String()
		if prev == next {
			t.Errorf("expected splitTiles to shuffled")
		}
		prev = next
	}
}
