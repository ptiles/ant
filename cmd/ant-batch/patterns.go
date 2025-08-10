package main

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"strconv"
	"strings"
)

type patterns struct {
	count int
}

func (p *patterns) parser() flagParser {
	return func(pattern string) error {
		if pattern == "" {
			return nil
		}

		patternsCount, err := strconv.Atoi(pattern)
		if err != nil {
			return err
		}

		p.count = patternsCount

		return nil
	}
}

func (p *patterns) skip() bool {
	return p.count == 0
}

func (p *patterns) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if p.skip() {
			yield("")
			return
		}
		if debug != nil {
			debug.WriteString("\nPatterns:")
		}

		precision := uint(10_000)
		for range p.count {
			pattern := fmt.Sprintf("%f",
				float64(precision-rand.UintN(precision))/float64(precision),
			)
			if !yield(fmt.Sprintf(" -p %s", pattern)) {
				return
			}
			if debug != nil {
				debug.WriteString(" ")
				debug.WriteString(pattern)
			}
		}
	}
}
