package pgrid

import (
	"image"
	"iter"
)

var axesRotation = [GridLinesTotal][GridLinesTotal]bool{}

// TODO: calculate this properly from angles instead of heuristic
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

//func (ga *GridAxes) update(gridLine0, gridLine1 GridLine) {
//	// Keep axes in ascending order
//	if gridLine0.Axis > gridLine1.Axis {
//		ga.Axis0, ga.Coords.Offset0 = gridLine1.Axis, gridLine1.Offset
//		ga.Axis1, ga.Coords.Offset1 = gridLine0.Axis, gridLine0.Offset
//	} else {
//		ga.Axis0, ga.Coords.Offset0 = gridLine0.Axis, gridLine0.Offset
//		ga.Axis1, ga.Coords.Offset1 = gridLine1.Axis, gridLine1.Offset
//	}
//}

// RunAxesColor is used in cmd/ant
func (f *Field) RunAxesColor(maxSteps uint64) iter.Seq2[GridAxes, uint8] {
	return func(yield func(GridAxes, uint8) bool) {
		initialTurn := f.InitialTurn()
		currLine, prevLine, pointSign := initialTurn.CurrLine, initialTurn.PrevLine, initialTurn.sign

		var currAxes GridAxes

		for range maxSteps {
			// Keep axes in ascending order
			//currAxes.update(currLine, prevLine) // TODO: inline if needed
			if currLine.Axis > prevLine.Axis {
				currAxes.Axis0, currAxes.Coords.Offset0 = prevLine.Axis, prevLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = currLine.Axis, currLine.Offset
			} else {
				currAxes.Axis0, currAxes.Coords.Offset0 = currLine.Axis, currLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = prevLine.Axis, prevLine.Offset
			}

			rule, pointColor := StepColor(currAxes, f.Limit)
			isRightTurn := f.Rules[rule]

			if !yield(currAxes, pointColor) {
				return
			}

			axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
			positiveSide := isRightTurn != axisRotation != pointSign

			nextLine, nextPointSign := f.nearestNeighbor(prevLine, currLine, positiveSide)

			currLine, prevLine, pointSign = nextLine, currLine, nextPointSign
		}
	}
}

// RunPoint was used in cmd/ant-dry
func (f *Field) RunPoint(maxSteps uint64) iter.Seq2[GridAxes, image.Point] {
	return func(yield func(GridAxes, image.Point) bool) {
		initialTurn := f.InitialTurn()
		currLine, prevLine, pointSign := initialTurn.CurrLine, initialTurn.PrevLine, initialTurn.sign

		var currAxes GridAxes

		for range maxSteps {
			// Keep axes in ascending order
			//currAxes.update(currLine, prevLine) // TODO: inline if needed
			if currLine.Axis > prevLine.Axis {
				currAxes.Axis0, currAxes.Coords.Offset0 = prevLine.Axis, prevLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = currLine.Axis, currLine.Offset
			} else {
				currAxes.Axis0, currAxes.Coords.Offset0 = currLine.Axis, currLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = prevLine.Axis, prevLine.Offset
			}

			rule := Step(currAxes, f.Limit)
			isRightTurn := f.Rules[rule]

			centerPoint := f.GetCenterPoint(currAxes)
			if !yield(currAxes, centerPoint) {
				return
			}

			axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
			positiveSide := isRightTurn != axisRotation != pointSign

			nextLine, nextPointSign := f.nearestNeighbor(prevLine, currLine, positiveSide)

			currLine, prevLine, pointSign = nextLine, currLine, nextPointSign
		}
	}
}

// RunAxes is used in cmd/ant-dry and cmd/ant-batch/path
func (f *Field) RunAxes(maxSteps uint64) iter.Seq[GridAxes] {
	return func(yield func(GridAxes) bool) {
		initialTurn := f.InitialTurn()
		currLine, prevLine, pointSign := initialTurn.CurrLine, initialTurn.PrevLine, initialTurn.sign

		var currAxes GridAxes

		for range maxSteps {
			// Keep axes in ascending order
			//currAxes.update(currLine, prevLine) // TODO: inline if needed
			if currLine.Axis > prevLine.Axis {
				currAxes.Axis0, currAxes.Coords.Offset0 = prevLine.Axis, prevLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = currLine.Axis, currLine.Offset
			} else {
				currAxes.Axis0, currAxes.Coords.Offset0 = currLine.Axis, currLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = prevLine.Axis, prevLine.Offset
			}

			rule := Step(currAxes, f.Limit)
			isRightTurn := f.Rules[rule]

			if !yield(currAxes) {
				return
			}

			axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
			positiveSide := isRightTurn != axisRotation != pointSign

			nextLine, nextPointSign := f.nearestNeighbor(prevLine, currLine, positiveSide)

			currLine, prevLine, pointSign = nextLine, currLine, nextPointSign
		}
	}
}

// RunTurns is used in cmd/ant-batch -ip
func (f *Field) RunTurns(maxSteps uint64) iter.Seq2[uint64, Turn] {
	return func(yield func(uint64, Turn) bool) {
		initialTurn := f.InitialTurn()
		currLine, prevLine, pointSign := initialTurn.CurrLine, initialTurn.PrevLine, initialTurn.sign

		var currAxes GridAxes

		for i := range maxSteps {
			// Keep axes in ascending order
			//currAxes.update(currLine, prevLine) // TODO: inline if needed
			if currLine.Axis > prevLine.Axis {
				currAxes.Axis0, currAxes.Coords.Offset0 = prevLine.Axis, prevLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = currLine.Axis, currLine.Offset
			} else {
				currAxes.Axis0, currAxes.Coords.Offset0 = currLine.Axis, currLine.Offset
				currAxes.Axis1, currAxes.Coords.Offset1 = prevLine.Axis, prevLine.Offset
			}

			rule := Step(currAxes, f.Limit)
			isRightTurn := f.Rules[rule]

			if !yield(i, Turn{CurrLine: currLine, PrevLine: prevLine, sign: pointSign}) {
				return
			}

			axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
			positiveSide := isRightTurn != axisRotation != pointSign

			nextLine, nextPointSign := f.nearestNeighbor(prevLine, currLine, positiveSide)

			currLine, prevLine, pointSign = nextLine, currLine, nextPointSign
		}
	}
}
