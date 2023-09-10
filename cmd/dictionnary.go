package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
	"unicode"
)

var (
	//go:embed dictionaries/french/ods8.txt.gz
	ods8 []byte

	//go:embed dictionaries/english/sowpods.txt.gz
	sowpods []byte
)

type indexedDict map[string][]string

func (id indexedDict) findWords(tiles tiles, d distribution) []string {
	r := make([]string, 0, len(tiles))

	blanks := 0
	for _, t := range tiles {
		if t.L == blank {
			blanks++
		} else {
			r = append(r, strings.ToUpper(t.L))
		}
	}
	if blanks > 0 {
		return id.findWordsWithBlanks(r, d, blanks)
	}
	slices.Sort(r)

	return id[strings.Join(r, "")]
}

func (id indexedDict) findWordsWithBlanks(r []string, d distribution, n int) []string {
	var words []string

	s := make([]string, 0, len(r)+n)
	for _, c := range combinationsWithReplacement(d.alphabet(), n) {
		s = s[:0]
		s = append(s, r...)
		s = append(s, c...)
		slices.Sort(s)

		if w, ok := id[strings.Join(s, "")]; ok {
			words = append(words, w...)
		}
	}
	sort.Strings(words)
	return words
}

func loadDictionaryFile(path string) (indexedDict, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	d, err := parseDictionary(f)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func checkWord(w string) error {
	for _, r := range w {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("'%c' is not a letter", r)
		}
	}
	return nil
}

func parseDictionary(r io.ReadCloser) (indexedDict, error) {
	var (
		dict = make(indexedDict)
		scan = bufio.NewScanner(r)
	)
	scan.Split(bufio.ScanLines)

	for i := 1; scan.Scan(); i++ {
		line := scan.Text()
		line = strings.TrimSpace(line)
		if err := checkWord(line); err != nil {
			return nil, fmt.Errorf("found invalid word %q at line %d: %s", line, i, err)
		}
		word := strings.ToUpper(line)
		r := []rune(word)
		slices.Sort(r)
		s := string(r)
		dict[s] = append(dict[s], word)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return dict, nil
}

// Port of Python3 eponymous function from itertools package.
// https://docs.python.org/3/library/itertools.html#itertools.combinations_with_replacement
func combinationsWithReplacement(s []string, r int) [][]string {
	n := len(s)
	if n == 0 || r == 0 {
		return nil
	}
	indices := make([]int, r)
	var combs [][]string
	for {
		c := make([]string, r)
		for i, idx := range indices {
			c[i] = s[idx]
		}
		combs = append(combs, c)

		// Find the rightmost index that
		// can be incremented.
		i := r - 1
		for ; i >= 0; i-- {
			if indices[i] != n-1 {
				break
			}
		}
		// If no index can be incremented,
		// we're done.
		if i < 0 {
			break
		}
		// Increment the index and set all
		// following indices to the same value.
		indices[i]++

		for j := i + 1; j < r; j++ {
			indices[j] = indices[i]
		}
	}
	return combs
}
