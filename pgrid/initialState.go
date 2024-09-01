package pgrid

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (f *Field) initialState() (GridPoint, GridLine, GridLine, bool, uint8) {
	re := regexp.MustCompile(`([A-X])(-?\d+)([+-]?)([A-X])(-?\d+)`)
	result := re.FindStringSubmatch(f.InitialPoint)

	currAx, currOff, dir, nextAx, nextOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := strings.Index(AxisCharacters, currAx)
	currOffset, _ := strconv.Atoi(currOff)
	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}

	nextAxis := strings.Index(AxisCharacters, nextAx)
	nextOffset, _ := strconv.Atoi(nextOff)
	nextLine := GridLine{Axis: uint8(nextAxis), Offset: offsetInt(nextOffset)}

	currAxIncreasing := dir != "-"
	currPointPoint := f.gridPointToPoint(currLine, nextLine)
	currPoint := f.makeGridPoint(currLine, nextLine, currPointPoint)

	prevPoint, prevLine := f.nearestNeighbor(currPoint.Offsets, nextLine, currLine, !currAxIncreasing)
	prevPointSign := distance(f.getLine(currLine), prevPoint.Point) < 0

	if f.verbose {
		fmt.Printf("Initial step: ")
		fmt.Printf("%s%d%s%d=>", AxisNames[currLine.Axis], currLine.Offset, AxisNames[prevLine.Axis], prevLine.Offset)
		fmt.Printf("%s%d%s%d\n", AxisNames[nextLine.Axis], nextLine.Offset, AxisNames[currLine.Axis], currLine.Offset)
	}

	return currPoint, currLine, prevLine, prevPointSign, 0
}
