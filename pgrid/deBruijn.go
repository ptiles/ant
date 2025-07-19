package pgrid

import (
	"image"
	"math"
)

//const DrawTilesAndPoints = true
//const deBruijnScale = 48
//const padding = deBruijnScale

const DrawTilesAndPoints = false
const deBruijnScale = 2
const Padding = deBruijnScale * 4
const LineScale = deBruijnScale * 2.5

//const inflation = 1
//const inflation = 2

var deBruijnX = [GridLinesTotal]float64{}
var deBruijnY = [GridLinesTotal]float64{}

func init() {
	floatLines := float64(GridLinesTotal)
	for i := range GridLinesTotal {
		floatI := float64(i) / floatLines
		deBruijnX[i] = math.Sin(2 * math.Pi * floatI)
		deBruijnY[i] = math.Cos(2 * math.Pi * floatI)
	}
}

func (gp *GridPoint) GetCenterPoint() image.Point {
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

func (f *Field) GetCenterPoint(ga GridAxes) image.Point {
	off0, off1 := float64(ga.Coords.Offset0), float64(ga.Coords.Offset1)

	x := (0.5+off0)*deBruijnX[ga.Axis0] + (0.5+off1)*deBruijnX[ga.Axis1]
	y := (0.5+off0)*deBruijnY[ga.Axis0] + (0.5+off1)*deBruijnY[ga.Axis1]

	for _, otl := range f.offsetsToLast[ga.Axis0][ga.Axis1] {
		off := math.Ceil(otl.zeroZero + off0*otl.ax0Delta + off1*otl.ax1Delta)
		x += off * deBruijnX[otl.targetAx]
		y += off * deBruijnY[otl.targetAx]
	}

	return image.Point{X: int(x * deBruijnScale), Y: int(y * deBruijnScale)}
}

const _0 = 0.05
const _1 = 0.95

func (gp *GridPoint) GetCornerPoints() [4]image.Point {
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
