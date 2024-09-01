package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"image"
	"math"
)

const X = 0
const Y = 1

type Field struct {
	Rules        []bool
	Limit        uint8
	InitialPoint string
	verbose      bool

	anchors          [5]Point
	anchorLines      [5]Line
	axisUnits        [5]Point
	normals          [5]Point
	intersectAnchors [5][5]Point
	intersectVectors [5][5]Point
}

func New(r float64, rules []bool, initialPoint string, verbose bool) *Field {
	phi := utils.FromDegrees(72)
	axisAngle0 := utils.FromDegrees(180 - 54)

	result := &Field{
		Rules: rules, Limit: uint8(len(rules)),
		InitialPoint: initialPoint, verbose: verbose,
	}

	for ax := range 5 {
		phiAx := phi * float64(ax)

		result.anchors[ax][X] = r * math.Cos(phiAx)
		result.anchors[ax][Y] = r * math.Sin(phiAx)

		result.axisUnits[ax][X] = 1 * math.Cos(axisAngle0+phiAx)
		result.axisUnits[ax][Y] = 1 * math.Sin(axisAngle0+phiAx)

		result.normals[ax][X] = 1 * math.Cos(0.5*phi+phiAx)
		result.normals[ax][Y] = 1 * math.Sin(0.5*phi+phiAx)

		anchorEnd := Point{
			result.anchors[ax][X] + result.axisUnits[ax][X],
			result.anchors[ax][Y] + result.axisUnits[ax][Y],
		}
		result.anchorLines[ax] = Line{result.anchors[ax], anchorEnd}
	}

	for ax0 := range uint8(5) {
		for ax1 := range uint8(5) {
			line0 := result.getLine(GridLine{ax0, 0})
			line1 := result.getLine(GridLine{ax1, 0})
			intersectAnchor := intersection(line0, line1)
			result.intersectAnchors[ax0][ax1] = intersectAnchor

			line01 := result.getLine(GridLine{ax0, 1})
			intersection01 := intersection(line01, line1)
			result.intersectVectors[ax0][ax1] = Point{
				intersection01[X] - intersectAnchor[X],
				intersection01[Y] - intersectAnchor[Y],
			}
		}
	}

	return result
}

var axisNames = [5]string{"A", "B", "C", "D", "E"}

type GridLine struct {
	Axis   uint8
	Offset int16
}

func (gl *GridLine) Sprint() string {
	return fmt.Sprintf("%s%d", axisNames[gl.Axis], gl.Offset)
}
func (gl *GridLine) Print() {
	fmt.Println(gl.Sprint())
}

func (f *Field) getLine(gl GridLine) Line {
	axis := f.axisUnits[gl.Axis]
	anchor := f.anchors[gl.Axis]
	normal := f.normals[gl.Axis]
	offset := float64(gl.Offset)

	point1 := Point{anchor[X] + normal[X]*offset, anchor[Y] + normal[Y]*offset}
	point2 := Point{point1[X] + axis[X], point1[Y] + axis[Y]}
	return Line{point1, point2}
}

type GridPoint struct {
	Axes    GridAxes
	Offsets GridOffsets
	Point   Point
}

type GridAxes struct {
	Axis0   uint8
	Axis1   uint8
	Offset0 int16
	Offset1 int16
}

type GridOffsets [5]int16

func (gp *GridPoint) Sprint() string {
	offsets := gp.Offsets
	ax0, ax1 := gp.Axes.Axis0, gp.Axes.Axis1
	return fmt.Sprintf(
		"%s%d:%s%d [A:%d, B:%d, C:%d, D:%d, E:%d]",
		axisNames[ax0], offsets[ax0], axisNames[ax1], offsets[ax1],
		offsets[0], offsets[1], offsets[2], offsets[3], offsets[4],
	)
}
func (gp *GridPoint) Print() {
	fmt.Println(gp.Sprint())
}

