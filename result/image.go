package result

import (
	"github.com/ptiles/ant/utils"
	"golang.org/x/image/draw"
	"image"
)

type Image struct {
	ResultRectN  image.Rectangle
	paddingRectN image.Rectangle
	resultImageS *image.RGBA
	ScaleFactor  int
	maxDimension int
	dynamic      bool
}

func NewImage(rectangle image.Rectangle, scaleFactor, maxDimension int) (i Image) {
	if rectangle.Empty() {
		i.dynamic = true
		i.resultImageS = image.NewRGBA(image.Rectangle{})
	} else {
		i.ResultRectN = rectangle
		i.ScaleFactor = scaleFactor
		i.resultImageS = image.NewRGBA(utils.RectDiv(rectangle, scaleFactor))
	}
	i.maxDimension = maxDimension

	return
}

func halveImage(dst, src *image.RGBA) {
	draw.ApproxBiLinear.Scale(
		dst, utils.RectDiv(src.Rect, 2),
		src, src.Rect,
		draw.Over, nil,
	)
}

func mergeImage(dst, src *image.RGBA, scaleFactor int) {
	draw.BiLinear.Scale(
		dst, utils.RectDiv(src.Rect, scaleFactor),
		src, src.Rect,
		draw.Over, nil,
	)
}

func (i *Image) Merge(modifiedImage *image.RGBA) {
	if i.dynamic {
		i.ResultRectN = i.ResultRectN.Union(modifiedImage.Rect)
		if !i.ResultRectN.In(i.paddingRectN) {
			if i.ScaleFactor == 0 {
				i.ScaleFactor = 1
			} else {
				i.ScaleFactor *= 2
				i.maxDimension *= 2
			}
			i.paddingRectN = utils.RectGrow(i.ResultRectN, i.maxDimension)
			newResultImageS := image.NewRGBA(utils.RectDiv(i.paddingRectN, i.ScaleFactor))
			halveImage(newResultImageS, i.resultImageS)
			i.resultImageS = newResultImageS
		}
	}
	mergeImage(i.resultImageS, modifiedImage, i.ScaleFactor)
}

func (i *Image) Draw(keepAlpha bool) (*image.NRGBA, image.Rectangle) {
	resultRectS := utils.RectDiv(i.ResultRectN, i.ScaleFactor)
	croppedRect := image.Rectangle{Min: image.Point{}, Max: resultRectS.Size()}
	croppedImage := image.NewNRGBA(croppedRect)

	draw.Draw(croppedImage, croppedRect, i.resultImageS, resultRectS.Min, draw.Over)

	if !keepAlpha {
		for y := range croppedImage.Rect.Dy() {
			yOffset := y * croppedImage.Stride
			for x := range croppedImage.Rect.Dx() {
				croppedImage.Pix[yOffset+x*4+3] = 255
			}
		}
	}

	return croppedImage, resultRectS
}
