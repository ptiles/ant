package pgrid

import (
	"math"
)

type Field struct {
	Rules        []bool
	Limit        uint8
	InitialPoint string

	offsetsToFirst allOffsetDeltas
	offsetsToLast  allOffsetDeltas
}

func New(radius float64, rules []bool, initialPoint string) *Field {
	gg := newGridGeometry(radius)

	return &Field{
		Rules:        rules,
		Limit:        uint8(len(rules)),
		InitialPoint: initialPoint,

		offsetsToFirst: gg.newOffsetsToFirst(),
		offsetsToLast:  gg.newOffsetsToLast(),
	}
}

type GridLine struct {
	Axis   uint8
	Offset offsetInt
}

type GridPoint struct {
	Axes    GridAxes
	Offsets GridOffsets
}

type offsetInt int32

type GridAxes struct {
	Axis0  uint8
	Axis1  uint8
	Coords GridCoords
}

type GridCoords struct {
	Offset0 offsetInt
	Offset1 offsetInt
}

type GridOffsets [GridLinesTotal]offsetInt

func (f *Field) nearestNeighbor(
	prevLine, currentLine GridLine,
	positiveSide bool,
) (GridLine, bool) {
	var nextLine GridLine
	var nextPointSign bool
	currentDistance := 1000000.0
	prevLineOffset := float64(prevLine.Offset)
	currentLineOffset := float64(currentLine.Offset)

	for i, otf := range f.offsetsToFirst[prevLine.Axis][currentLine.Axis] {
		otl := f.offsetsToLast[prevLine.Axis][currentLine.Axis][i]
		off := otl.zeroZero + prevLineOffset*otl.ax0Delta + currentLineOffset*otl.ax1Delta
		nextLineOffset := math.Ceil(off)

		dist := otf.zeroZero + currentLineOffset*otf.ax0Delta + nextLineOffset*otf.ax1Delta - prevLineOffset

		if (positiveSide && dist > 0) || (!positiveSide && dist < 0) {
			absDist := math.Abs(dist)
			if absDist < currentDistance {
				currentDistance = absDist
				nextLine = GridLine{otf.targetAx, offsetInt(nextLineOffset)}
				nextPointSign = true
			}
		} else {
			absDist := math.Abs(dist - otf.ax1Delta)
			if absDist < currentDistance {
				currentDistance = absDist
				nextLine = GridLine{otf.targetAx, offsetInt(nextLineOffset) - 1}
				nextPointSign = false
			}
		}
	}

	return nextLine, nextPointSign
}
