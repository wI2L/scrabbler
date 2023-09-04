package cmd

import (
	"testing"
)

func Test_tiles_splitByKind(t *testing.T) {
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
		tiles := tilesFromWord(word, french)
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
