package pgrid

import (
	"image"
	"math"

	"github.com/ptiles/ant/geom"
)

var linear [GridLinesTotal][GridLinesTotal]geom.Point

func init() {
	for ax0, ax1 := range AxesAll() {
		linear[ax0][ax1].X = LineScale * threeAxesOffsetDeg(ax0, ax1, 0)
		linear[ax0][ax1].Y = LineScale * threeAxesOffsetDeg(ax0, ax1, 90)
	}
}

func (ga *GridAxes) GetCenterPointFloat() (float64, float64) {
	off0 := float64(ga.Coords.Offset0)
	off1 := float64(ga.Coords.Offset1)

	return off0*linear[ga.Axis0][ga.Axis1].X + off1*linear[ga.Axis1][ga.Axis0].X,
		off0*linear[ga.Axis0][ga.Axis1].Y + off1*linear[ga.Axis1][ga.Axis0].Y
}

func (ga *GridAxes) GetCenterPoint() image.Point {
	x, y := ga.GetCenterPointFloat()

	return image.Point{X: int(x), Y: int(y)}
}

func (ga *GridAxes) GetCenterPointRound() image.Point {
	x, y := ga.GetCenterPointFloat()

	return image.Point{X: int(math.Round(x)), Y: int(math.Round(y))}
}
