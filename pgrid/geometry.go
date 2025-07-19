package pgrid

import (
	"github.com/ptiles/ant/geom"
)

type Geometry struct {
	offsetsToFirst allOffsetDeltas
	offsetsToLast  allOffsetDeltas
}

type gridGeometry struct {
	anchors [GridLinesTotal]geom.Point
	normals [GridLinesTotal]geom.Point
	units   [GridLinesTotal]geom.Point
}

func newGridGeometry(radius float64) Geometry {
	gg := gridGeometry{}

	alpha := 360 / float64(GridLinesTotal)
	rightAngle := float64(90)

	for ax := range GridLinesTotal {
		alphaAx := alpha * float64(ax)

		gg.anchors[ax].X = radius * geom.Cos(alphaAx)
		gg.anchors[ax].Y = radius * geom.Sin(alphaAx)

		gg.normals[ax].X = geom.Cos(alphaAx + alpha/2)
		gg.normals[ax].Y = geom.Sin(alphaAx + alpha/2)

		gg.units[ax].X = geom.Cos(alphaAx + alpha/2 + rightAngle)
		gg.units[ax].Y = geom.Sin(alphaAx + alpha/2 + rightAngle)
	}

	offsetsToFirst := gg.newOffsetsToFirst()
	offsetsToLast := gg.newOffsetsToLast()

	if GridLinesTotal == 5 {
		//printOffsets("offsetsToFirst", offsetsToFirst)
		//printOffsets("offsetsToLast", offsetsToLast)
		//os.Exit(1)
		updateOffsetsToFirst(&offsetsToFirst)
		updateOffsetsToLast(&offsetsToLast)
	}

	return Geometry{offsetsToFirst: offsetsToFirst, offsetsToLast: offsetsToLast}
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

	for ax0, ax1 := range AxesAll() {
		for i, axT := range otherAxes(ax0, ax1) {
			result[ax0][ax1][i] = gg.newOffsetDeltas(ax1, axT, ax0)
			result[ax0][ax1][i].targetAx = axT
		}
	}

	return result
}

func (gg *gridGeometry) newOffsetsToLast() allOffsetDeltas {
	result := allOffsetDeltas{}

	for ax0, ax1 := range AxesAll() {
		for i, axT := range otherAxes(ax0, ax1) {
			result[ax0][ax1][i] = gg.newOffsetDeltas(ax0, ax1, axT)
			result[ax0][ax1][i].targetAx = axT
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

	axA0B0Point := geom.Intersection(axA0Line, axB0Line)
	axA0B0Offset := geom.Distance(axT0Line, axA0B0Point)

	axA1B0Point := geom.Intersection(axA1Line, axB0Line)
	axA1Delta := geom.Distance(axT0Line, axA1B0Point) - axA0B0Offset

	axA0B1Point := geom.Intersection(axA0Line, axB1Line)
	axB1Delta := geom.Distance(axT0Line, axA0B1Point) - axA0B0Offset

	return offsetDeltas{zeroZero: axA0B0Offset, ax0Delta: axA1Delta, ax1Delta: axB1Delta}
}

func (gg *gridGeometry) getLine(gl GridLine) geom.Line {
	anchor := gg.anchors[gl.Axis]
	normal := gg.normals[gl.Axis]
	unit := gg.units[gl.Axis]
	offset := float64(gl.Offset)

	pointA := geom.Point{X: anchor.X + normal.X*offset, Y: anchor.Y + normal.Y*offset}
	pointB := geom.Point{X: pointA.X + unit.X, Y: pointA.Y + unit.Y}
	return geom.Line{A: pointA, B: pointB}
}
