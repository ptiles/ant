package pgrid

import (
	"image"
	"image/color"
)

func (f *Field) ModifiedImagesStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	prevPointPoint, currPoint, prevLine, currLine, currPointColor := f.initialState()
	initialPoint := currPoint.getCenterPoint()
	currentImage := image.NewRGBA(pointRect(initialPoint, 256))

	for range maxSteps {
		prevPointPoint, currPoint, prevLine, currLine, currPointColor = f.next(prevPointPoint, currPoint, prevLine, currLine)
		point := currPoint.getCenterPoint()
		if isOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(pointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[currPointColor])
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
