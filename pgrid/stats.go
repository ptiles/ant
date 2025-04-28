package pgrid

import (
	"github.com/ptiles/ant/utils"
	"image"
	"iter"
)

func Uniq() (uint64, int) {
	uPoints := uint64(0)
	uMaps := 0

	for ax0, ax1 := range AxesCanon() {
		uArr := aValues[ax0][ax1]
		for _, dMap := range uArr.Maps {
			if dMap != nil {
				uPoints += uint64(len(dMap))
				uMaps += 1
			}
		}
	}

	return uPoints, uMaps
}

func (f *Field) Rect() image.Rectangle {
	rect := image.Rectangle{}
	for ax0, ax1 := range AxesCanon() {
		rect = rect.Union(f.RectUpArray(ax0, ax1))
	}
	return rect
}

func (f *Field) RectUpArray(ax0, ax1 uint8) image.Rectangle {
	rect := image.Rectangle{}
	uArr := aValues[ax0][ax1]

	for _, minPoint := range uArr.MinCoords() {
		rect = f.RectUnionMinPoint(rect, ax0, ax1, minPoint)
	}

	//if rect.Empty() {
	//	return rect
	//}

	return utils.SnapRect(rect, Padding)
}

func (f *Field) rectCornerPoints(ax0, ax1 uint8, minPoint GridCoords) (image.Point, image.Point, image.Point, image.Point) {
	rSize := offsetInt(downMask)

	ga00 := GridAxes{Axis0: ax0, Axis1: ax1, Coords: GridCoords{
		Offset0: minPoint.Offset0, Offset1: minPoint.Offset1,
	}}
	ga01 := GridAxes{Axis0: ax0, Axis1: ax1, Coords: GridCoords{
		Offset0: minPoint.Offset0, Offset1: minPoint.Offset1 + rSize,
	}}
	ga10 := GridAxes{Axis0: ax0, Axis1: ax1, Coords: GridCoords{
		Offset0: minPoint.Offset0 + rSize, Offset1: minPoint.Offset1,
	}}
	ga11 := GridAxes{Axis0: ax0, Axis1: ax1, Coords: GridCoords{
		Offset0: minPoint.Offset0 + rSize, Offset1: minPoint.Offset1 + rSize,
	}}

	return f.GetCenterPoint(ga00), f.GetCenterPoint(ga01), f.GetCenterPoint(ga10), f.GetCenterPoint(ga11)
}

func (f *Field) RectUnionMinPoint(rect image.Rectangle, ax0, ax1 uint8, minPoint GridCoords) image.Rectangle {
	ga00, ga01, ga10, ga11 := f.rectCornerPoints(ax0, ax1, minPoint)

	r0 := image.Rectangle{Min: ga00, Max: ga11}.Canon()
	r1 := image.Rectangle{Min: ga01, Max: ga10}.Canon()

	return rect.Union(r0.Union(r1))
}

func (f *Field) RectIntersectMinPoint(rect image.Rectangle, ax0, ax1 uint8, minPoint GridCoords) image.Rectangle {
	ga00, ga01, ga10, ga11 := f.rectCornerPoints(ax0, ax1, minPoint)

	r0 := image.Rectangle{Min: ga00, Max: ga11}.Canon()
	r1 := image.Rectangle{Min: ga01, Max: ga10}.Canon()

	return rect.Intersect(r0.Union(r1))
}

func (ua *upArray) MinCoords() iter.Seq2[upInt, GridCoords] {
	return func(yield func(upInt, GridCoords) bool) {
		for off1 := range ua.Max.Offset1 - ua.Min.Offset1 {
			for off0 := range ua.Max.Offset0 - ua.Min.Offset0 {
				i := off1*ua.Stride + off0
				//if len(ua.Maps[i]) > 0 {
				if ua.Maps[i] != nil {
					coords := GridCoords{
						Offset0: offsetInt((off0 + ua.Min.Offset0) << bits),
						Offset1: offsetInt((off1 + ua.Min.Offset1) << bits),
					}
					if !yield(i, coords) {
						return
					}
				}
			}
		}
	}
}
