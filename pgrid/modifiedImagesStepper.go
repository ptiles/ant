package pgrid

import (
	"github.com/ptiles/ant/utils"
	"image"
	"image/color"
)

func (f *Field) ModifiedImagesStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	currPoint, currLine, prevLine, prevPointSign := f.initialState()
	var pointColor uint8
	initialPoint := currPoint.getCenterPoint()
	currentImage := image.NewRGBA(utils.PointRect(initialPoint, 256))

	for range maxSteps {
		currPoint, currLine, prevLine, prevPointSign, pointColor = f.next(currPoint, currLine, prevLine, prevPointSign)
		point := currPoint.getCenterPoint()
		if utils.IsOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(utils.PointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[pointColor])
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
