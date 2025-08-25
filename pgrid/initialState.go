package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"image"
)

type Turn struct {
	CurrLine GridLine
	PrevLine GridLine
	sign     bool
}

func (f *Field) InitialTurn() Turn {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := utils.ParseInitialPoint(f.InitialPoint)

	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}
	prevLine := GridLine{Axis: uint8(prevAxis), Offset: offsetInt(prevOffset)}

	return Turn{
		CurrLine: currLine,
		PrevLine: prevLine,
		sign:     prevPointSign,
	}
}

func (t Turn) String() string {
	prevPointSignString := "-"
	if t.sign {
		prevPointSignString = "+"
	}
	return fmt.Sprintf("%s%s%s", t.CurrLine.String(), prevPointSignString, t.PrevLine.String())
}

func (f *Field) InitialPointOutside(r image.Rectangle) bool {
	return !f.InitialCenterPoint().In(r)
}

func (f *Field) InitialCenterPoint() image.Point {
	currAxis, currOffset, _, prevAxis, prevOffset := utils.ParseInitialPoint(f.InitialPoint)
	return f.GetCenterPoint(GridAxes{
		Axis0: uint8(currAxis), Axis1: uint8(prevAxis),
		Coords: GridCoords{
			Offset0: offsetInt(currOffset), Offset1: offsetInt(prevOffset),
		},
	})
}

func (f *Field) SeedString(seedDropBits uint8) string {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := utils.ParseInitialPoint(f.InitialPoint)

	// Same seed for five symmetric points
	//return fmt.Sprintf(
	//	"%d%t%d%d",
	//	(int(GridLinesTotal)+currAxis-prevAxis)%int(GridLinesTotal),
	//	prevPointSign, currOffset>>seedDropBits, prevOffset>>seedDropBits,
	//)

	// Different seeds
	return fmt.Sprintf(
		"%s%s%t%d%d",
		AxisNames[currAxis], AxisNames[prevAxis],
		prevPointSign, currOffset>>seedDropBits, prevOffset>>seedDropBits,
	)
}
