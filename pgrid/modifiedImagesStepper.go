package pgrid

import (
	"image"
	"image/color"
)

func (f *Field) ModifiedImagesStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	prevPoint, currPoint, prevLine, currLine, prevPointColor := f.initialState()
	initialPoint := f.getCenterPoint(&prevPoint)
	currentImage := image.NewRGBA(pointRect(initialPoint, 256))

	for range maxSteps {
		prevPoint, currPoint, prevLine, currLine, prevPointColor = f.next(prevPoint, currPoint, prevLine, currLine)
		point := f.getCenterPoint(&prevPoint)
		if isOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(pointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[prevPointColor])
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
