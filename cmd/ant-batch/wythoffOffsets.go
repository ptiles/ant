package main

import (
	"fmt"
	"github.com/ptiles/ant/seq"
	"github.com/ptiles/ant/utils"
	"iter"
	"slices"
	"strings"
)

func (fl *Flags) WythoffOffsets(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.initialPointWythoff == "" || fl.initialDirection == "" ||
			fl.initialAxis1 == "" || fl.initialAxis2 == "" {
			return
		}

		rangeMin, rangeMax, rangeDelta, _ := utils.ParseRangeStr(fl.initialPointWythoff)
		offsets := slices.Sorted(seq.WythoffDelta(rangeMin, rangeMax, rangeDelta))

		if fl.debug {
			debug.WriteString(fmt.Sprintf(
				"\nWythoffOffsets: %s%s%s %d:%v",
				fl.initialAxis1, fl.initialDirection, fl.initialAxis2,
				len(offsets), offsets,
			))
		}

		for _, initialOffset1 := range offsets {
			for _, initialOffset2 := range offsets {
				yield(fmt.Sprintf(
					" -i %s%d%s%s%d",
					fl.initialAxis1, initialOffset1,
					fl.initialDirection,
					fl.initialAxis2, initialOffset2,
				))
			}
		}
	}
}
