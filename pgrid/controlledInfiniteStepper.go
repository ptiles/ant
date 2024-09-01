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
	prevPointPoint, currPoint, prevLine, currLine, currPointColor := f.initialState()
	initialPoint := currPoint.getCenterPoint()
	currentImage := image.NewRGBA(pointRect(initialPoint, 256))

	step := 0
	shouldReset := false
	shouldRun := true

	for shouldRun {
		prevPointPoint, currPoint, prevLine, currLine, currPointColor = f.next(prevPointPoint, currPoint, prevLine, currLine)

		point := currPoint.getCenterPoint()
		if !shouldReset && isOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(pointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[currPointColor])

		select {
		case command := <-commandCh:
			if command == Reset {
				shouldReset = true
			}
		default:
		}

		if shouldReset && currPointColor == 0 {
			shouldReset = false

			currentImage = image.NewRGBA(pointRect(point, 256))

			ResetValues()
		}

		step += 1
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
