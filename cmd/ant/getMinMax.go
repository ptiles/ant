package main

import (
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/store"
	"math"
)

func getMinMax(field *pgrid.Field, minWidth, minHeight int) (int, int) {
	maxX := minWidth / 2
	maxY := minHeight / 2

	//// TODO: make more accurate per-axis-pair minima and maxima
	//offsetPairs := [4][2]int16{
	//	{store.MinOffset0, store.MinOffset1},
	//	{store.MinOffset0, store.MaxOffset1},
	//	{store.MaxOffset0, store.MinOffset1},
	//	{store.MaxOffset0, store.MaxOffset1},
	//}
	//for pa := uint8(0); pa < 10; pa++ {
	//	axis0, axis1 := store.UnpackAxes(pa)
	//	for _, offsetPair := range offsetPairs {
	//		off0, off1 := offsetPair[0], offsetPair[1]
	//		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
	//		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
	//		gp := field.MakeGridPoint(line0, line1)
	//		point := field.GetCenterPoint(&gp)
	//		pointX := int(math.Abs(point[0]))
	//		if pointX > maxX {
	//			maxX = pointX
	//		}
	//		pointY := int(math.Abs(point[1]))
	//		if pointY > maxY {
	//			maxY = pointY
	//		}
	//	}
	//}

	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
		gp := field.MakeGridPoint(line0, line1)
		point := field.GetCenterPoint(&gp)
		pointX := int(math.Abs(point[0]))
		if pointX > maxX {
			maxX = pointX
		}
		pointY := int(math.Abs(point[1]))
		if pointY > maxY {
			maxY = pointY
		}
	})

	maxX = (maxX/128 + 1) * 128
	maxY = (maxY/128 + 1) * 128

	return maxX, maxY
}
