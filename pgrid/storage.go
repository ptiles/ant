package pgrid

type axisValues map[GridCoords]uint8

var values = [GridLinesTotal * GridLinesTotal]axisValues{}

func init() {
	ResetValues()
}

func Get(axes GridAxes) uint8 {
	return values[axes.Axis0*GridLinesTotal+axes.Axis1][axes.Coords]
}

func Set(axes GridAxes, value uint8) {
	values[axes.Axis0*GridLinesTotal+axes.Axis1][axes.Coords] = value
}

func ResetValues() {
	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			values[ax0*GridLinesTotal+ax1] = axisValues{}
		}
	}
}
