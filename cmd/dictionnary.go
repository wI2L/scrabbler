package cmd

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
)

type indexedDict map[string][]string

func loadDictionary(path string) (indexedDict, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	d, err := readDictFile(f)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func readDictFile(r io.ReadCloser) (indexedDict, error) {
	var (
		dict = make(indexedDict)
		scan = bufio.NewScanner(r)
	)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		word := strings.ToUpper(scan.Text())
		r := []rune(word)

		sort.Slice(r, func(i int, j int) bool {
			return r[i] < r[j]
		})
		s := string(r)
		dict[s] = append(dict[s], word)
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return dict, nil
}

func (d indexedDict) findWords(tiles tiles) []string {
	r := make([]rune, 0, len(tiles))

	for _, t := range tiles {
		r = append(r, t.L)
	}
	sort.Slice(r, func(i int, j int) bool {
		return r[i] < r[j]
	})
	return d[string(r)]
}