func (f *Field) gridPointToPoint(gridLine0, gridLine1 GridLine) Point {
	anchor := f.intersectAnchors[gridLine0.Axis][gridLine1.Axis]

	vector0 := f.intersectVectors[gridLine0.Axis][gridLine1.Axis]
	vector1 := f.intersectVectors[gridLine1.Axis][gridLine0.Axis]

	off0, off1 := float64(gridLine0.Offset), float64(gridLine1.Offset)

	return Point{
		anchor[X] + vector0[X]*off0 + vector1[X]*off1,
		anchor[Y] + vector0[Y]*off0 + vector1[Y]*off1,
	}
}

func (f *Field) makeGridPoint(gridLine0, gridLine1 GridLine, point Point) GridPoint {
	if gridLine0.Axis > gridLine1.Axis {
		gridLine0, gridLine1 = gridLine1, gridLine0
	}

	gridPoint := GridPoint{
		Axes: GridAxes{
			gridLine0.Axis, gridLine1.Axis, gridLine0.Offset, gridLine1.Offset,
		},
		Point: point,
	}

	gridPoint.Offsets[gridLine0.Axis] = gridLine0.Offset
	gridPoint.Offsets[gridLine1.Axis] = gridLine1.Offset

	for ax := range uint8(5) {
		if ax == gridLine0.Axis || ax == gridLine1.Axis {
			continue
		}
		dist := distance(f.anchorLines[ax], gridPoint.Point)
		gridPoint.Offsets[ax] = int16(math.Ceil(dist))
	}

	return gridPoint
}

var deBruijnConstants = [5][2]float64{
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

func deBruijn(floatOffsets *[5]float64) (float64, float64) {
	x := 0 +
		floatOffsets[0]*deBruijnConstants[0][0] +
		floatOffsets[1]*deBruijnConstants[1][0] +
		floatOffsets[2]*deBruijnConstants[2][0] +
		floatOffsets[3]*deBruijnConstants[3][0] +
		floatOffsets[4]*deBruijnConstants[4][0]

	y := 0 +
		floatOffsets[0]*deBruijnConstants[0][1] +
		floatOffsets[1]*deBruijnConstants[1][1] +
		floatOffsets[2]*deBruijnConstants[2][1] +
		floatOffsets[3]*deBruijnConstants[3][1] +
		floatOffsets[4]*deBruijnConstants[4][1]

	// Flip X and Y axes to make stars vertically symmetrical
	return y, x
}

func (gp *GridPoint) getCenterPoint() image.Point {
	var floatOffsets [5]float64
	for i := range 5 {
		floatOffsets[i] = float64(gp.Offsets[i])
	}

	floatOffsets[gp.Axes.Axis0] += 0.5
	floatOffsets[gp.Axes.Axis1] += 0.5

	x, y := deBruijn(&floatOffsets)
	return image.Point{X: int(x * 2), Y: int(y * 2)}
	//return geom.Point{x, y}
}

func (f *Field) nearestNeighbor(currentPoint GridPoint, prevLine, currentLine GridLine, positiveSide bool) (GridPoint, GridLine) {
	var nextLine GridLine
	var nextPointPoint Point
	currentDistance := 1000000.0
	prevLineLine := f.getLine(prevLine)

	for axis := range uint8(5) {
		if axis == prevLine.Axis || axis == currentLine.Axis {
			continue
		}
		axisOffset := currentPoint.Offsets[axis]
		for _, offset := range [2]int16{axisOffset, axisOffset - 1} {
			line := GridLine{axis, offset}
			pointPoint := f.gridPointToPoint(currentLine, line)
			dist := distance(prevLineLine, pointPoint)
			//fmt.Printf("Neighbor %s %d ; %s %d ", axisNames[currentLine.Axis], currentLine.Offset, axisNames[line.Axis], line.Offset)
			//fmt.Printf("distance=%.1f\n", dist)
			if (positiveSide && dist > 0) || (!positiveSide && dist < 0) {
				if math.Abs(dist) < currentDistance {
					currentDistance = math.Abs(dist)
					nextLine = line
					nextPointPoint = pointPoint
				}
			}
		}
	}

	nextPoint := f.makeGridPoint(currentLine, nextLine, nextPointPoint)
	//fmt.Printf("currentDistance=%.1f, positiveSide=%t\n", currentDistance, positiveSide)
	return nextPoint, nextLine
}
