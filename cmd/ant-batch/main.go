package main

import (
	"fmt"
	"iter"
	"os"
	"slices"
	"strings"
)

func main() {
	fl := parseFlags()

	var debug *strings.Builder

	if fl.debug {
		debug = &strings.Builder{}
	}

	antNames := slices.Collect(fl.names.seq(debug))
	antPatterns := slices.Collect(fl.patterns.seq(debug))

	initialPoints := concatSeq(
		fl.grid.seq(debug),
		fl.interval.seq(debug),
		fl.lines.seq(debug),
		fl.list.seq(debug),
		fl.near.seq(debug),
		fl.path.seq(debug),
		fl.rect.seq(debug),
		fl.wythoff.seq(debug),
	)

	if debug != nil {
		fmt.Fprintln(os.Stderr, "Values used:", debug.String())
	}

	for initialPoint := range initialPoints {
		for _, pattern := range antPatterns {
			for _, name := range antNames {
				fmt.Print(name, pattern, initialPoint, "\n")
			}
		}
	}

	if debug != nil {
		fmt.Fprintln(os.Stderr, "Values used:", debug.String())
	}
}

func concatSeq(seqs ...iter.Seq[string]) iter.Seq[string] {
	return func(yield func(string) bool) {
		empty := true
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
				empty = false
			}
		}
		if empty {
			yield("")
		}
	}
}
