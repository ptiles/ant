package utils

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var AxisNames = []string{
	"A", "B", "C", "D", "E",
	"F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y",
}

func ParseInitialPoint(initialPoint string) (int, int, bool, int, int) {
	re := regexp.MustCompile(`([A-X])(-?\d+)([+-])([A-X])(-?\d+)`)
	result := re.FindStringSubmatch(initialPoint)

	currAx, currOff, dir, prevAx, prevOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := slices.Index(AxisNames, currAx)
	currOffset, _ := strconv.Atoi(currOff)

	currAxIncreasing := dir != "-"

	prevAxis := slices.Index(AxisNames, prevAx)
	prevOffset, _ := strconv.Atoi(prevOff)

	return currAxis, currOffset, currAxIncreasing, prevAxis, prevOffset
}

func ParseInitialAxes(initialPoint string) (string, string, string) {
	re := regexp.MustCompile(`([A-X])([+-])([A-X])`)
	result := re.FindStringSubmatch(initialPoint)

	return result[1], result[2], result[3]
}

func ParseStepsStr(stepsStr string) (uint64, uint64, uint64, error) {
	re := regexp.MustCompile(`((\d+)-)?(\d+)(%(\d+))?`)
	result := re.FindStringSubmatch(strings.Replace(stepsStr, "_", "", -1))

	minSt, _ := strconv.ParseUint(result[2], 0, 64)
	maxSt, _ := strconv.ParseUint(result[3], 0, 64)
	incSt, _ := strconv.ParseUint(result[5], 0, 64)

	return minSt, maxSt, incSt, nil
}
