package cmd

import "testing"

func Test_tiles_drawByKind(t *testing.T) {
	const dc = 5
	b := newBag(french)

	for _, tt := range []struct {
		name     string
		kind     letterKind
		tiles    *rack
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
			draw := b.drawByKind(tt.kind, dc, nil)

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

func Test_tiles_shuffle(t *testing.T) {
	b := newBag(french)

	// A new tile collection is shuffled at creation.
	// Snapshot the tiles order and compare it with another
	// shuffle several times, to ensure proper shuffling.
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

func Test_tiles_drawByKind_full(t *testing.T) {
	b := newBag(french)

	for !b.isEmpty() {
		v := b.drawByKind(kindVowel, 3, nil)
		c := b.drawByKind(kindConsonant, 4, nil)

		t.Log(v, c)
	}
	if len(b.vowels) != 0 {
		t.Errorf("expected bag to contain no vowels")
	}
	if len(b.consonants) != 0 {
		t.Errorf("expected bag to contain no consonants")
	}
}

func Test_rack_splitByKind(t *testing.T) {
	for _, word := range []string{
		"aggrandizes",
		"archivelog",
		"carburetors",
		"chromosomes",
		"guidance",
		"hemophilia",
		"overhauls",
		"pretrain",
		"prosperous",
		"teenage",
		"toiletries",
		"tracheotomies",
		"weaponizing",
	} {
		tiles := tilesFromWord(word, english)
		vowels, consonants := tiles.splitByKind()

		for _, tile := range vowels {
			if tile.kind() != kindVowel {
				t.Errorf("expected tile to be a vowel")
			}
		}
		for _, tile := range consonants {
			if tile.kind() != kindConsonant {
				t.Errorf("expected tile to be a consonant")
			}
		}
	}
}
