package output

import (
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/utils/ximage"
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
}

func NewImage(rectangle image.Rectangle, scaleFactor, maxDimension int) *Image {
	rectS := ximage.RectDiv(rectangle, scaleFactor)
	return &Image{
		ResultRectN:  rectangle,
		imageRectN:   rectangle,
		imageS:       image.NewRGBA(rectS),
		ScaleFactor:  scaleFactor,
		maxDimension: maxDimension,
		dynamic:      rectangle.Empty(),
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
	i.imageRectN = ximage.RectGrow(i.ResultRectN, i.maxDimension)
	imageRectS := ximage.RectDiv(i.imageRectN, i.ScaleFactor)

	newResultImageS := image.NewRGBA(imageRectS)
	draw.ApproxBiLinear.Scale(
		newResultImageS, ximage.RectDiv(i.imageS.Rect, 2),
		i.imageS, i.imageS.Rect,
		draw.Over, nil,
	)
	i.imageS = newResultImageS
}

func (i *Image) mergeImage(modifiedImage *image.RGBA) {
	draw.BiLinear.Scale(
		i.imageS, ximage.RectDiv(modifiedImage.Rect, i.ScaleFactor),
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
		}
	}
	i.mergeImage(modifiedImage)
}

func (i *Image) SaveImages(fileName string, keepAlpha bool) image.Rectangle {
	outputRectS := ximage.RectDiv(i.outputRectN(), i.ScaleFactor)
	croppedRect := image.Rectangle{Min: image.Point{}, Max: outputRectS.Size()}

	croppedImg := image.NewNRGBA(croppedRect)
	sourcePoint := outputRectS.Min

	draw.Draw(croppedImg, croppedImg.Rect, i.imageS, sourcePoint, draw.Over)

	if !keepAlpha {
		ximage.RemoveAlpha(croppedImg)
	}

	ximage.SavePNG(croppedImg, fileName)

	return croppedImg.Rect
}
