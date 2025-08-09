package output

import (
	"github.com/StephaneBunel/bresenham"
	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/seq"
	"image"
	"image/color"
	"math"
	"slices"
)

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

func AxisLine(ax uint8, offset int, scaleFactor float64) geom.Line {
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
	return AxisLine(ax, 0, 2000)
}

func (i *Image) AxisRange(ax uint8) (int, int) {
	return AxisRange(ax, i.outputRectN())
}

func AxisRange(ax uint8, rect image.Rectangle) (int, int) {
	cornerPoints := [4]image.Point{
		rect.Min, rect.Max,
		{X: rect.Min.X, Y: rect.Max.Y},
		{X: rect.Max.X, Y: rect.Min.Y},
	}

	cornerOffsets := make([]float64, 4)
	axLine := zeroAxUnitLine(ax)

	for j, cornerPoint := range cornerPoints {
		cornerOffsets[j] = geom.Distance(axLine, geom.NewPoint(cornerPoint)) / pgrid.LineScale
	}

	minOff := int(math.Floor(slices.Min(cornerOffsets))) - 10
	maxOff := int(math.Ceil(slices.Max(cornerOffsets))) + 10

	return minOff, maxOff
}

func (i *Image) DrawAxis(croppedImage *image.NRGBA, ax uint8, off int, color color.RGBA) {
	axLine := AxisLine(ax, off, float64(i.ScaleFactor))

	var edgePoints [2]geom.Point
	j := 0

	for _, edge := range i.edges {
		edgePoint := geom.Intersection(axLine, edge)
		if edge.SegmentContains(edgePoint) {
			edgePoints[j] = edgePoint
			j += 1
		}
		if j == 2 {
			bresenham.DrawLine(
				croppedImage,
				int(math.Round(edgePoints[0].X)), int(math.Round(edgePoints[0].Y)),
				int(math.Round(edgePoints[1].X)), int(math.Round(edgePoints[1].Y)),
				color,
			)
			return
		}
	}
}

func (i *Image) drawPartialGrid(gridImage *image.NRGBA, c color.RGBA, minColumn, maxColumn int) {
	for ax := range pgrid.GridLinesTotal {
		minOffset, maxOffset := i.AxisRange(ax)
		for off := range seq.WythoffMinMaxColumn(minOffset, maxOffset, minColumn, maxColumn) {
			i.DrawAxis(gridImage, ax, off, c)
		}
	}
}

func (i *Image) DrawGrid(gridSize int) *image.NRGBA {
	gridImage := image.NewNRGBA(i.imageS.Rect)

	i.drawPartialGrid(gridImage, color.RGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff}, gridSize, gridSize+3)
	i.drawPartialGrid(gridImage, color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff}, gridSize+3, gridSize+5)
	i.drawPartialGrid(gridImage, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, gridSize+5, math.MaxInt)

	return gridImage
}
