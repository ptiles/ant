package pgrid

import (
	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/store"
	"math"
)

func fromDegrees(deg int) float64 {
	return float64(deg) * math.Pi / 180.0
}

const X = 0
const Y = 1

const A = 0
const B = 1
const C = 2
const D = 3
const E = 4

var axisIndices = [5]uint8{A, B, C, D, E}

type Field struct {
	dist             float64
	anchors          [5]geom.Point
	anchorLines      [5]geom.Line
	axisUnits        [5]geom.Point
	normals          [5]geom.Point
	intersectAnchors [5][5]geom.Point
	intersectVectors [5][5]geom.Point
}

func New(r, dist float64, phi0degrees int) Field {
	phi0 := fromDegrees(phi0degrees)
	phi := fromDegrees(72)
	axisAngle0 := fromDegrees(180 - 54 + phi0degrees)

	result := Field{dist: dist}

	for ax := range axisIndices {
		phiAx := phi * float64(ax)

		result.anchors[ax][X] = r * math.Cos(phi0+phiAx)
		result.anchors[ax][Y] = r * math.Sin(phi0+phiAx)

		result.axisUnits[ax][X] = 1 * math.Cos(axisAngle0+phiAx)
		result.axisUnits[ax][Y] = 1 * math.Sin(axisAngle0+phiAx)

		result.normals[ax][X] = dist * math.Cos(phi0+0.5*phi+phiAx)
		result.normals[ax][Y] = dist * math.Sin(phi0+0.5*phi+phiAx)

		anchorEnd := geom.Point{
			result.anchors[ax][X] + result.axisUnits[ax][X],
			result.anchors[ax][Y] + result.axisUnits[ax][Y],
		}
		result.anchorLines[ax] = geom.Line{result.anchors[ax], anchorEnd}
	}

	for ax0 := uint8(0); ax0 < 5; ax0++ {
		for ax1 := uint8(0); ax1 < 5; ax1++ {
			line0 := result.GetLine(GridLine{ax0, 0})
			line1 := result.GetLine(GridLine{ax1, 0})
			intersectAnchor := geom.Intersection(line0, line1)
			result.intersectAnchors[ax0][ax1] = intersectAnchor

			line01 := result.GetLine(GridLine{ax0, 1})
			intersection01 := geom.Intersection(line01, line1)
			result.intersectVectors[ax0][ax1] = geom.Point{
				intersection01[X] - intersectAnchor[X],
				intersection01[Y] - intersectAnchor[Y],
			}
		}
	}

	return result
}

type GridLine struct {
	Axis   uint8
	Offset int16
	//Line   geom.Line
}

func (f *Field) GetLine(gl GridLine) geom.Line {
	axis := f.axisUnits[gl.Axis]
	anchor := f.anchors[gl.Axis]
	normal := f.normals[gl.Axis]
	distance := float64(gl.Offset)

	point1 := geom.Point{anchor[X] + normal[X]*distance, anchor[Y] + normal[Y]*distance}
	point2 := geom.Point{point1[X] + axis[X], point1[Y] + axis[Y]}
	return geom.Line{point1, point2}
}

type GridPoint struct {
	axes         [5]bool
	offsets      [5]float64
	Point        geom.Point
	PackedCoords store.PackedCoordinates
}

func (f *Field) MakeGridPoint(gridLine0, gridLine1 GridLine) GridPoint {
	gridPoint := GridPoint{}
	offset0 := float64(gridLine0.Offset)
	gridPoint.offsets[gridLine0.Axis] = offset0
	gridPoint.axes[gridLine0.Axis] = true
	offset1 := float64(gridLine1.Offset)
	gridPoint.offsets[gridLine1.Axis] = offset1
	gridPoint.axes[gridLine1.Axis] = true

	//gridPoint.Point = geom.Intersection(gridLine0.Line, gridLine1.Line)
	anchor := f.intersectAnchors[gridLine0.Axis][gridLine1.Axis]

	vector0 := f.intersectVectors[gridLine0.Axis][gridLine1.Axis]
	vector1 := f.intersectVectors[gridLine1.Axis][gridLine0.Axis]
	gridPoint.Point = geom.Point{
		anchor[X] + vector0[X]*offset0 + vector1[X]*offset1,
		anchor[Y] + vector0[Y]*offset0 + vector1[Y]*offset1,
	}

	for ax := range gridPoint.offsets {
		if !gridPoint.axes[ax] {
			distance := geom.Distance(f.anchorLines[ax], gridPoint.Point)
			gridPoint.offsets[ax] = distance / f.dist
		}
	}

	gridPoint.PackedCoords = store.PackCoordinates(
		gridLine0.Axis, gridLine1.Axis, gridLine0.Offset, gridLine1.Offset,
	)

	return gridPoint
}

func (f *Field) NearestNeighbor(currentPoint GridPoint, prevLine, currentLine GridLine, positiveSide bool) (GridPoint, GridLine) {
	var nextLineResult GridLine
	var nextPointResult GridPoint
	currentDistance := 1000000.0
	prevLineLine := f.GetLine(prevLine)

	for _, axis := range axisIndices {
		if axis != prevLine.Axis && axis != currentLine.Axis {
			axisOffset := currentPoint.offsets[axis]
			for _, offset := range [2]float64{math.Ceil(axisOffset), math.Floor(axisOffset)} {
				nextLine := GridLine{axis, int16(offset)}
				// TODO: look how many times MakeGridPoint gets called unnecessarily here
				nextPoint := f.MakeGridPoint(currentLine, nextLine)
				distance := geom.Distance(prevLineLine, nextPoint.Point)
				//axisNames := [5]string{"A", "B", "C", "D", "E"}
				//fmt.Printf("Neighbor %s %d ; %s %d ", axisNames[currentLine.Axis], currentLine.Offset, axisNames[nextLine.Axis], nextLine.Offset)
				//fmt.Printf("distance=%.1f\n", distance)
				if (positiveSide && distance > 0) || (!positiveSide && distance < 0) {
					if math.Abs(distance) < currentDistance {
						currentDistance = math.Abs(distance)
						nextLineResult = nextLine
						nextPointResult = nextPoint
					}
				}
			}
		}
	}

	//fmt.Printf("currentDistance=%.1f, positiveSide=%t\n", currentDistance, positiveSide)
	return nextPointResult, nextLineResult
}

var axesRotation = [5][5]bool{
	{true, true, true, false, false},
	{false, true, true, true, false},
	{false, false, true, true, true},
	{true, false, false, true, true},
	{true, true, false, false, true},
}

func (f *Field) NextPoint(prevPoint, currPoint GridPoint, prevLine, currLine GridLine, isRightTurn bool) (GridPoint, GridPoint, GridLine, GridLine) {
	axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
	prevPointSign := geom.Distance(f.GetLine(currLine), prevPoint.Point) < 0
	positiveSide := isRightTurn != axisRotation != prevPointSign

	nextPoint, nextLine := f.NearestNeighbor(currPoint, prevLine, currLine, positiveSide)
	return currPoint, nextPoint, currLine, nextLine
}
