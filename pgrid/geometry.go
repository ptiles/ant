package pgrid

import (
	"iter"
	"math"
)

type Point [2]float64

type Line [2]Point

type gridGeometry struct {
	anchors [GridLinesTotal]Point
	normals [GridLinesTotal]Point
	units   [GridLinesTotal]Point
}

const X = 0
const Y = 1

func newGridGeometry(r float64) gridGeometry {
	result := gridGeometry{}

	phi := 2 * math.Pi / float64(GridLinesTotal)
	rightAngle := 0.5 * math.Pi

	for ax := range GridLinesTotal {
		phiAx := phi * float64(ax)

		result.anchors[ax][X] = r * math.Cos(phiAx)
		result.anchors[ax][Y] = r * math.Sin(phiAx)

		result.normals[ax][X] = math.Cos(phiAx + phi/2)
		result.normals[ax][Y] = math.Sin(phiAx + phi/2)

		result.units[ax][X] = math.Cos(phiAx + phi/2 + rightAngle)
		result.units[ax][Y] = math.Sin(phiAx + phi/2 + rightAngle)
	}

	return result
}

func otherAxes(ax0, ax1 uint8) iter.Seq2[uint8, uint8] {
	return func(yield func(uint8, uint8) bool) {
		if ax0 == ax1 {
			return
		}

		i := uint8(0)
		for ax := range GridLinesTotal {
			if ax == ax0 || ax == ax1 {
				continue
			}
			if !yield(i, ax) {
				return
			}
			i += 1
		}
	}
}

type offsetDeltas struct {
	targetAx uint8
	zeroZero float64
	ax0Delta float64
	ax1Delta float64
}
type allOffsetDeltas [GridLinesTotal][GridLinesTotal][GridLinesTotal - 2]offsetDeltas

func (gg *gridGeometry) newOffsetsToFirst() allOffsetDeltas {
	result := allOffsetDeltas{}

	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			for i, axT := range otherAxes(ax0, ax1) {
				result[ax0][ax1][i] = gg.newOffsetDeltas(ax1, axT, ax0)
				result[ax0][ax1][i].targetAx = axT
			}
		}
	}

	return result
}

func (gg *gridGeometry) newOffsetsToLast() allOffsetDeltas {
	result := allOffsetDeltas{}

	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			for i, axT := range otherAxes(ax0, ax1) {
				result[ax0][ax1][i] = gg.newOffsetDeltas(ax0, ax1, axT)
				result[ax0][ax1][i].targetAx = axT
			}
		}
	}

	return result
}

func (gg *gridGeometry) newOffsetDeltas(axA, axB, axT uint8) offsetDeltas {
	axA0Line := gg.getLine(GridLine{axA, 0})
	axA1Line := gg.getLine(GridLine{axA, 1})
	axB0Line := gg.getLine(GridLine{axB, 0})
	axB1Line := gg.getLine(GridLine{axB, 1})
	axT0Line := gg.getLine(GridLine{axT, 0})

	axA0B0Point := intersection(axA0Line, axB0Line)
	axA0B0Offset := distance(axT0Line, axA0B0Point)

	axA1B0Point := intersection(axA1Line, axB0Line)
	axA1Delta := distance(axT0Line, axA1B0Point) - axA0B0Offset

	axA0B1Point := intersection(axA0Line, axB1Line)
	axB1Delta := distance(axT0Line, axA0B1Point) - axA0B0Offset

	return offsetDeltas{zeroZero: axA0B0Offset, ax0Delta: axA1Delta, ax1Delta: axB1Delta}
}

func (gg *gridGeometry) getLine(gl GridLine) Line {
	anchor := gg.anchors[gl.Axis]
	normal := gg.normals[gl.Axis]
	unit := gg.units[gl.Axis]
	offset := float64(gl.Offset)

	point1 := Point{anchor[X] + normal[X]*offset, anchor[Y] + normal[Y]*offset}
	point2 := Point{point1[X] + unit[X], point1[Y] + unit[Y]}
	return Line{point1, point2}
}

func intersection(line1, line2 Line) Point {
	line1point1, line1point2 := line1[0], line1[1]
	line2point1, line2point2 := line2[0], line2[1]

	x1, y1 := line1point1[X], line1point1[Y]
	x2, y2 := line1point2[X], line1point2[Y]
	x3, y3 := line2point1[X], line2point1[Y]
	x4, y4 := line2point2[X], line2point2[Y]

	dx12 := x1 - x2
	dy12 := y1 - y2
	dx34 := x3 - x4
	dy34 := y3 - y4

	den := dx12*dy34 - dy12*dx34

	return Point{
		((x1*y2-y1*x2)*dx34 - dx12*(x3*y4-y3*x4)) / den,
		((x1*y2-y1*x2)*dy34 - dy12*(x3*y4-y3*x4)) / den,
	}
}

func distance(line Line, point Point) float64 {
	x1, y1 := line[0][X], line[0][Y]
	x2, y2 := line[1][X], line[1][Y]

	x0, y0 := point[X], point[Y]

	x10 := x1 - x0
	y10 := y1 - y0
	x21 := x2 - x1
	y21 := y2 - y1

	return x21*y10 - x10*y21
}
