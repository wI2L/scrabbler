package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

const maxPredicateRetries = 50

type drawPredicate interface {
	Reset(draw rack)
	Take(t tile, drawCount int) bool
}

type predicateList struct {
	value   []drawPredicate
	changed bool
}

func (pl predicateList) String() string { return "" }
func (pl predicateList) Type() string   { return "key=[val],..." }

func (pl *predicateList) Set(val string) error {
	ss := strings.Split(val, ",")
	ps := make([]drawPredicate, 0, len(ss))

	for _, pair := range ss {
		p := strings.TrimSpace(pair)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) < 1 {
			return fmt.Errorf("invalid predicate name")
		}
		name := kv[0]

		switch name {
		case "dup-vowels":
			if len(kv) != 2 {
				return fmt.Errorf("predicate '%s' requires a value", name)
			}
			threshold, err := strconv.Atoi(kv[1])
			if err != nil {
				return err
			}
			ps = append(ps, &duplicateVowelsPredicate{threshold: threshold})
		default:
			return fmt.Errorf("unknown predicate: %s", name)
		}
	}
	if !pl.changed {
		pl.value = ps
	} else {
		pl.value = append(pl.value, ps...)
	}
	pl.changed = true

	return nil
}

// duplicateVowelsPredicate is a predicate that
// prevent repetitive vowels pick during a draw.
type duplicateVowelsPredicate struct {
	draw      []tile
	threshold int
}

func (p *duplicateVowelsPredicate) Reset(draw rack) {
	p.draw = p.draw[:0]
	for _, t := range draw {
		p.draw = append(p.draw, t)
	}
}

func (p *duplicateVowelsPredicate) Take(t tile, _ int) bool {
	if t.kind() != kindVowel {
		return true
	}
	n := 0
	for _, v := range p.draw {
		if v.L == t.L {
			n++
		}
	}
	if n >= p.threshold {
		return false
	}
	p.draw = append(p.draw, t)

	return true
}
