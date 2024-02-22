package pgrid

import (
	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/store"
	"github.com/ptiles/ant/utils"
	"math"
)

const X = 0
const Y = 1

const A = 0
const B = 1
const C = 2
const D = 3
const E = 4

type Field struct {
	dist             float64
	anchors          [5]geom.Point
	anchorLines      [5]geom.Line
	axisUnits        [5]geom.Point
	normals          [5]geom.Point
	intersectAnchors [5][5]geom.Point
	intersectVectors [5][5]geom.Point
}

func New(r, dist float64) Field {
	phi := utils.FromDegrees(72)
	axisAngle0 := utils.FromDegrees(180 - 54)

	result := Field{dist: dist}

	for ax := 0; ax < 5; ax++ {
		phiAx := phi * float64(ax)

		result.anchors[ax][X] = r * math.Cos(phiAx)
		result.anchors[ax][Y] = r * math.Sin(phiAx)

		result.axisUnits[ax][X] = 1 * math.Cos(axisAngle0+phiAx)
		result.axisUnits[ax][Y] = 1 * math.Sin(axisAngle0+phiAx)

		result.normals[ax][X] = dist * math.Cos(0.5*phi+phiAx)
		result.normals[ax][Y] = dist * math.Sin(0.5*phi+phiAx)

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
	Offsets      [5]int16
	Point        geom.Point
	PackedCoords store.PackedCoordinates
}

func (f *Field) MakeGridPoint(gridLine0, gridLine1 GridLine) GridPoint {
	gridPoint := GridPoint{}
	offset0 := gridLine0.Offset
	gridPoint.Offsets[gridLine0.Axis] = offset0
	gridPoint.axes[gridLine0.Axis] = true
	offset1 := gridLine1.Offset
	gridPoint.Offsets[gridLine1.Axis] = offset1
	gridPoint.axes[gridLine1.Axis] = true

	anchor := f.intersectAnchors[gridLine0.Axis][gridLine1.Axis]

	vector0 := f.intersectVectors[gridLine0.Axis][gridLine1.Axis]
	vector1 := f.intersectVectors[gridLine1.Axis][gridLine0.Axis]
	off0, off1 := float64(offset0), float64(offset1)
	gridPoint.Point = geom.Point{
		anchor[X] + vector0[X]*off0 + vector1[X]*off1,
		anchor[Y] + vector0[Y]*off0 + vector1[Y]*off1,
	}

	for ax := range gridPoint.Offsets {
		if !gridPoint.axes[ax] {
			distance := geom.Distance(f.anchorLines[ax], gridPoint.Point)
			gridPoint.Offsets[ax] = int16(math.Ceil(distance / f.dist))
		}
	}

	gridPoint.PackedCoords = store.PackCoordinates(
		gridLine0.Axis, gridLine1.Axis, gridLine0.Offset, gridLine1.Offset,
	)

	return gridPoint
}

var DeBrujinConstants = [5][2]float64{
	{
		math.Cos(2 * math.Pi * float64(0) / 5),
		math.Sin(2 * math.Pi * float64(0) / 5),
	},
	{
		math.Cos(2 * math.Pi * float64(1) / 5),
		math.Sin(2 * math.Pi * float64(1) / 5),
	},
	{
		math.Cos(2 * math.Pi * float64(2) / 5),
		math.Sin(2 * math.Pi * float64(2) / 5),
	},
	{
		math.Cos(2 * math.Pi * float64(3) / 5),
		math.Sin(2 * math.Pi * float64(3) / 5),
	},
	{
		math.Cos(2 * math.Pi * float64(4) / 5),
		math.Sin(2 * math.Pi * float64(4) / 5),
	},
}

func DeBrujin(floatOffsets *[5]float64) (float64, float64) {
	x := 0 +
		floatOffsets[0]*DeBrujinConstants[0][0] +
		floatOffsets[1]*DeBrujinConstants[1][0] +
		floatOffsets[2]*DeBrujinConstants[2][0] +
		floatOffsets[3]*DeBrujinConstants[3][0] +
		floatOffsets[4]*DeBrujinConstants[4][0]

	y := 0 +
		floatOffsets[0]*DeBrujinConstants[0][1] +
		floatOffsets[1]*DeBrujinConstants[1][1] +
		floatOffsets[2]*DeBrujinConstants[2][1] +
		floatOffsets[3]*DeBrujinConstants[3][1] +
		floatOffsets[4]*DeBrujinConstants[4][1]

	// Flip X and Y axes to make stars vertically symmetrical
	return y, x
}

func (f *Field) GetCenterPoint(gp *GridPoint) geom.Point {
	var floatOffsets [5]float64
	for i := 0; i < 5; i++ {
		floatOffsets[i] = float64(gp.Offsets[i])
	}

	axis0, axis1 := store.UnpackAxes(gp.PackedCoords.PackedAxes)
	floatOffsets[axis0] += 0.5
	floatOffsets[axis1] += 0.5

	x, y := DeBrujin(&floatOffsets)
	return geom.Point{x * 3, y * 3}
}

func (f *Field) NearestNeighbor(currentPoint GridPoint, prevLine, currentLine GridLine, positiveSide bool) (GridPoint, GridLine) {
	var nextLineResult GridLine
	var nextPointResult GridPoint
	currentDistance := 1000000.0
	prevLineLine := f.GetLine(prevLine)

	for axis := uint8(0); axis < 5; axis++ {
		if axis != prevLine.Axis && axis != currentLine.Axis {
			axisOffset := currentPoint.Offsets[axis]
			for _, offset := range [2]int16{axisOffset, axisOffset - 1} {
				nextLine := GridLine{axis, offset}
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
