package pgrid

import (
	"image"
	"math"
)

//const drawTilesAndPoints = true
//const deBruijnScale = 48
//const padding = deBruijnScale

const drawTilesAndPoints = false
const deBruijnScale = 2
const padding = deBruijnScale * 4

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

const _0 = 0.05
const _1 = 0.95

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
