package pgrid

import (
	"github.com/ptiles/ant/seq"
	"github.com/ptiles/ant/utils"
	"image"
	"iter"
	"math"
	"sort"
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

type Bounds [GridLinesTotal]struct {
	Axis     uint8     `json:"axis"`
	Min      offsetInt `json:"min"`
	MinCount int       `json:"minCount"`
	Max      offsetInt `json:"max"`
	MaxCount int       `json:"maxCount"`

	Counts OffsetCounts `json:"counts"`
}

type BoundsSize [GridLinesTotal]offsetInt

func GetBounds(limit int) (Bounds, BoundsSize, offsetInt, offsetInt) {
	bounds := Bounds{}
	for ax := range GridLinesTotal {
		bounds[ax].Axis = ax
		bounds[ax].Min = math.MaxInt32
		bounds[ax].Max = math.MinInt32
	}

	for ax0, ax1 := range AxesCanon() {
		uArr := aValues[ax0][ax1]
		for i, dMap := range uArr.Maps {
			baseOff0 := offsetInt(uArr.Min.Offset0+upInt(i)%uArr.Stride) << bits
			baseOff1 := offsetInt(uArr.Min.Offset1+upInt(i)/uArr.Stride) << bits
			for dCoord := range dMap {
				off0 := offsetInt(dCoord.Offset0) + baseOff0
				off1 := offsetInt(dCoord.Offset1) + baseOff1

				if off0 < bounds[ax0].Min {
					bounds[ax0].Min = off0
					bounds[ax0].MinCount = 1
				} else if off0 == bounds[ax0].Min {
					bounds[ax0].MinCount += 1
				}
				if off1 < bounds[ax1].Min {
					bounds[ax1].Min = off1
					bounds[ax1].MinCount = 1
				} else if off1 == bounds[ax1].Min {
					bounds[ax1].MinCount += 1
				}

				if off0 > bounds[ax0].Max {
					bounds[ax0].Max = off0
					bounds[ax0].MaxCount = 1
				} else if off0 == bounds[ax0].Max {
					bounds[ax0].MaxCount += 1
				}
				if off1 > bounds[ax1].Max {
					bounds[ax1].Max = off1
					bounds[ax1].MaxCount = 1
				} else if off1 == bounds[ax1].Max {
					bounds[ax1].MaxCount += 1
				}
			}
		}
	}

	sizeMin := offsetInt(math.MaxInt32)
	sizeMax := offsetInt(math.MinInt32)
	sizes := BoundsSize{}
	for ax := range GridLinesTotal {
		diff := bounds[ax].Max - bounds[ax].Min
		sizes[ax] = diff
		if diff < sizeMin {
			sizeMin = diff
		}
		if diff > sizeMax {
			sizeMax = diff
		}
	}

	topCounts := TopCounts(limit)
	for ax := range GridLinesTotal {
		bounds[ax].Counts = topCounts[ax]
	}

	return bounds, sizes, sizeMin, sizeMax
}

type OffsetCount struct {
	Offset offsetInt `json:"offset"`
	Count  int       `json:"count"`
	Row    int       `json:"wythoffRow,omitempty"`
	Col    int       `json:"wythoffCol,omitempty"`
}
type OffsetCounts []OffsetCount

func (oc OffsetCounts) Len() int           { return len(oc) }
func (oc OffsetCounts) Less(i, j int) bool { return oc[i].Count > oc[j].Count }
func (oc OffsetCounts) Swap(i, j int)      { oc[i], oc[j] = oc[j], oc[i] }

func TopCounts(limit int) [GridLinesTotal]OffsetCounts {
	countsMap := [GridLinesTotal]map[offsetInt]int{}
	for ax := range GridLinesTotal {
		countsMap[ax] = make(map[offsetInt]int)
	}

	for ax0, ax1 := range AxesCanon() {
		uArr := aValues[ax0][ax1]
		for i, dMap := range uArr.Maps {
			baseOff0 := offsetInt(uArr.Min.Offset0+upInt(i)%uArr.Stride) << bits
			baseOff1 := offsetInt(uArr.Min.Offset1+upInt(i)/uArr.Stride) << bits
			for dCoord := range dMap {
				off0 := offsetInt(dCoord.Offset0) + baseOff0
				off1 := offsetInt(dCoord.Offset1) + baseOff1

				countsMap[ax0][off0] += 1
				countsMap[ax1][off1] += 1
			}
		}
	}

	var topCounts [GridLinesTotal]OffsetCounts
	for ax := range GridLinesTotal {
		if limit > 0 {
			topCounts[ax] = axisCounts(countsMap[ax])[:limit]
		} else {
			topCounts[ax] = axisCounts(countsMap[ax])
		}
	}

	return topCounts
}

func axisCounts(countsMap map[offsetInt]int) OffsetCounts {
	countsSlice := make(OffsetCounts, len(countsMap))
	i := 0
	for offset, count := range countsMap {
		rev := seq.WythoffReverse[int(offset)]
		countsSlice[i] = OffsetCount{Offset: offset, Count: count, Row: rev.Row, Col: rev.Col}
		i++
	}
	sort.Sort(countsSlice)
	return countsSlice
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
