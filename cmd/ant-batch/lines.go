package main

import (
	"fmt"
	"iter"
	"strings"
)

type lines struct {
	lines []string
}

func (l *lines) parser() flagParser {
	return func(lines string) error {
		if lines == "" {
			return nil
		}

		l.lines = strings.Split(lines, ",")

		return nil
	}
}

func (l *lines) skip() bool {
	return len(l.lines) == 0
}

func (l *lines) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if l.skip() {
			return
		}

		if debug != nil {
			debug.WriteString(fmt.Sprint("\nListLines: ", len(l.lines), l.lines))
		}

		for i, initialLine1 := range l.lines {
			for _, initialLine2 := range l.lines[i+1:] {
				if initialLine1[0] == initialLine2[0] {
					continue
				}

				for _, result := range [4]string{
					fmt.Sprintf(" -i %s%s%s", initialLine1, "+", initialLine2),
					fmt.Sprintf(" -i %s%s%s", initialLine1, "-", initialLine2),
					fmt.Sprintf(" -i %s%s%s", initialLine2, "+", initialLine1),
					fmt.Sprintf(" -i %s%s%s", initialLine2, "-", initialLine1),
				} {
					if !yield(result) {
						return
					}
				}
			}
		}
	}
}
