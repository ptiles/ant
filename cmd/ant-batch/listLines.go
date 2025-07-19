package main

import (
	"fmt"
	"iter"
	"strings"
)

func (fl *Flags) ListLines(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.initialLines == "" {
			return
		}

		lines := strings.Split(fl.initialLines, ",")
		if fl.debug {
			debug.WriteString(fmt.Sprint("\nListLines: ", len(lines), lines))
		}

		for i, initialLine1 := range lines {
			for _, initialLine2 := range lines[i+1:] {
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
