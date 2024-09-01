package pgrid

import (
	"image"
	"image/color"
)

type CommandType int

const (
	Reset CommandType = iota
)

func (f *Field) ControlledInfiniteStepper(modifiedImagesCh chan<- *image.RGBA, commandCh <-chan CommandType, palette []color.RGBA) {
	prevPoint, currPoint, prevLine, currLine, prevPointColor := f.initialState()
	initialPoint := f.getCenterPoint(&prevPoint)
	currentImage := image.NewRGBA(pointRect(initialPoint, 256))

	step := 0
	shouldReset := false
	shouldRun := true

	for shouldRun {
		prevPoint, currPoint, prevLine, currLine, prevPointColor = f.step(prevPoint, currPoint, prevLine, currLine)

		point := f.getCenterPoint(&prevPoint)
		if !shouldReset && isOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(pointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[prevPointColor])

		select {
		case command := <-commandCh:
			if command == Reset {
				shouldReset = true
			}
		default:
		}

		if shouldReset && prevPointColor == 0 {
			shouldReset = false

			currentImage = image.NewRGBA(pointRect(point, 256))

			ResetValues()
		}

		step += 1
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
