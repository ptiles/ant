package wgrid

import (
	"fmt"
	"image"
	"iter"
	"math"
	"slices"
	"strings"

	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/pgrid/axis"
	"github.com/ptiles/ant/seq"
)

type WythoffGrid struct {
	rect   image.Rectangle
	edges  []geom.Line
	Ranges [pgrid.GridLinesTotal]Range
}

type Range struct {
	Min     int
	Max     int
	reverse seq.NumColSlice
}

func New(rect image.Rectangle) WythoffGrid {
	var ranges [pgrid.GridLinesTotal]Range
	for ax := range pgrid.GridLinesTotal {
		rMin, rMax := axisRange(ax, rect)
		ranges[ax].Min = rMin
		ranges[ax].Max = rMax
		ranges[ax].reverse = seq.WythoffReverseSorted.Slice(rMin, rMax)
	}

	return WythoffGrid{
		rect:   rect,
		edges:  edgeLines(rect),
		Ranges: ranges,
	}
}

func edgeLines(rect image.Rectangle) []geom.Line {
	if rect.Empty() {
		return nil
	}

	nw := geom.Point{X: float64(rect.Min.X), Y: float64(rect.Min.Y)}
	ne := geom.Point{X: float64(rect.Max.X), Y: float64(rect.Min.Y)}
	sw := geom.Point{X: float64(rect.Min.X), Y: float64(rect.Max.Y)}
	se := geom.Point{X: float64(rect.Max.X), Y: float64(rect.Max.Y)}

	return []geom.Line{
		{nw, ne}, {ne, se}, {sw, se}, {nw, sw},
	}
}

var sinCos [pgrid.GridLinesTotal]struct{ sin, cos float64 }

func init() {
	alpha := 360 / float64(pgrid.GridLinesTotal)

	for ax := range pgrid.GridLinesTotal {
		alphaAx := alpha * float64(ax)
		sinCos[ax].sin = geom.Sin(alphaAx)
		sinCos[ax].cos = geom.Cos(alphaAx)
	}
}

func axisLine(ax uint8, offset int) geom.Line {
	aX, aY := float64(offset)*pgrid.LineScale, -0.5
	bX, bY := float64(offset)*pgrid.LineScale, 0.5
	sin := sinCos[ax].sin
	cos := sinCos[ax].cos

	return geom.Line{
		A: geom.Point{X: aX*cos - aY*sin, Y: aX*sin + aY*cos},
		B: geom.Point{X: bX*cos - bY*sin, Y: bX*sin + bY*cos},
	}
}

func zeroAxUnitLine(ax uint8) geom.Line {
	return axisLine(ax, 0)
}

func axisRange(ax uint8, rect image.Rectangle) (int, int) {
	cornerPoints := [4]image.Point{
		rect.Min, rect.Max,
		{X: rect.Min.X, Y: rect.Max.Y},
		{X: rect.Max.X, Y: rect.Min.Y},
	}

	cornerOffsets := make([]float64, 4)
	axLine := zeroAxUnitLine(ax)

	for j, cornerPoint := range cornerPoints {
		cornerOffsets[j] = geom.Distance(geom.NewPoint(cornerPoint), axLine) / pgrid.LineScale
	}

	minOff := int(math.Floor(slices.Min(cornerOffsets))) - 10
	maxOff := int(math.Ceil(slices.Max(cornerOffsets))) + 10

	return minOff, maxOff
}

func (wg WythoffGrid) Contains(point image.Point) bool {
	return point.In(wg.rect)
}

type AxesMap map[uint8]int

func (am AxesMap) String() string {
	var sb strings.Builder

	for ax := range pgrid.GridLinesTotal {
		if offset, ok := am[ax]; ok {
			fmt.Fprintf(&sb, "  %s%7d  |", axis.Name[ax], offset)
		} else {
			fmt.Fprintf(&sb, "            |")
		}
	}

	return sb.String()
}

func (wg WythoffGrid) IntersectionsMap(minColumn, maxColumn int) map[image.Point]AxesMap {
	intersections := map[image.Point]AxesMap{}

	for ax0, ax1 := range pgrid.AxesCanon() {
		reverse0 := wg.Ranges[ax0].reverse
		reverse1 := wg.Ranges[ax1].reverse

		for off0 := range reverse0.MinMaxColumn(minColumn, maxColumn) {
			for off1 := range reverse1.MinMaxColumn(minColumn, maxColumn) {
				gridAxes := pgrid.GridAxes{
					Axis0: ax0, Axis1: ax1,
					Coords: pgrid.GridCoords{
						Offset0: pgrid.OffsetInt(off0), Offset1: pgrid.OffsetInt(off1),
					},
				}
				point := gridAxes.GetCenterPoint()

				if wg.Contains(point) {
					if intersections[point] == nil {
						intersections[point] = make(AxesMap, pgrid.GridLinesTotal)
					}
					intersections[point][ax0] = off0
					intersections[point][ax1] = off1
				}
			}
		}
	}

	return intersections
}

func (wg WythoffGrid) EdgePoints(minColumn, maxColumn int) iter.Seq[[2]image.Point] {
	return func(yield func([2]image.Point) bool) {
		var edgePoints [2]image.Point

		for ax := range pgrid.GridLinesTotal {
			reverse := wg.Ranges[ax].reverse

			for off := range reverse.MinMaxColumn(minColumn, maxColumn) {
				axLine := axisLine(ax, off)
				j := 0
				for _, edge := range wg.edges {
					edgePoint := geom.Intersection(axLine, edge)
					if edge.SegmentContains(edgePoint) {
						edgePoints[j] = edgePoint.Round()
						j += 1
					}
					if j == 2 {
						if !yield(edgePoints) {
							return
						}
						break
					}
				}
			}
		}
	}
}
