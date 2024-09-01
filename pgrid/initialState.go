package pgrid

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (f *Field) initialState() (GridPoint, GridPoint, GridLine, GridLine, uint8) {
	re := regexp.MustCompile(`([A-E])(-?\d+)([+-]?)([A-E])(-?\d+)`)
	result := re.FindStringSubmatch(f.InitialPoint)

	currAx, currOff, dir, nextAx, nextOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := strings.Index("ABCDE", currAx)
	currOffset, _ := strconv.Atoi(currOff)
	currLine := GridLine{Axis: uint8(currAxis), Offset: int16(currOffset)}

	nextAxis := strings.Index("ABCDE", nextAx)
	nextOffset, _ := strconv.Atoi(nextOff)
	nextLine := GridLine{Axis: uint8(nextAxis), Offset: int16(nextOffset)}

	currAxIncreasing := dir != "-"
	currPointPoint := f.gridPointToPoint(currLine, nextLine)
	currPoint := f.makeGridPoint(currLine, nextLine, currPointPoint)

	prevPoint, prevLine := f.nearestNeighbor(currPoint, nextLine, currLine, !currAxIncreasing)

	if f.verbose {
		axisNames := [5]string{"A", "B", "C", "D", "E"}
		fmt.Printf("Initial step: ")
		fmt.Printf("%s%d%s%d=>", axisNames[currLine.Axis], currLine.Offset, axisNames[prevLine.Axis], prevLine.Offset)
		fmt.Printf("%s%d%s%d\n", axisNames[nextLine.Axis], nextLine.Offset, axisNames[currLine.Axis], currLine.Offset)
	}

	return prevPoint, currPoint, prevLine, currLine, 0
}
