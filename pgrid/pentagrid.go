package pgrid

import (
	"github.com/ptiles/ant/canv"
	"github.com/ptiles/ant/geom"
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

const AxisCount = uint8(5)

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

// _   AB   AC   AD   AE |   BC    BD    BE |   CD    CE |   DE
// _ a[0]                | a[1]             | a[2]       | else
// _ a[1] a[2] a[3] a[4] | a[2]  a[3]  a[4] | a[3]  a[4] | a[4]
// _    0   +1   +2   +3 |    4    +1    +2 |    7    +1 |    9
// _    0    1    2    3 |    4     5     6 |    7     8 |    9
//func axesOffset(a [5]bool) int {
//	if a[0] {
//		return qi(a[2], 1, 0) + qi(a[3], 2, 0) + qi(a[4], 4, 0)
//	}
//	if a[1] {
//		return 4 + qi(a[3], 1, 0) + qi(a[4], 2, 0)
//	}
//	if a[2] {
//		return 7 + qi(a[4], 1, 0)
//	}
//	return 9
//}

func aIndex(a0, a1 uint8) uint8 {
	if a0 == A {
		return a1 - 1
	}
	if a0 == B {
		return a1 + 2
	}
	if a0 == C {
		return a1 + 4
	}
	return 9
}

var values = [1 << 8][1 << 8][10]uint8{}

var a0ByIndex = [10]uint8{A, A, A, A, B, B, B, C, C, D}
var a1ByIndex = [10]uint8{B, C, D, E, C, D, E, D, E, E}

func (f Field) GetPointsByOffsets(off0, off1 int16) ([]geom.Point, []uint8) {
	points := make([]geom.Point, 0, 10)
	colors := make([]uint8, 0, 10)
	for index := uint8(0); index < 10; index++ {
		color := values[uint8(off0)][uint8(off1)][index]
		if color > 0 {
			line0 := f.MakeGridLine(a0ByIndex[index], off0).Line
			line1 := f.MakeGridLine(a1ByIndex[index], off1).Line
			point := geom.Intersection(line0, line1)
			points = append(points, point)
			colors = append(colors, color)
		}
	}
	return points, colors
}

var MinOffset0 = int16(0)
var MaxOffset0 = int16(0)

var MinOffset1 = int16(0)
var MaxOffset1 = int16(0)

func (gp GridPoint) Get() uint8 {
	return values[uint8(gp.off0)][uint8(gp.off1)][gp.aIndex]
}

func (gp GridPoint) Set(value uint8) {
	off0, off1 := gp.off0, gp.off1

	if off0 > MaxOffset0 {
		MaxOffset0 = off0
	}
	if off0 < MinOffset0 {
		MinOffset0 = off0
	}

	if off1 > MaxOffset1 {
		MaxOffset1 = off1
	}
	if off1 < MinOffset0 {
		MinOffset1 = off1
	}

	values[uint8(off0)][uint8(off1)][gp.aIndex] = value
}

type GridPoint struct {
	axes    [5]bool
	offsets [5]float64
	Point   geom.Point
	name    string
	off0    int16
	off1    int16
	aIndex  uint8
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

	if gridLine0.Axis < gridLine1.Axis {
		gridPoint.off0, gridPoint.off1 = gridLine0.Offset, gridLine1.Offset
		gridPoint.aIndex = aIndex(gridLine0.Axis, gridLine1.Axis)
	} else {
		gridPoint.off0, gridPoint.off1 = gridLine1.Offset, gridLine0.Offset
		gridPoint.aIndex = aIndex(gridLine1.Axis, gridLine0.Axis)
	}

	return gridPoint
}

func (f Field) DrawGridPoint(gridPoint GridPoint, prefix string) {
	f.canvasFile.DrawPoint(gridPoint.Point, prefix+gridPoint.name, 0)
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
