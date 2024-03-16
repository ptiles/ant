package main

import (
	"fmt"
	"strconv"
	"strings"
)

func numToName(num uint64, bitWidth int) string {
	format := fmt.Sprintf("%%0%ds", bitWidth)
	binary := fmt.Sprintf(format, strconv.FormatUint(num, 2))
	return strings.Replace(strings.Replace(binary, "0", "L", -1), "1", "R", -1)
}

func main() {
	initialPoint := "A1421B1790"
	maxSteps := 20000000
	maxBitWidth := 7

	for bitWidth := 3; bitWidth <= maxBitWidth; bitWidth++ {
		for num := uint64(1); num < 1<<bitWidth-1; num++ {
			name := numToName(num, bitWidth)
			fmt.Printf("-n %s -i %s -s %d\n", name, initialPoint, maxSteps)
		}
	}
}
