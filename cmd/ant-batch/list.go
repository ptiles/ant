package main

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"iter"
	"strings"
)

type list struct {
	axis1       string
	axis2       string
	direction   string
	axesPresent bool

	offsets []string
}

func (l *list) offsetsParser() flagParser {
	return func(offsets string) error {
		if offsets == "" {
			return nil
		}

		l.offsets = strings.Split(offsets, ",")

		return nil
	}
}

func (l *list) axesParser() flagParser {
	return func(axes string) error {
		if axes == "" {
			return nil
		}

		axis1, direction, axis2 := utils.ParseInitialAxes(axes)
		if axis1 == "" || direction == "" || axis2 == "" {
			return nil
		}

		l.axis1 = axis1
		l.axis2 = axis2
		l.direction = direction
		l.axesPresent = true

		return nil
	}
}

func (l *list) skip() bool {
	return !l.axesPresent || len(l.offsets) == 0
}

func (l *list) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if l.skip() {
			return
		}

		if debug != nil {
			debug.WriteString(fmt.Sprint("\nListOffsets: ", len(l.offsets), l.offsets))
		}

		for _, initialOffset1 := range l.offsets {
			for _, initialOffset2 := range l.offsets {
				if !yield(fmt.Sprintf(
					" -i %s%s%s%s%s",
					l.axis1, initialOffset1,
					l.direction,
					l.axis2, initialOffset2,
				)) {
					return
				}
			}
		}
	}
}
