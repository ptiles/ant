package pgrid

import (
	"github.com/ptiles/ant/utils"
)

func (f *Field) initialState() (GridPoint, GridLine, GridLine, bool, uint8) {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := utils.ParseInitialPoint(f.InitialPoint)

	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}
	prevLine := GridLine{Axis: uint8(prevAxis), Offset: offsetInt(prevOffset)}

	currPoint := f.makeGridPoint(currLine, prevLine, f.gridPointToPoint(currLine, prevLine))

	//fmt.Printf(
	//	"Initial step: %s %s %s %t\n",
	//	currPoint.String(), currLine.String(), prevLine.String(), prevPointSign,
	//)

	return currPoint, currLine, prevLine, prevPointSign, 1
}
