package pgrid

import (
	"github.com/ptiles/ant/utils"
	"image"
	"image/color"
)

type CommandType int

const (
	Reset CommandType = iota
)

func (f *Field) ControlledInfiniteStepper(modifiedImagesCh chan<- *image.RGBA, commandCh <-chan CommandType, palette []color.RGBA) {
	currPoint, currLine, prevLine, prevPointSign, pointColor := f.initialState()
	initialPoint := currPoint.getCenterPoint()
	currentImage := image.NewRGBA(utils.PointRect(initialPoint, 256))

	step := 0
	shouldReset := false
	shouldRun := true

	for shouldRun {
		currPoint, currLine, prevLine, prevPointSign, pointColor = f.next(currPoint, currLine, prevLine, prevPointSign)

		point := currPoint.getCenterPoint()
		if !shouldReset && utils.IsOutside(point, currentImage.Rect) {
			modifiedImagesCh <- currentImage
			currentImage = image.NewRGBA(utils.PointRect(point, 256))
		}
		currentImage.Set(point.X, point.Y, palette[pointColor])

		select {
		case command := <-commandCh:
			if command == Reset {
				shouldReset = true
			}
		default:
		}

		if shouldReset && pointColor == 0 {
			shouldReset = false

			currentImage = image.NewRGBA(utils.PointRect(point, 256))

			ResetValues()
		}

		step += 1
	}
	modifiedImagesCh <- currentImage
	close(modifiedImagesCh)
}
