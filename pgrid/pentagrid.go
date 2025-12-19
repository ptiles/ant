package pgrid

import (
	"github.com/ptiles/ant/pgrid/parse"
)

type Field struct {
	geometry      Geometry
	Rules         []bool
	Limit         uint8
	currAxis      uint8
	currOffset    int
	prevPointSign bool
	prevAxis      uint8
	prevOffset    int
}

func New(radius float64, rules []bool, initialPoint string) *Field {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := parse.InitialPoint(initialPoint)

	return &Field{
		Rules:         rules,
		Limit:         uint8(len(rules)),
		currAxis:      currAxis,
		currOffset:    currOffset,
		prevPointSign: prevPointSign,
		prevAxis:      prevAxis,
		prevOffset:    prevOffset,
		geometry:      newGeometry(radius),
	}
}

type GridLine struct {
	Axis   uint8
	Offset OffsetInt
}

type GridPoint struct {
	Axes    GridAxes
	Offsets GridOffsets
}

type OffsetInt int32

type GridAxes struct {
	Axis0  uint8
	Axis1  uint8
	Coords GridCoords
}

type GridCoords struct {
	Offset0 OffsetInt
	Offset1 OffsetInt
}

type GridOffsets [GridLinesTotal]OffsetInt
