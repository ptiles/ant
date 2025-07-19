package output

import (
	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/utils"
	"golang.org/x/image/draw"
	"image"
)

type Image struct {
	ResultRectN  image.Rectangle // Area with active pixels
	imageRectN   image.Rectangle
	imageS       *image.RGBA
	ScaleFactor  int
	maxDimension int
	dynamic      bool
	edges        []geom.Line
}

func NewImage(rectangle image.Rectangle, scaleFactor, maxDimension int) *Image {
	rectS := utils.RectDiv(rectangle, scaleFactor)
	return &Image{
		ResultRectN:  rectangle,
		imageRectN:   rectangle,
		imageS:       image.NewRGBA(rectS),
		ScaleFactor:  scaleFactor,
		maxDimension: maxDimension,
		dynamic:      rectangle.Empty(),
		edges:        edgeLines(rectS),
	}
}

func (i *Image) outputRectN() image.Rectangle {
	if i.dynamic {
		return i.ResultRectN
	}
	return i.imageRectN
}

func (i *Image) RectCenteredString() string {
	return utils.RectCenteredString(i.ResultRectN, i.ScaleFactor)
}

func (i *Image) halveImage() {
	i.imageRectN = utils.RectGrow(i.ResultRectN, i.maxDimension)
	imageRectS := utils.RectDiv(i.imageRectN, i.ScaleFactor)

	newResultImageS := image.NewRGBA(imageRectS)
	draw.ApproxBiLinear.Scale(
		newResultImageS, utils.RectDiv(i.imageS.Rect, 2),
		i.imageS, i.imageS.Rect,
		draw.Over, nil,
	)
	i.imageS = newResultImageS
}

func (i *Image) mergeImage(modifiedImage *image.RGBA) {
	draw.BiLinear.Scale(
		i.imageS, utils.RectDiv(modifiedImage.Rect, i.ScaleFactor),
		modifiedImage, modifiedImage.Rect,
		draw.Over, nil,
	)
}

func (i *Image) Merge(modifiedImage *image.RGBA) {
	i.ResultRectN = i.ResultRectN.Union(modifiedImage.Rect)
	if i.dynamic {
		if !i.ResultRectN.In(i.imageRectN) {
			if i.ScaleFactor == 0 {
				i.ScaleFactor = 1
			} else {
				i.ScaleFactor *= 2
				i.maxDimension *= 2
			}
			i.halveImage()
			i.edges = edgeLines(i.imageS.Rect)
		}
	}
	i.mergeImage(modifiedImage)
}

func (i *Image) SaveImages(fileName, withGridFileName string, gridSize int, keepAlpha bool) image.Rectangle {
	cropped := newCropped(i)

	cropped.draw(i.imageS)
	if !keepAlpha {
		cropped.removeAlpha()
	}
	if fileName != "" {
		cropped.savePNG(fileName)
	}

	if withGridFileName != "" && gridSize > 0 {
		cropped.draw(i.DrawGrid(gridSize))
		cropped.savePNG(withGridFileName)
	}

	return cropped.dst.Rect
}

func (i *Image) SaveGridOnly(gridOnlyFileName string, gridSize int, keepAlpha bool) {
	cropped := newCropped(i)
	cropped.draw(i.DrawGrid(gridSize))
	if !keepAlpha {
		cropped.removeAlpha()
	}
	cropped.savePNG(gridOnlyFileName)
}
