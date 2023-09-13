package cmd

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const dictDir = "dictionaries/"

func Test_loadDictionaryFile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	entries, err := os.ReadDir(dictDir)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string

	for _, e := range entries {
		if e.IsDir() {
			files, err := os.ReadDir(filepath.Join(dictDir, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range files {
				fp := filepath.Join(dictDir, e.Name(), f.Name())
				if filepath.Ext(fp) == ".gz" {
					paths = append(paths, fp)
				}
			}
		}
	}
	for _, path := range paths {
		filename := filepath.Base(path)

		t.Run(filename, func(t *testing.T) {
			dict, err := loadDictionaryFile(path)
			if err != nil {
				t.Fatal(err)
			}
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := f.Close(); err != nil {
					t.Error(err)
				}
			})
			r, err := gzip.NewReader(f)
			if err != nil {
				t.Fatal(err)
			}
			linesCount, err := wc(r)
			if err != nil {
				t.Fatal(err)
			}
			var wordsCount uint
			for _, v := range dict {
				wordsCount += uint(len(v))
			}
			if linesCount != wordsCount {
				t.Errorf("expected %d word in dictionary, got %d", linesCount, wordsCount)
			}
		})
	}
}

//goland:noinspection SpellCheckingInspection
func Test_indexedDict_findWords_french(t *testing.T) {
	dict := frenchDict(t)

	for _, tt := range []struct {
		draw  string
		words []string
	}{
		{
			draw: "OCBWYO",
			words: []string{
				"COWBOY",
			},
		},
		{
			draw: "UOSERSP",
			words: []string{
				"POSEURS",
				"POUSSER",
				"SOUPERS",
			},
		},
		{
			draw:  "XYZABCD",
			words: nil,
		},
	} {
		tiles := tilesFromWord(tt.draw, french)
		words := dict.findWords(tiles, french)

		if want, got := len(tt.words), len(words); want != got {
			t.Errorf("expected %d words, got %d", want, got)
		}

		for i, w := range tt.words {
			if !strings.EqualFold(w, words[i]) {
				t.Errorf("expected word #%d to be %q", i, w)
			}
		}
	}
}

//goland:noinspection SpellCheckingInspection
func Test_indexedDict_findWordsWithBlanks_french(t *testing.T) {
	d := frenchDict(t)

	for _, tt := range []struct {
		letters    []string
		blankCount int // count of blanks
		words      []string
	}{
		{
			[]string{"P", "A", "T", "T", "E", "S"},
			1,
			[]string{
				"OPTATES",
				"PANTETS",
				"PATATES",
				"PATENTS",
				"PATITES",
				"PATTEES",
				"PATTUES",
				"PESTAIT",
				"PESTANT",
				"PETANTS",
				"PETATES",
				"POTATES",
				"PRESTAT",
				"TAPATES",
				"TAPITES",
				"TAPOTES",
				"TIPATES",
				"TOPATES",
				"TYPATES",
			},
		},
		{
			[]string{"P", "T", "T", "E", "S"},
			2,
			[]string{
				"OPTATES",
				"PANTETS",
				"PATATES",
				"PATENTS",
				"PATITES",
				"PATTEES",
				"PATTUES",
				"PEOTTES",
				"PESETTE",
				"PESTAIT",
				"PESTANT",
				"PESTENT",
				"PETANTS",
				"PETATES",
				"PETIOTS",
				"PETITES",
				"PISTENT",
				"PONTETS",
				"POSTENT",
				"POTATES",
				"PRESTAT",
				"PRETEST",
				"PROTETS",
				"PUTIETS",
				"PUTTEES",
				"PUTTERS",
				"SEPTETS",
				"SPITENT",
				"SPITTEE",
				"SPITTER",
				"SPITTES",
				"SPITTEZ",
				"STIPITE",
				"TAPATES",
				"TAPITES",
				"TAPOTES",
				"TIPATES",
				"TOPATES",
				"TOUPETS",
				"TYPATES",
				"TYPOTES",
			},
		},
	} {
		words := d.findWordsWithBlanks(tt.letters, french, tt.blankCount)

		if got, want := len(words), len(tt.words); want != got {
			t.Errorf("got %d words, want %d", got, want)
		}
		if !reflect.DeepEqual(words, tt.words) {
			t.Error("word lists mismatch")
		}
	}
}

func Test_combinationsWithReplacement(t *testing.T) {
	for _, tt := range []struct {
		letters []string
		r       int
		combs   [][]string
	}{
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			1,
			[][]string{
				{"A"},
				{"B"},
				{"C"},
				{"D"},
				{"E"},
				{"F"},
			},
		},
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			2,
			[][]string{
				{"A", "A"},
				{"A", "B"},
				{"A", "C"},
				{"A", "D"},
				{"A", "E"},
				{"A", "F"},
				{"B", "B"},
				{"B", "C"},
				{"B", "D"},
				{"B", "E"},
				{"B", "F"},
				{"C", "C"},
				{"C", "D"},
				{"C", "E"},
				{"C", "F"},
				{"D", "D"},
				{"D", "E"},
				{"D", "F"},
				{"E", "E"},
				{"E", "F"},
				{"F", "F"},
			},
		},
	} {
		combs := combinationsWithReplacement(tt.letters, tt.r)
		if got, want := len(combs), len(tt.combs); want != got {
			t.Errorf("got %d combinations, want %d", got, want)
		}
		if !reflect.DeepEqual(combs, tt.combs) {
			t.Error("combinations mismatch")
		}
	}
}

var cachedDict indexedDict

func frenchDict(t *testing.T) indexedDict {
	t.Helper()

	if cachedDict != nil {
		return cachedDict
	}
	path := filepath.Join(dictDir, "french/ods8.txt.gz")

	dict, err := loadDictionaryFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cachedDict = dict

	return dict
}

func tilesFromWord(word string, d distribution) tiles {
	tiles := make(tiles, 0, len(word))

	for _, r := range word {
		for _, v := range d.letters {
			if v.L == string(r) {
				tiles = append(tiles, tile{letter: v})
			}
		}
	}
	return tiles
}

func wc(r io.Reader) (uint, error) {
	const nl = '\n'

	var c uint
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		var pos int
		for {
			i := bytes.IndexByte(buf[pos:], nl)
			if i == -1 || n == pos {
				break
			}
			pos += i + 1
			c++
		}
		if err == io.EOF {
			break
		}
	}
	return c, nil
}
