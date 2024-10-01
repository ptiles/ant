package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
)

func (f *Field) initialState() (GridPoint, GridLine, GridLine, bool) {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := utils.ParseInitialPoint(f.InitialPoint)

	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}
	prevLine := GridLine{Axis: uint8(prevAxis), Offset: offsetInt(prevOffset)}

	currPoint := f.makeGridPoint(currLine, prevLine, f.gridPointToPoint(currLine, prevLine))

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
