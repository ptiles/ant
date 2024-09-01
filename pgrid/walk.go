package pgrid

var axesRotation = [5][5]bool{
	{true, true, true, false, false},
	{false, true, true, true, false},
	{false, false, true, true, true},
	{true, false, false, true, true},
	{true, true, false, false, true},
}

func (f *Field) nextPoint(prevPoint, currPoint GridPoint, prevLine, currLine GridLine, isRightTurn bool) (GridPoint, GridPoint, GridLine, GridLine) {
	axisRotation := axesRotation[prevLine.Axis][currLine.Axis]
	prevPointSign := distance(f.getLine(currLine), prevPoint.Point) < 0
	positiveSide := isRightTurn != axisRotation != prevPointSign

	nextPoint, nextLine := f.nearestNeighbor(currPoint, prevLine, currLine, positiveSide)
	return currPoint, nextPoint, currLine, nextLine
}

func (f *Field) walk(axes GridAxes) (bool, uint8) {
	value := Get(axes)
	newValue := (value + 1) % f.Limit
	Set(axes, newValue)
	return f.Rules[value], newValue
}

func (f *Field) step(prevPoint, currPoint GridPoint, prevLine, currLine GridLine) (GridPoint, GridPoint, GridLine, GridLine, uint8) {
	isRightTurn, prevPointColor := f.walk(currPoint.Axes)
	prevPoint, currPoint, prevLine, currLine = f.nextPoint(prevPoint, currPoint, prevLine, currLine, isRightTurn)
	return prevPoint, currPoint, prevLine, currLine, prevPointColor
}
