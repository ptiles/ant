package main

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"strings"
)

func (fl *Flags) Patterns(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.patternsCount == 0 {
			yield("")
			return
		}
		if fl.debug {
			debug.WriteString("\nPatterns:")
		}

		precision := uint(10_000)
		for range fl.patternsCount {
			pattern := fmt.Sprintf("%f",
				float64(precision-rand.UintN(precision))/float64(precision),
			)
			if !yield(fmt.Sprintf(" -p %s", pattern)) {
				return
			}
			if fl.debug {
				debug.WriteString(" ")
				debug.WriteString(pattern)
			}
		}
	}
}
