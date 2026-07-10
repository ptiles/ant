package pgrid

import (
	"image"

	"github.com/ptiles/ant/geom"
)

var linear [GridLinesTotal][GridLinesTotal]geom.Point

func init() {
	for ax0, ax1 := range AxesAll() {
		linear[ax0][ax1].X = LineScale * threeAxesOffsetDeg(ax0, ax1, 0)
		linear[ax0][ax1].Y = LineScale * threeAxesOffsetDeg(ax0, ax1, 90)
	}
}

func (ga *GridAxes) GetCenterPoint() image.Point {
	off0 := float64(ga.Coords.Offset0)
	off1 := float64(ga.Coords.Offset1)

	return image.Point{
		X: int(off0*linear[ga.Axis0][ga.Axis1].X +
			off1*linear[ga.Axis1][ga.Axis0].X),
		Y: int(off0*linear[ga.Axis0][ga.Axis1].Y +
			off1*linear[ga.Axis1][ga.Axis0].Y),
	}
}
