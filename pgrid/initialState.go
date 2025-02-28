package pgrid

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
)

func (f *Field) initialState() (GridPoint, GridLine, GridLine, bool) {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := ParseInitialPoint(f.InitialPoint)

	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}
	prevLine := GridLine{Axis: uint8(prevAxis), Offset: offsetInt(prevOffset)}

	currPoint := f.makeGridPoint(currLine, prevLine)

	//fmt.Printf(
	//	"Initial step: %s %s %s %t\n",
	//	currPoint.String(), currLine.String(), prevLine.String(), prevPointSign,
	//)

	return currPoint, currLine, prevLine, prevPointSign
}

func (f *Field) initialStateString(currLine GridLine, prevLine GridLine, prevPointSign bool) string {
	prevPointSignString := "-"
	if prevPointSign {
		prevPointSignString = "+"
	}

	return fmt.Sprintf("%s%s%s", currLine.String(), prevPointSignString, prevLine.String())
}

func ParseInitialPoint(initialPoint string) (int, int, bool, int, int) {
	re := regexp.MustCompile(`([A-X])(-?\d+)([+-]?)([A-X])(-?\d+)`)
	result := re.FindStringSubmatch(initialPoint)

	currAx, currOff, dir, prevAx, prevOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := slices.Index(AxisNames, currAx)
	currOffset, _ := strconv.Atoi(currOff)

	currAxIncreasing := dir != "-"

	prevAxis := slices.Index(AxisNames, prevAx)
	prevOffset, _ := strconv.Atoi(prevOff)

	return currAxis, currOffset, currAxIncreasing, prevAxis, prevOffset
}
