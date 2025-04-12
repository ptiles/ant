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

func New(r float64, rules []bool, initialPoint string) *Field {
	gg := newGridGeometry(r)

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

func (f *Field) makeGridPoint(gridLine0, gridLine1 GridLine) GridPoint {
	if gridLine0.Axis > gridLine1.Axis {
		gridLine0, gridLine1 = gridLine1, gridLine0
	}

	gridPoint := GridPoint{
		Axes: GridAxes{
			Axis0: gridLine0.Axis, Axis1: gridLine1.Axis,
			Coords: GridCoords{
				Offset0: gridLine0.Offset, Offset1: gridLine1.Offset,
			},
		},
	}

	gridPoint.Offsets[gridLine0.Axis] = gridLine0.Offset
	gridPoint.Offsets[gridLine1.Axis] = gridLine1.Offset

	off0, off1 := float64(gridLine0.Offset), float64(gridLine1.Offset)
	for _, otl := range f.offsetsToLast[gridLine0.Axis][gridLine1.Axis] {
		off := otl.zeroZero + off0*otl.ax0Delta + off1*otl.ax1Delta
		gridPoint.Offsets[otl.targetAx] = offsetInt(math.Ceil(off))
	}

	return gridPoint
}

func (f *Field) nearestNeighbor(
	currentPointOffsets GridOffsets,
	prevLine, currentLine GridLine,
	positiveSide bool,
) (GridPoint, GridLine, bool) {
	var nextLine GridLine
	var nextPointSign bool
	currentDistance := 1000000.0
	prevLineOffset := float64(prevLine.Offset)
	currentLineOffset := float64(currentLine.Offset)

	for _, otf := range f.offsetsToFirst[prevLine.Axis][currentLine.Axis] {
		nextAxis := otf.targetAx

		nextOffset := currentPointOffsets[nextAxis]
		dist := otf.zeroZero + currentLineOffset*otf.ax0Delta + float64(nextOffset)*otf.ax1Delta - prevLineOffset

		if (positiveSide && dist > 0) || (!positiveSide && dist < 0) {
			absDist := math.Abs(dist)
			if absDist < currentDistance {
				currentDistance = absDist
				nextLine = GridLine{nextAxis, nextOffset}
				nextPointSign = true
			}
		}

		nextOffset -= 1
		dist -= otf.ax1Delta

		if (positiveSide && dist > 0) || (!positiveSide && dist < 0) {
			absDist := math.Abs(dist)
			if absDist < currentDistance {
				currentDistance = absDist
				nextLine = GridLine{nextAxis, nextOffset}
				nextPointSign = false
			}
		}
	}

	nextPoint := f.makeGridPoint(currentLine, nextLine)
	return nextPoint, nextLine, nextPointSign
}
