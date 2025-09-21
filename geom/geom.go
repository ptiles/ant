package geom

import (
	"image"
	"math"
)

type Point struct{ X, Y float64 }

func (p Point) Mul(scaleFactor int) Point {
	return Point{X: p.X * float64(scaleFactor), Y: p.Y * float64(scaleFactor)}
}

func (p Point) Round() image.Point {
	return image.Point{X: int(math.Round(p.X)), Y: int(math.Round(p.Y))}
}

type Line struct{ A, B Point }

func (l Line) SegmentContains(p Point) bool {
	lineRect := image.Rect(int(l.A.X), int(l.A.Y), int(l.B.X), int(l.B.Y)).Inset(-1)

	return p.Round().In(lineRect)
}

func FromDeg(deg float64) float64 {
	return deg / 180 * math.Pi
}

func Sin(deg float64) float64 {
	return math.Sin(FromDeg(deg))
}

func Cos(deg float64) float64 {
	return math.Cos(FromDeg(deg))
}

func Intersection(line1, line2 Line) Point {
	line1pointA, line1pointB := line1.A, line1.B
	line2pointA, line2pointB := line2.A, line2.B

	x1A, y1A := line1pointA.X, line1pointA.Y
	x1B, y1B := line1pointB.X, line1pointB.Y
	x2A, y2A := line2pointA.X, line2pointA.Y
	x2B, y2B := line2pointB.X, line2pointB.Y

	dx1 := x1A - x1B
	dy1 := y1A - y1B
	dx2 := x2A - x2B
	dy2 := y2A - y2B

	den := dx1*dy2 - dy1*dx2

	return Point{
		X: ((x1A*y1B-y1A*x1B)*dx2 - dx1*(x2A*y2B-y2A*x2B)) / den,
		Y: ((x1A*y1B-y1A*x1B)*dy2 - dy1*(x2A*y2B-y2A*x2B)) / den,
	}
}

func NewPoint(p image.Point) Point {
	return Point{X: float64(p.X), Y: float64(p.Y)}
}

func Distance(line Line, point Point) float64 {
	x1, y1 := line.A.X, line.A.Y
	x2, y2 := line.B.X, line.B.Y

	x0, y0 := point.X, point.Y

	x10 := x1 - x0
	y10 := y1 - y0
	x21 := x2 - x1
	y21 := y2 - y1

	return x21*y10 - x10*y21
}
