package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
)

func (f *Field) initialState() (GridPoint, GridLine, GridLine, bool, uint8) {
	currAxis, currOffset, currAxIncreasing, nextAxis, nextOffset := utils.ParseInitialPoint(f.InitialPoint)

	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}
	nextLine := GridLine{Axis: uint8(nextAxis), Offset: offsetInt(nextOffset)}

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
