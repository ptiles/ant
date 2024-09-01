package pgrid

var axesRotation = [5][5]bool{
	{true, true, true, false, false},
	{false, true, true, true, false},
	{false, false, true, true, true},
	{true, false, false, true, true},
	{true, true, false, false, true},
}

func (f *Field) next(prevPoint, currPoint GridPoint, prevLine, currLine GridLine) (GridPoint, GridPoint, GridLine, GridLine, uint8) {
	isRightTurn, prevPointColor := f.step(currPoint.Axes)

	axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
	prevPointSign := distance(f.getLine(currLine), prevPoint.Point) < 0
	positiveSide := isRightTurn != axisRotation != prevPointSign

	nextPoint, nextLine := f.nearestNeighbor(currPoint, prevLine, currLine, positiveSide)

	return currPoint, nextPoint, currLine, nextLine, prevPointColor
}

func (f *Field) step(axes GridAxes) (bool, uint8) {
	value := Get(axes)
	newValue := (value + 1) % f.Limit
	Set(axes, newValue)
	return f.Rules[value], newValue
}
