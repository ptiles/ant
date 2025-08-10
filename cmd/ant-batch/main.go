package main

import (
	"fmt"
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

	var initialPoints []string
	initialPoints = slices.AppendSeq(initialPoints, fl.interval.seq(debug))
	initialPoints = slices.AppendSeq(initialPoints, fl.lines.seq(debug))
	initialPoints = slices.AppendSeq(initialPoints, fl.list.seq(debug))
	initialPoints = slices.AppendSeq(initialPoints, fl.near.seq(debug))
	initialPoints = slices.AppendSeq(initialPoints, fl.path.seq(debug))
	initialPoints = slices.AppendSeq(initialPoints, fl.wythoff.seq(debug))
	if len(initialPoints) == 0 {
		initialPoints = []string{""}
	}

	for _, name := range antNames {
		for _, pattern := range antPatterns {
			for _, initialPoint := range initialPoints {
				fmt.Print(name, pattern, initialPoint, "\n")
			}
		}
	}

	if debug != nil {
		fmt.Fprintln(os.Stderr, "Values used:", debug.String())
	}
}
