package main

import (
	"fmt"
	"iter"
	"regexp"
	"strconv"
	"strings"
)

func (fl *Flags) AntNames(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.antNameRange == "" {
			yield("")
			return
		}
		if fl.debug {
			debug.WriteString("\nAntNames:")
		}

		minBitWidth, maxBitWidth := parseNameRange(fl.antNameRange)
		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			maxNum := uint64(1<<bitWidth) - 1
			for num := uint64(1); num < maxNum; num++ {
				name := numToName(num, bitWidth)
				if !yield(fmt.Sprintf(" -n %s", name)) {
					return
				}
				if fl.debug {
					debug.WriteString(" ")
					debug.WriteString(name)
				}
			}
		}
	}
}

func parseNameRange(antNameRange string) (int, int) {
	result := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(antNameRange)
	minBitWidth, _ := strconv.Atoi(result[1])
	maxBitWidth, _ := strconv.Atoi(result[2])

	return minBitWidth, maxBitWidth
}

func numToName(num uint64, bitWidth int) string {
	format := fmt.Sprintf("%%0%ds", bitWidth)
	binary := fmt.Sprintf(format, strconv.FormatUint(num, 2))
	return strings.Replace(strings.Replace(binary, "0", "L", -1), "1", "R", -1)
}
