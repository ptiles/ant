package pgrid

import (
	"github.com/schollz/progressbar/v3"
	"iter"
	"strings"
)

var axesRotation = [GridLinesTotal][GridLinesTotal]bool{}

func init() {
	for ax0 := range GridLinesTotal {
		for i := range GridLinesTotal {
			ax1 := (ax0 + i) % GridLinesTotal
			// {true, true, true, false, false}
			// {false, true, true, true, false}
			axesRotation[ax0][ax1] = i <= GridLinesTotal/2
		}
	}
}

func (f *Field) next(currPoint GridPoint, currLine, prevLine GridLine, prevPointSign bool) (GridPoint, GridLine, GridLine, bool, uint8) {
	isRightTurn, currPointColor := f.step(currPoint.Axes)

	axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
	positiveSide := isRightTurn != axisRotation != prevPointSign

	nextPoint, nextLine := f.nearestNeighbor(currPoint.Offsets, prevLine, currLine, positiveSide)
	currPointSign := distance(f.getLine(nextLine), currPoint.Point) < 0

	return nextPoint, nextLine, currLine, currPointSign, currPointColor
}

func (f *Field) step(axes GridAxes) (bool, uint8) {
	value := Get(axes)
	newValue := (value + 1) % f.Limit
	Set(axes, newValue)
	return f.Rules[value], newValue
}

func (f *Field) Run(maxSteps uint64, bar *progressbar.ProgressBar) iter.Seq2[GridPoint, uint8] {
	return func(yield func(GridPoint, uint8) bool) {
		initialTotal := 25
		initialCounter := initialTotal
		initialPoints := make([]string, initialTotal)

		currPoint, currLine, prevLine, pointSign, _ := f.initialState()

		for step := range maxSteps {
			if initialCounter > 0 && step%10_000 == 0 {
				initialPoints[initialTotal-initialCounter] = f.initialStateString(currLine, prevLine, pointSign)
				initialCounter -= 1
				if initialCounter == 0 {
					bar.Clear()
					println(strings.Join(initialPoints, ","))
				}
			}

			isRightTurn, pointColor := f.step(currPoint.Axes)

			if !yield(currPoint, pointColor) {
				return
			}

			axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
			positiveSide := isRightTurn != axisRotation != pointSign

			nextPoint, nextLine := f.nearestNeighbor(currPoint.Offsets, prevLine, currLine, positiveSide)
			pointSign = distance(f.getLine(nextLine), currPoint.Point) < 0

			currPoint, currLine, prevLine = nextPoint, nextLine, currLine
		}
	}
}
