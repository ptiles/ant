package wgrid

import (
	"fmt"
	"github.com/StephaneBunel/bresenham"
	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/pgrid/axis"
	"github.com/ptiles/ant/seq"
	"github.com/ptiles/ant/utils/ximage"
	"image"
	"image/color"
	"math"
	"slices"
	"strings"
)

type WythoffGrid struct {
	RectN       image.Rectangle
	RectS       image.Rectangle
	ScaleFactor int
	Edges       []geom.Line
	Ranges      [pgrid.GridLinesTotal]Range
}

type Range struct {
	Min int
	Max int
}

func New(rectN image.Rectangle, scaleFactor int) WythoffGrid {
	rectS := ximage.RectDiv(rectN, scaleFactor)

	var ranges [pgrid.GridLinesTotal]Range
	for ax := range pgrid.GridLinesTotal {
		rMin, rMax := axisRange(ax, rectN)
		ranges[ax].Min = rMin
		ranges[ax].Max = rMax
	}

	return WythoffGrid{
		RectN:       rectN,
		RectS:       rectS,
		ScaleFactor: scaleFactor,
		Edges:       edgeLines(rectS),
		Ranges:      ranges,
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

func axisLine(ax uint8, offset int, scaleFactor float64) geom.Line {
	aX, aY := float64(offset)*pgrid.LineScale, float64(-1_000)
	bX, bY := float64(offset)*pgrid.LineScale, float64(1_000)
	sin := sinCos[ax].sin
	cos := sinCos[ax].cos

	return geom.Line{
		A: geom.Point{X: (aX*cos - aY*sin) / scaleFactor, Y: (aX*sin + aY*cos) / scaleFactor},
		B: geom.Point{X: (bX*cos - bY*sin) / scaleFactor, Y: (bX*sin + bY*cos) / scaleFactor},
	}
}

func zeroAxUnitLine(ax uint8) geom.Line {
	return axisLine(ax, 0, 2000)
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

func (wg WythoffGrid) Intersection(ax0 uint8, off0 int, ax1 uint8, off1 int) (image.Point, bool) {
	point := Intersection(ax0, off0, ax1, off1, wg.ScaleFactor)
	return point, point.In(wg.RectS)
}

func Intersection(ax0 uint8, off0 int, ax1 uint8, off1 int, scaleFactor int) image.Point {
	line0 := axisLine(ax0, off0, float64(scaleFactor))
	line1 := axisLine(ax1, off1, float64(scaleFactor))

	return geom.Intersection(line0, line1).Round()
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
		minOffset0, maxOffset0 := wg.Ranges[ax0].Min, wg.Ranges[ax0].Max
		minOffset1, maxOffset1 := wg.Ranges[ax1].Min, wg.Ranges[ax1].Max

		for off0 := range seq.WythoffMinMaxColumn(minOffset0, maxOffset0, minColumn, maxColumn) {
			for off1 := range seq.WythoffMinMaxColumn(minOffset1, maxOffset1, minColumn, maxColumn) {
				if point, in := wg.Intersection(ax0, off0, ax1, off1); in {
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

func (wg WythoffGrid) drawAxisSegment(ax uint8, off int, gridImage *image.NRGBA, c color.RGBA) {
	axLine := axisLine(ax, off, float64(wg.ScaleFactor))

	var edgePoints [2]image.Point
	j := 0

	for _, edge := range wg.Edges {
		edgePoint := geom.Intersection(axLine, edge)
		if edge.SegmentContains(edgePoint) {
			edgePoints[j] = edgePoint.Round()
			j += 1
		}
		if j == 2 {
			bresenham.DrawLine(gridImage, edgePoints[0].X, edgePoints[0].Y, edgePoints[1].X, edgePoints[1].Y, c)
			return
		}
	}
}

func (wg WythoffGrid) DrawGrid(gridImage *image.NRGBA, c color.RGBA, minColumn, maxColumn int) {
	for ax := range pgrid.GridLinesTotal {
		minOffset, maxOffset := wg.Ranges[ax].Min, wg.Ranges[ax].Max
		for off := range seq.WythoffMinMaxColumn(minOffset, maxOffset, minColumn, maxColumn) {
			wg.drawAxisSegment(ax, off, gridImage, c)
		}
	}
}

func (wg WythoffGrid) DrawMultiGrid(gridSize int) *image.NRGBA {
	gridImage := image.NewNRGBA(wg.RectS)

	wg.DrawGrid(gridImage, color.RGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff}, gridSize, gridSize+3)
	wg.DrawGrid(gridImage, color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff}, gridSize+3, gridSize+5)
	wg.DrawGrid(gridImage, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, gridSize+5, math.MaxInt)

	return gridImage
}
