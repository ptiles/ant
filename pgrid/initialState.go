package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/pgrid/axis"
	"image"
)

type Turn struct {
	CurrLine GridLine
	PrevLine GridLine
	sign     bool
}

func (f *Field) InitialTurn() Turn {
	currLine := GridLine{Axis: f.currAxis, Offset: OffsetInt(f.currOffset)}
	prevLine := GridLine{Axis: f.prevAxis, Offset: OffsetInt(f.prevOffset)}

	return Turn{
		CurrLine: currLine,
		PrevLine: prevLine,
		sign:     f.prevPointSign,
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
	return f.GetCenterPoint(GridAxes{
		Axis0: f.currAxis, Axis1: f.prevAxis,
		Coords: GridCoords{
			Offset0: OffsetInt(f.currOffset), Offset1: OffsetInt(f.prevOffset),
		},
	})
}

func (f *Field) SeedString(seedDropBits uint8) string {
	// Same seed for five symmetric points
	//return fmt.Sprintf(
	//	"%d%t%d%d",
	//	(int(GridLinesTotal)+currAxis-prevAxis)%int(GridLinesTotal),
	//	prevPointSign, currOffset>>seedDropBits, prevOffset>>seedDropBits,
	//)

	// Different seeds
	return fmt.Sprintf(
		"%s%s%t%d%d",
		axis.Name[f.currAxis], axis.Name[f.prevAxis],
		f.prevPointSign, f.currOffset>>seedDropBits, f.prevOffset>>seedDropBits,
	)
}
