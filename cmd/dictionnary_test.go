package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func Test_readDictFile(t *testing.T) {
	const base = "../dictionaries/"

	entries, err := os.ReadDir(base)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string

	for _, e := range entries {
		if e.IsDir() {
			files, err := os.ReadDir(filepath.Join(base, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range files {
				fp := filepath.Join(base, e.Name(), f.Name())
				fmt.Println(fp)
				if filepath.Ext(fp) == ".txt" {
					paths = append(paths, fp)
				}
			}
		}
	}
	for _, path := range paths {
		filename := filepath.Base(path)

		t.Run(filename, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := f.Close(); err != nil {
					t.Error(err)
				}
			})
			dict, err := readDictFile(f)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				t.Fatal(err)
			}
			linesCount, err := wc(f)
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

func Test_indexedDict_findWords_french(t *testing.T) {
	const path = "../dictionaries/french/ods8.txt"

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	dict, err := readDictFile(f)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		draw  string
		words []string
	}{
		{
			draw: "ocbwyo",
			words: []string{
				"cowboy",
			},
		},
		{
			draw: "uosersp",
			words: []string{
				"pousser",
				"poseurs",
				"soupers",
			},
		},
		{
			draw:  "xyzabcd",
			words: nil,
		},
	} {
		tiles := tilesFromWord(tt.draw, french)
		words := dict.findWords(tiles)

		if want, got := len(tt.words), len(words); want != got {
			t.Errorf("expected %d words, got %d", want, got)
		}
		sort.Strings(words)
		sort.Strings(tt.words)

		for i, w := range tt.words {
			if !strings.EqualFold(w, words[i]) {
				t.Errorf("expected word #%d to be %q", i, w)
			}
		}
	}
}

func tilesFromWord(word string, d distribution) tiles {
	tiles := make(tiles, 0, len(word))

	for _, r := range []rune(strings.ToUpper(word)) {
		if v, ok := d[r]; ok {
			tiles = append(tiles, tile{L: r, letterProps: v})
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
