package pgrid

import (
	"image"
	"image/color"
)

func (f *Field) ModifiedImagesStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	currPoint, currLine, prevLine, prevPointSign, pointColor := f.initialState()
	initialPoint := currPoint.getCenterPoint()
	currentImage := image.NewRGBA(pointRect(initialPoint, 256))

	for range maxSteps {
		currPoint, currLine, prevLine, prevPointSign, pointColor = f.next(currPoint, currLine, prevLine, prevPointSign)
		point := currPoint.getCenterPoint()
		if isOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(pointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[pointColor])
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
