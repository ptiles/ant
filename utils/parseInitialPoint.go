package utils

import (
	"regexp"
	"strconv"
	"strings"
)

const AxisCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXY"

func ParseInitialPoint(initialPoint string) (int, int, bool, int, int) {
	re := regexp.MustCompile(`([A-X])(-?\d+)([+-]?)([A-X])(-?\d+)`)
	result := re.FindStringSubmatch(initialPoint)

	currAx, currOff, dir, prevAx, prevOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := strings.Index(AxisCharacters, currAx)
	currOffset, _ := strconv.Atoi(currOff)

	currAxIncreasing := dir != "-"

	prevAxis := strings.Index(AxisCharacters, prevAx)
	prevOffset, _ := strconv.Atoi(prevOff)

	return currAxis, currOffset, currAxIncreasing, prevAxis, prevOffset
}
