package main

import (
	"fmt"
	"iter"
	"strings"
)

func (fl *Flags) ListOffsets(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.initialOffsets == "" || fl.initialDirection == "" ||
			fl.initialAxis1 == "" || fl.initialAxis2 == "" {
			return
		}

		offsets := strings.Split(fl.initialOffsets, ",")
		if fl.debug {
			debug.WriteString(fmt.Sprint("\nListOffsets: ", len(offsets), offsets))
		}

		for _, initialOffset1 := range offsets {
			for _, initialOffset2 := range offsets {
				yield(fmt.Sprintf(
					" -i %s%s%s%s%s",
					fl.initialAxis1, initialOffset1,
					fl.initialDirection,
					fl.initialAxis2, initialOffset2,
				))
			}
		}
	}
}
