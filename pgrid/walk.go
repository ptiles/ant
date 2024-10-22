package pgrid

import (
	"image"
	"iter"
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

func (f *Field) step0(axes GridAxes) bool {
	value := Get(axes)
	Set0(axes, (value+1)%f.Limit)
	return f.Rules[value]
}

func (f *Field) Run(maxSteps uint64) iter.Seq2[GridPoint, uint8] {
	return func(yield func(GridPoint, uint8) bool) {
		currPoint, currLine, prevLine, pointSign := f.initialState()

		for range maxSteps {
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

func (f *Field) RunAxes(maxSteps uint64) iter.Seq[GridAxes] {
	return func(yield func(GridAxes) bool) {
		currPoint, currLine, prevLine, pointSign := f.initialState()

		for range maxSteps {
			isRightTurn := f.step0(currPoint.Axes)

			if !yield(currPoint.Axes) {
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

func (f *Field) RunPoint(maxSteps uint64) iter.Seq2[GridAxes, image.Point] {
	return func(yield func(GridAxes, image.Point) bool) {
		currPoint, currLine, prevLine, pointSign := f.initialState()

		for range maxSteps {
			isRightTurn := f.step0(currPoint.Axes)

			if !yield(currPoint.Axes, currPoint.getCenterPoint()) {
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
