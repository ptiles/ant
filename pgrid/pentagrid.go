package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"image"
	"math"
	"os"
)

const X = 0
const Y = 1

const GRID_LINES_TOTAL = uint8(5)

type Field struct {
	Rules        []bool
	Limit        uint8
	InitialPoint string
	verbose      bool

	anchors          [GRID_LINES_TOTAL]Point
	anchorLines      [GRID_LINES_TOTAL]Line
	axisUnits        [GRID_LINES_TOTAL]Point
	normals          [GRID_LINES_TOTAL]Point
	intersectAnchors [GRID_LINES_TOTAL][GRID_LINES_TOTAL]Point
	intersectVectors [GRID_LINES_TOTAL][GRID_LINES_TOTAL]Point
}

func New(r float64, rules []bool, initialPoint string, verbose bool) *Field {
	phi := utils.FromDegrees(360 / int(GRID_LINES_TOTAL))
	axisAngle0 := utils.FromDegrees(90) + phi/2

	result := &Field{
		Rules: rules, Limit: uint8(len(rules)),
		InitialPoint: initialPoint, verbose: verbose,
	}

	for ax := range GRID_LINES_TOTAL {
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

	for ax0 := range GRID_LINES_TOTAL {
		for ax1 := range GRID_LINES_TOTAL {
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

var axisNames = [GRID_LINES_TOTAL]string{"A", "B", "C", "D", "E"}

type GridLine struct {
	Axis   uint8
	Offset offsetInt
}

func (gl *GridLine) String() string {
	return fmt.Sprintf("%s%d", axisNames[gl.Axis], gl.Offset)
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

type offsetInt int16

type GridAxes struct {
	Axis0   uint8
	Axis1   uint8
	Offset0 offsetInt
	Offset1 offsetInt
}

type GridOffsets [GRID_LINES_TOTAL]offsetInt

func (gp *GridPoint) String() string {
	offsets := gp.Offsets
	ax0, ax1 := gp.Axes.Axis0, gp.Axes.Axis1
	return fmt.Sprintf(
		"%s%d:%s%d [A:%d, B:%d, C:%d, D:%d, E:%d]",
		axisNames[ax0], offsets[ax0], axisNames[ax1], offsets[ax1],
		offsets[0], offsets[1], offsets[2], offsets[3], offsets[4],
	)
}
func (gp *GridPoint) Print() {
	fmt.Println(gp)
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

	for ax := range GRID_LINES_TOTAL {
		if ax == gridLine0.Axis || ax == gridLine1.Axis {
			continue
		}
		dist := distance(f.anchorLines[ax], point)
		gridPoint.Offsets[ax] = offsetInt(math.Ceil(dist))
	}

	return gridPoint
}

var deBruijnX = [GRID_LINES_TOTAL]float64{}
var deBruijnY = [GRID_LINES_TOTAL]float64{}

func init() {
	if GRID_LINES_TOTAL < 5 || GRID_LINES_TOTAL%2 == 0 {
		fmt.Println("GRID_LINES_TOTAL should odd and at least 5")
		os.Exit(1)
	}

	floatLines := float64(GRID_LINES_TOTAL)
	for i := range GRID_LINES_TOTAL {
		floatI := float64(i)
		deBruijnX[i] = math.Sin(2 * math.Pi * floatI / floatLines)
		deBruijnY[i] = math.Cos(2 * math.Pi * floatI / floatLines)
	}
}

const deBruijnScale = 2

func (gp *GridPoint) getCenterPoint() image.Point {
	x := 0.5*deBruijnX[gp.Axes.Axis0] + 0.5*deBruijnX[gp.Axes.Axis1]
	y := 0.5*deBruijnY[gp.Axes.Axis0] + 0.5*deBruijnY[gp.Axes.Axis1]

	for i := range GRID_LINES_TOTAL {
		floatOffset := float64(gp.Offsets[i])
		x += floatOffset * deBruijnX[i]
		y += floatOffset * deBruijnY[i]
	}

	return image.Point{X: int(x * deBruijnScale), Y: int(y * deBruijnScale)}
}

func (f *Field) nearestNeighbor(currentPointOffsets GridOffsets, prevLine, currentLine GridLine, positiveSide bool) (GridPoint, GridLine) {
	var nextLine GridLine
	var nextPointPoint Point
	currentDistance := 1000000.0
	prevLineLine := f.getLine(prevLine)

	for axis := range GRID_LINES_TOTAL {
		if axis == prevLine.Axis || axis == currentLine.Axis {
			continue
		}
		axisOffset := currentPointOffsets[axis]
		for _, offset := range [2]offsetInt{axisOffset, axisOffset - 1} {
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
