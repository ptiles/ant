package parse

import (
	"github.com/ptiles/ant/pgrid/axis"
	"regexp"
	"strconv"
)

func InitialPoint(initialPoint string) (uint8, int, bool, uint8, int) {
	re := regexp.MustCompile(`([A-X])(-?\d+)([:+-])([A-X])(-?\d+)`)
	result := re.FindStringSubmatch(initialPoint)

	currAx, currOff, dir, prevAx, prevOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := axis.Index(currAx)
	currOffset, _ := strconv.Atoi(currOff)

	currAxIncreasing := dir != "-"

	prevAxis := axis.Index(prevAx)
	prevOffset, _ := strconv.Atoi(prevOff)

	return currAxis, currOffset, currAxIncreasing, prevAxis, prevOffset
}

func InitialAxes(initialPoint string) (string, string, string) {
	re := regexp.MustCompile(`([A-X])([+-])([A-X])`)
	result := re.FindStringSubmatch(initialPoint)

	return result[1], result[2], result[3]
}
