package pgrid

var axesRotation = [GRID_LINES_TOTAL][GRID_LINES_TOTAL]bool{}

func init() {
	for ax0 := range GRID_LINES_TOTAL {
		for i := range GRID_LINES_TOTAL {
			ax1 := (ax0 + i) % GRID_LINES_TOTAL
			// {true, true, true, false, false}
			// {false, true, true, true, false}
			axesRotation[ax0][ax1] = i <= GRID_LINES_TOTAL/2
		}
	}
}

func (f *Field) next(prevPointPoint Point, currPoint GridPoint, prevLine, currLine GridLine) (Point, GridPoint, GridLine, GridLine, uint8) {
	isRightTurn, currPointColor := f.step(currPoint.Axes)

	axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
	prevPointSign := distance(f.getLine(currLine), prevPointPoint) < 0
	positiveSide := isRightTurn != axisRotation != prevPointSign

	nextPoint, nextLine := f.nearestNeighbor(currPoint.Offsets, prevLine, currLine, positiveSide)

	return currPoint.Point, nextPoint, currLine, nextLine, currPointColor
}

func (f *Field) step(axes GridAxes) (bool, uint8) {
	value := Get(axes)
	newValue := (value + 1) % f.Limit
	Set(axes, newValue)
	return f.Rules[value], newValue
}
