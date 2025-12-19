package pgrid

import (
	"github.com/ptiles/ant/geom"
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

var deBruijn [GridLinesTotal]geom.Point

func init() {
	floatLines := float64(GridLinesTotal)
	for i := range GridLinesTotal {
		floatI := float64(i) / floatLines
		deBruijn[i].X = math.Cos(2 * math.Pi * floatI)
		deBruijn[i].Y = math.Sin(2 * math.Pi * floatI)
	}
}

func (gp *GridPoint) GetCenterPoint() image.Point {
	x := 0.5*deBruijn[gp.Axes.Axis0].X + 0.5*deBruijn[gp.Axes.Axis1].X
	y := 0.5*deBruijn[gp.Axes.Axis0].Y + 0.5*deBruijn[gp.Axes.Axis1].Y

	for i := range GridLinesTotal {
		// floatOffset := float64(gp.Offsets[i]) * inflation
		floatOffset := float64(gp.Offsets[i])
		x += floatOffset * deBruijn[i].X
		y += floatOffset * deBruijn[i].Y
	}

	return image.Point{X: int(x * deBruijnScale), Y: int(y * deBruijnScale)}
}

func (f *Field) GetCenterPoint(ga GridAxes) image.Point {
	off0, off1 := float64(ga.Coords.Offset0), float64(ga.Coords.Offset1)

	x := (0.5+off0)*deBruijn[ga.Axis0].X + (0.5+off1)*deBruijn[ga.Axis1].X
	y := (0.5+off0)*deBruijn[ga.Axis0].Y + (0.5+off1)*deBruijn[ga.Axis1].Y

	for _, delta := range f.geometry[ga.Axis0][ga.Axis1].deltas {
		off := math.Ceil(delta.zeroZero + off0*delta.ax0Delta + off1*delta.ax1Delta)
		x += off * deBruijn[delta.targetAx].X
		y += off * deBruijn[delta.targetAx].Y
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
		x += floatOffset * deBruijn[i].X
		y += floatOffset * deBruijn[i].Y
	}

	// TODO: prepare this in init() and store in counter-clockwise order
	dax0x := deBruijn[gp.Axes.Axis0].X
	dax0y := deBruijn[gp.Axes.Axis0].Y
	dax1x := deBruijn[gp.Axes.Axis1].X
	dax1y := deBruijn[gp.Axes.Axis1].Y

	return [4]image.Point{
		{X: int((x + _0*dax0x + _0*dax1x) * deBruijnScale), Y: int((y + _0*dax0y + _0*dax1y) * deBruijnScale)},
		{X: int((x + _0*dax0x + _1*dax1x) * deBruijnScale), Y: int((y + _0*dax0y + _1*dax1y) * deBruijnScale)},
		{X: int((x + _1*dax0x + _1*dax1x) * deBruijnScale), Y: int((y + _1*dax0y + _1*dax1y) * deBruijnScale)},
		{X: int((x + _1*dax0x + _0*dax1x) * deBruijnScale), Y: int((y + _1*dax0y + _0*dax1y) * deBruijnScale)},
	}
}
