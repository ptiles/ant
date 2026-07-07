package ximage

import (
	"image"
	"image/color"

	"github.com/StephaneBunel/bresenham"
)

func DrawSegment(img bresenham.Plotter, segment [2]image.Point, c color.Color) {
	bresenham.DrawLine(img, segment[0].X, segment[0].Y, segment[1].X, segment[1].Y, c)
}

func DrawQuad(img bresenham.Plotter, quad [4]image.Point, c color.Color) {
	bresenham.DrawLine(img, quad[0].X, quad[0].Y, quad[1].X, quad[1].Y, c)
	bresenham.DrawLine(img, quad[1].X, quad[1].Y, quad[2].X, quad[2].Y, c)
	bresenham.DrawLine(img, quad[2].X, quad[2].Y, quad[3].X, quad[3].Y, c)
	bresenham.DrawLine(img, quad[3].X, quad[3].Y, quad[0].X, quad[0].Y, c)
}

func DrawSquare(img bresenham.Plotter, center image.Point, halfSize int, c color.Color) {
	DrawQuad(img, [4]image.Point{
		{X: center.X - halfSize, Y: center.Y - halfSize},
		{X: center.X + halfSize, Y: center.Y - halfSize},
		{X: center.X + halfSize, Y: center.Y + halfSize},
		{X: center.X - halfSize, Y: center.Y + halfSize},
	}, c)
}

func DrawSquareThick(img bresenham.Plotter, center image.Point, halfSize, border int, c color.Color) {
	for b := range border {
		DrawSquare(img, center, halfSize+b, c)
	}
}
