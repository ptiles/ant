package pgrid

import (
	"fmt"
	"image"
	"math"
)

type Field struct {
	Rules        []bool
	Limit        uint8
	InitialPoint string

	offsetsToLast  allOffsetDeltas
	offsetsToFirst allOffsetDeltas
}

func New(r float64, rules []bool, initialPoint string) *Field {
	gg := newGridGeometry(r)

	return &Field{
		Rules:        rules,
		Limit:        uint8(len(rules)),
		InitialPoint: initialPoint,

		offsetsToLast:  gg.newOffsetsToLast(),
		offsetsToFirst: gg.newOffsetsToFirst(),
	}
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
		AxisNames[ga.Axis0], ga.Coords.Offset0, AxisNames[ga.Axis1], ga.Coords.Offset1,
	)
}

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
		// floatOffset := float64(gp.Offsets[i]) * inflation
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
		// floatOffset := float64(gp.Offsets[i]) * inflation
		floatOffset := float64(gp.Offsets[i])
		x += floatOffset * deBruijnX[i]
		y += floatOffset * deBruijnY[i]
	}

	// TODO: prepare this in init() and store in counter-clockwise order
	dax0x := deBruijnX[gp.Axes.Axis0]
	dax0y := deBruijnY[gp.Axes.Axis0]
	dax1x := deBruijnX[gp.Axes.Axis1]
	dax1y := deBruijnY[gp.Axes.Axis1]

	return [4]image.Point{
		{X: int((x + _0*dax0x + _0*dax1x) * deBruijnScale), Y: int((y + _0*dax0y + _0*dax1y) * deBruijnScale)},
		{X: int((x + _0*dax0x + _1*dax1x) * deBruijnScale), Y: int((y + _0*dax0y + _1*dax1y) * deBruijnScale)},
		{X: int((x + _1*dax0x + _1*dax1x) * deBruijnScale), Y: int((y + _1*dax0y + _1*dax1y) * deBruijnScale)},
		{X: int((x + _1*dax0x + _0*dax1x) * deBruijnScale), Y: int((y + _1*dax0y + _0*dax1y) * deBruijnScale)},
	}
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
