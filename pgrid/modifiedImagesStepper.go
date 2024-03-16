package pgrid

import (
	"image"
	"image/color"
)

func isOutside(x, y int, rect image.Rectangle) bool {
	return x < rect.Min.X+16 || y < rect.Min.Y+16 || x > rect.Max.X-16 || y > rect.Max.Y-16
}

func (f *Field) ModifiedImagesStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	prevPoint, currPoint, prevLine, currLine, prevPointColor := f.initialState()
	xs, ys := f.getCenterPoint(&prevPoint)
	currentImage := image.NewRGBA(image.Rect(xs-128, ys-128, xs+128, ys+128))

	for step := 0; step < maxSteps; step++ {
		prevPoint, currPoint, prevLine, currLine, prevPointColor = f.step(prevPoint, currPoint, prevLine, currLine)
		x, y := f.getCenterPoint(&prevPoint)
		if isOutside(x, y, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(image.Rect(x-128, y-128, x+128, y+128))
		}
		currentImage.Set(x, y, palette[prevPointColor])

	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
