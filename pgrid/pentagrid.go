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

	anchors          [GridLinesTotal]Point
	anchorLines      [GridLinesTotal]Line
	axisUnits        [GridLinesTotal]Point
	normals          [GridLinesTotal]Point
	intersectAnchors [GridLinesTotal][GridLinesTotal]Point
	intersectVectors [GridLinesTotal][GridLinesTotal]Point
}

func New(r float64, rules []bool, initialPoint string) *Field {
	phi := utils.FromDegrees(360 / int(GridLinesTotal))
	axisAngle0 := utils.FromDegrees(90) + phi/2

	result := &Field{
		Rules: rules, Limit: uint8(len(rules)), InitialPoint: initialPoint,
	}

	for ax := range GridLinesTotal {
		phiAx := phi * float64(ax)

		result.anchors[ax][X] = r * math.Cos(phiAx)
		result.anchors[ax][Y] = r * math.Sin(phiAx)

		result.axisUnits[ax][X] = math.Cos(axisAngle0 + phiAx)
		result.axisUnits[ax][Y] = math.Sin(axisAngle0 + phiAx)

		result.normals[ax][X] = math.Cos(phi/2 + phiAx)
		result.normals[ax][Y] = math.Sin(phi/2 + phiAx)

		anchorEnd := Point{
			result.anchors[ax][X] + result.axisUnits[ax][X],
			result.anchors[ax][Y] + result.axisUnits[ax][Y],
		}
		result.anchorLines[ax] = Line{result.anchors[ax], anchorEnd}
	}

	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
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

var AxisNames = [25]string{
	"A", "B", "C", "D", "E",
	"F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y",
}

type GridLine struct {
	Axis   uint8
	Offset offsetInt
}

func (gl *GridLine) String() string {
	return fmt.Sprintf("%s%d", AxisNames[gl.Axis], gl.Offset)
}
func (gl *GridLine) Print() {
	fmt.Println(gl)
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

type offsetInt int32

type GridAxes struct {
	Axis0   uint8
	Axis1   uint8
	Offset0 offsetInt
	Offset1 offsetInt
}

type GridOffsets [GridLinesTotal]offsetInt

func (gp *GridPoint) String() string {
	offsets := gp.Offsets
	ax0, ax1 := gp.Axes.Axis0, gp.Axes.Axis1
	return fmt.Sprintf(
		"[A:%d, B:%d, C:%d, D:%d, E:%d] %s%d:%s%d",
		offsets[0], offsets[1], offsets[2], offsets[3], offsets[4],
		AxisNames[ax0], offsets[ax0], AxisNames[ax1], offsets[ax1],
	)
}

func (ga *GridAxes) String() string {
	return fmt.Sprintf(
		"%s%d:%s%d",
		AxisNames[ga.Axis0], ga.Offset0, AxisNames[ga.Axis1], ga.Offset1,
	)
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

	for ax := range GridLinesTotal {
		if ax == gridLine0.Axis || ax == gridLine1.Axis {
			continue
		}
		dist := distance(f.anchorLines[ax], point)
		gridPoint.Offsets[ax] = offsetInt(math.Ceil(dist))
	}

	return gridPoint
}

var deBruijnX = [GridLinesTotal]float64{}
var deBruijnY = [GridLinesTotal]float64{}

func init() {
	if GridLinesTotal%2 == 0 || GridLinesTotal < 5 || GridLinesTotal > 25 {
		fmt.Println("GridLinesTotal should be odd number between 5 and 25")
	}

	floatLines := float64(GridLinesTotal)
	for i := range GridLinesTotal {
		floatI := float64(i)
		deBruijnX[i] = math.Sin(2 * math.Pi * floatI / floatLines)
		deBruijnY[i] = math.Cos(2 * math.Pi * floatI / floatLines)
	}
}

func (gp *GridPoint) getCenterPoint() image.Point {
	x := 0.5*deBruijnX[gp.Axes.Axis0] + 0.5*deBruijnX[gp.Axes.Axis1]
	y := 0.5*deBruijnY[gp.Axes.Axis0] + 0.5*deBruijnY[gp.Axes.Axis1]

	for i := range GridLinesTotal {
		floatOffset := float64(gp.Offsets[i])
		x += floatOffset * deBruijnX[i]
		y += floatOffset * deBruijnY[i]
	}

	return image.Point{X: int(x * deBruijnScale), Y: int(y * deBruijnScale)}
}

func (gp *GridPoint) getCornerPoints() [4]image.Point {
	x := float64(0)
	y := float64(0)

	for i := range GridLinesTotal {
		floatOffset := float64(gp.Offsets[i])
		x += floatOffset * deBruijnX[i]
		y += floatOffset * deBruijnY[i]
	}

	// TODO: prepare this in init() and store in counter-clockwise order
	_0 := 0.05
	_1 := 0.95
	dax0x := deBruijnX[gp.Axes.Axis0]
	dax0y := deBruijnY[gp.Axes.Axis0]
	dax1x := deBruijnX[gp.Axes.Axis1]
	dax1y := deBruijnY[gp.Axes.Axis1]

	return [4]image.Point{
		//{X: int((x+_0*dax0x+_0*dax1x)*deBruijnScale) - 1, Y: int((y + _0*dax0y + _0*dax1y) * deBruijnScale)},
		//{X: int((x + _0*dax0x + _0*dax1x) * deBruijnScale), Y: int((y+_0*dax0y+_0*dax1y)*deBruijnScale) - 1},
		//{X: int((x+_0*dax0x+_0*dax1x)*deBruijnScale) + 1, Y: int((y + _0*dax0y + _0*dax1y) * deBruijnScale)},
		//{X: int((x + _0*dax0x + _0*dax1x) * deBruijnScale), Y: int((y+_0*dax0y+_0*dax1y)*deBruijnScale) + 1},

		{X: int((x + _0*dax0x + _0*dax1x) * deBruijnScale), Y: int((y + _0*dax0y + _0*dax1y) * deBruijnScale)},
		{X: int((x + _0*dax0x + _1*dax1x) * deBruijnScale), Y: int((y + _0*dax0y + _1*dax1y) * deBruijnScale)},
		{X: int((x + _1*dax0x + _1*dax1x) * deBruijnScale), Y: int((y + _1*dax0y + _1*dax1y) * deBruijnScale)},
		{X: int((x + _1*dax0x + _0*dax1x) * deBruijnScale), Y: int((y + _1*dax0y + _0*dax1y) * deBruijnScale)},
	}
}

func (f *Field) nearestNeighbor(currentPointOffsets GridOffsets, prevLine, currentLine GridLine, positiveSide bool) (GridPoint, GridLine) {
	var nextLine GridLine
	var nextPointPoint Point
	currentDistance := 1000000.0
	prevLineLine := f.getLine(prevLine)

	for axis := range GridLinesTotal {
		if axis == prevLine.Axis || axis == currentLine.Axis {
			continue
		}
		axisOffset := currentPointOffsets[axis]
		for _, offset := range [2]offsetInt{axisOffset, axisOffset - 1} {
			line := GridLine{axis, offset}
			pointPoint := f.gridPointToPoint(currentLine, line)
			dist := distance(prevLineLine, pointPoint)
			//fmt.Printf("Neighbor %s %d ; %s %d ", AxisNames[currentLine.Axis], currentLine.Offset, AxisNames[line.Axis], line.Offset)
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
