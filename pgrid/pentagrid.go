package pgrid

import (
	"github.com/ptiles/ant/canv"
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
	dist       float64
	anchors    [5]geom.Point
	Axes       [5]geom.Point
	normals    [5]geom.Point
	canvasFile *canv.Canvas
}

func New(r, dist float64, phi0degrees int, canvasFile *canv.Canvas) Field {
	phi0 := fromDegrees(phi0degrees)
	phi := fromDegrees(72)
	axis0 := fromDegrees(180 - 54 + phi0degrees)

	var anchors [5]geom.Point
	var axes [5]geom.Point
	var normals [5]geom.Point

	for ax := 0; ax < len(anchors); ax++ {
		phiAx := phi * float64(ax)

		anchors[ax][X] = r * math.Cos(phi0+phiAx)
		anchors[ax][Y] = r * math.Sin(phi0+phiAx)

		axes[ax][X] = 1 * math.Cos(axis0+phiAx)
		axes[ax][Y] = 1 * math.Sin(axis0+phiAx)

		normals[ax][X] = dist * math.Cos(phi0+0.5*phi+phiAx)
		normals[ax][Y] = dist * math.Sin(phi0+0.5*phi+phiAx)
	}

	return Field{dist, anchors, axes, normals, canvasFile}
}

type GridLine struct {
	Axis   uint8
	Offset int16
	Line   geom.Line
}

func (f Field) MakeGridLine(ax uint8, off int16) GridLine {
	axis := f.Axes[ax]
	anchor := f.anchors[ax]
	normal := f.normals[ax]
	distance := float64(off)

	point1 := geom.Point{anchor[X] + normal[X]*distance, anchor[Y] + normal[Y]*distance}
	point2 := geom.Point{point1[X] + axis[X], point1[Y] + axis[Y]}
	line := geom.Line{point1, point2}

	return GridLine{ax, off, line}
}

type GridPoint struct {
	axes         [5]bool
	offsets      [5]float64
	Point        geom.Point
	name         string
	PackedCoords store.PackedCoordinates
}

func (f Field) MakeGridPoint(gridLine0, gridLine1 GridLine, name string) GridPoint {
	gridPoint := GridPoint{name: name}
	gridPoint.offsets[gridLine0.Axis] = float64(gridLine0.Offset)
	gridPoint.axes[gridLine0.Axis] = true
	gridPoint.offsets[gridLine1.Axis] = float64(gridLine1.Offset)
	gridPoint.axes[gridLine1.Axis] = true

	gridPoint.Point = geom.Intersection(gridLine0.Line, gridLine1.Line)

	for ax := range gridPoint.offsets {
		if !gridPoint.axes[ax] {
			anchorEnd := geom.Point{f.anchors[ax][X] + f.Axes[ax][X], f.anchors[ax][Y] + f.Axes[ax][Y]}
			anchorLine := geom.Line{f.anchors[ax], anchorEnd}
			distance := geom.Distance(anchorLine, gridPoint.Point)
			gridPoint.offsets[ax] = distance / f.dist
		}
	}

	gridPoint.PackedCoords = store.PackCoordinates(
		gridLine0.Axis, gridLine1.Axis, gridLine0.Offset, gridLine1.Offset,
	)

	return gridPoint
}

type Neighbor struct {
	nextLine  GridLine
	nextPoint GridPoint
	distance  float64
}

func (f Field) nearestNeighbors(currPoint GridPoint, prevLine, currLine GridLine, positiveSide bool) []Neighbor {
	// TODO: use plain array for this
	neighbors := make([]Neighbor, 0, 12)
	for _, axis := range axisIndices {
		if axis != prevLine.Axis && axis != currLine.Axis {
			axisOffset := currPoint.offsets[axis]
			for _, offset := range [...]float64{math.Ceil(axisOffset) + 1, math.Floor(axisOffset) - 1, math.Ceil(axisOffset), math.Floor(axisOffset)} {
				nextLine := f.MakeGridLine(axis, int16(offset))
				nextPoint := f.MakeGridPoint(currLine, nextLine, "")
				distance := geom.Distance(prevLine.Line, nextPoint.Point)
				//axisNames := [5]string{"A", "B", "C", "D", "E"}
				//fmt.Printf("Neighbor %s %d ; %s %d ", axisNames[currLine.Axis], currLine.Offset, axisNames[nextLine.Axis], nextLine.Offset)
				//fmt.Printf("distance=%.1f\n", distance)
				if (positiveSide && distance > 0) || (!positiveSide && distance < 0) {
					neighbors = append(neighbors, Neighbor{nextLine, nextPoint, distance})
				}
			}
		}
	}
	return neighbors
}

func (f Field) NearestNeighbor(currentPoint GridPoint, prevLine, currentLine GridLine, positiveSide bool) (GridPoint, GridLine) {
	var nextLine GridLine
	var nextPoint GridPoint
	currentDistance := 1000000.0
	for _, neighbor := range f.nearestNeighbors(currentPoint, prevLine, currentLine, positiveSide) {
		distance := neighbor.distance
		if math.Abs(distance) < currentDistance {
			currentDistance = math.Abs(distance)
			nextLine = neighbor.nextLine
			nextPoint = neighbor.nextPoint
		}
	}
	//fmt.Printf("currentDistance=%.1f, positiveSide=%t\n", currentDistance, positiveSide)
	return nextPoint, nextLine
}

var axesRotation = [5][5]bool{
	{true, true, true, false, false},
	{false, true, true, true, false},
	{false, false, true, true, true},
	{true, false, false, true, true},
	{true, true, false, false, true},
}

func (f Field) NextPoint(prevPoint, currPoint GridPoint, prevLine, currLine GridLine, isRightTurn bool) (GridPoint, GridPoint, GridLine, GridLine) {
	axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
	prevPointSign := geom.Distance(currLine.Line, prevPoint.Point) < 0
	positiveSide := isRightTurn != axisRotation != prevPointSign

	nextPoint, nextLine := f.NearestNeighbor(currPoint, prevLine, currLine, positiveSide)
	return currPoint, nextPoint, currLine, nextLine
}
