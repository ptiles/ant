package main

import (
	"image"
	"image/color"
	"math"

	"github.com/ptiles/ant/utils/ximage"
	"github.com/ptiles/ant/wgrid"
)

func DrawGrid(wg *wgrid.WythoffGrid, gridImage *image.NRGBA, minColumn, maxColumn int, c color.RGBA) {
	for edgePoints := range wg.EdgePoints(minColumn, maxColumn) {
		ximage.DrawSegment(gridImage, edgePoints, c)
	}
}

func DrawMultiGrid(wg *wgrid.WythoffGrid, gridImage *image.NRGBA, gridSize int) {
	DrawGrid(wg, gridImage, gridSize, gridSize+3, color.RGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff})
	DrawGrid(wg, gridImage, gridSize+3, gridSize+5, color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff})
	DrawGrid(wg, gridImage, gridSize+5, math.MaxInt, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
}
