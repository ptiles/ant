package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid/parse"
	"github.com/ptiles/ant/seq"
	"github.com/ptiles/ant/utils"
	"iter"
	"slices"
	"strings"
)

type wythoff struct {
	axis1       string
	axis2       string
	axesPresent bool

	direction      string
	rangeMin       int
	rangeMax       int
	rangeDelta     int
	offsetsPresent bool
}

func (w *wythoff) intervalParser() flagParser {
	return func(interval string) error {
		if interval == "" || !w.axesPresent {
			return nil
		}

		rangeMin, rangeMax, rangeDelta, err := utils.ParseRangeDeltaStr(interval)
		if err != nil {
			return err
		}

		w.rangeMin = rangeMin
		w.rangeMax = rangeMax
		w.rangeDelta = rangeDelta
		w.offsetsPresent = true

		return nil
	}
}

func (w *wythoff) axesParser() flagParser {
	return func(axes string) error {
		if axes == "" {
			return nil
		}

		axis1, direction, axis2 := parse.InitialAxes(axes)
		if axis1 == "" || direction == "" || axis2 == "" {
			return nil
		}

		w.axis1 = axis1
		w.axis2 = axis2
		w.direction = direction
		w.axesPresent = true

		return nil
	}
}

func (w *wythoff) skip() bool {
	return !w.axesPresent || !w.offsetsPresent
}

func (w *wythoff) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if w.skip() {
			return
		}

		offsets := slices.Sorted(seq.WythoffDelta(w.rangeMin, w.rangeMax, w.rangeDelta))

		if debug != nil {
			debug.WriteString(fmt.Sprintf(
				"\nWythoffOffsets: %s%s%s [%d]%v",
				w.axis1, w.direction, w.axis2,
				len(offsets), offsets,
			))
		}

		for _, initialOffset1 := range offsets {
			for _, initialOffset2 := range offsets {
				yield(fmt.Sprintf(
					" -i %s%d%s%s%d",
					w.axis1, initialOffset1,
					w.direction,
					w.axis2, initialOffset2,
				))
			}
		}
	}
}
