package main

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"iter"
	"strconv"
	"strings"
)

type names struct {
	rangeMin int
	rangeMax int
	present  bool
}

func (n *names) parser() flagParser {
	return func(interval string) error {
		if interval == "" {
			return nil
		}

		rangeMin, rangeMax, err := utils.ParseRangeStr(interval)
		if err != nil {
			return err
		}

		if rangeMin == 0 {
			rangeMin = rangeMax
		}

		n.rangeMin = rangeMin
		n.rangeMax = rangeMax
		n.present = true

		return nil
	}
}

func (n *names) skip() bool {
	return !n.present
}

func (n *names) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if n.skip() {
			yield("")
			return
		}

		if debug != nil {
			debug.WriteString("\nAntNames:")
		}

		minBitWidth, maxBitWidth := n.rangeMin, n.rangeMax

		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			maxNum := uint64(1<<bitWidth) - 1
			for num := uint64(1); num < maxNum; num++ {
				name := numToName(num, bitWidth)
				if !yield(fmt.Sprintf(" -n %s", name)) {
					return
				}
				if debug != nil {
					debug.WriteString(" ")
					debug.WriteString(name)
				}
			}
		}
	}
}

func numToName(num uint64, bitWidth int) string {
	format := fmt.Sprintf("%%0%ds", bitWidth)
	binary := fmt.Sprintf(format, strconv.FormatUint(num, 2))
	return strings.Replace(strings.Replace(binary, "0", "L", -1), "1", "R", -1)
}
