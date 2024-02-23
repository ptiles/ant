package pgrid

import (
	"github.com/ptiles/ant/store"
)

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func (f *Field) GetMinMax(minWidth, minHeight int) (int, int) {
	maxX := minWidth / 2
	maxY := minHeight / 2

	//// TODO: make more efficient per-axis-pair minima and maxima
	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := GridLine{Axis: axis0, Offset: off0}
		line1 := GridLine{Axis: axis1, Offset: off1}
		gp := f.MakeGridPoint(line0, line1)
		x, y := f.GetCenterPoint(&gp)
		pointX := absInt(x)
		if pointX > maxX {
			maxX = pointX
		}
		pointY := absInt(y)
		if pointY > maxY {
			maxY = pointY
		}
	})

	maxX = (maxX / 128) * 128
	maxY = (maxY / 128) * 128

	return maxX, maxY
}
