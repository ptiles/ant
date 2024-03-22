package pgrid

import (
	"image"
	"image/color"
)

func pointSnap(p image.Point, gridSize int) image.Point {
	return image.Point{X: p.X / gridSize * gridSize, Y: p.Y / gridSize * gridSize}
}

func rectGrow(r image.Rectangle, maxDimension int) image.Rectangle {
	halfImage := image.Point{X: maxDimension / 2, Y: maxDimension / 2}
	centerPoint := r.Min.Add(r.Max).Div(2)
	return image.Rectangle{Min: centerPoint.Sub(halfImage), Max: centerPoint.Add(halfImage)}
}

func pointRect(p image.Point, maxDimension int) image.Rectangle {
	gridSize := 16
	snapped := pointSnap(p, gridSize)
	return rectGrow(image.Rectangle{Min: snapped, Max: snapped.Add(image.Point{X: gridSize, Y: gridSize})}, maxDimension)
}

func isOutside(point image.Point, rect image.Rectangle) bool {
	return point.X < rect.Min.X+16 || point.Y < rect.Min.Y+16 ||
		point.X > rect.Max.X-16 || point.Y > rect.Max.Y-16
}

func (f *Field) ModifiedImagesStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	prevPoint, currPoint, prevLine, currLine, prevPointColor := f.initialState()
	initialPoint := f.getCenterPoint(&prevPoint)
	currentImage := image.NewRGBA(pointRect(initialPoint, 256))

	for step := 0; step < maxSteps; step++ {
		prevPoint, currPoint, prevLine, currLine, prevPointColor = f.step(prevPoint, currPoint, prevLine, currLine)
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
