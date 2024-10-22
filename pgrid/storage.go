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

func Set0(axes GridAxes, value uint8) {
	if value == 0 {
		delete(values[axes.Axis0*GridLinesTotal+axes.Axis1], axes.Coords)
	} else {
		values[axes.Axis0*GridLinesTotal+axes.Axis1][axes.Coords] = value
	}
}

func ResetValues() {
	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			values[ax0*GridLinesTotal+ax1] = axisValues{}
		}
	}
}

func ClearZeros() {
	for a := range values {
		for k, v := range values[a] {
			if v == 0 {
				delete(values[a], k)
			}
		}
	}
}

func Uniq() (uniq uint64) {
	for a := range values {
		uniq += uint64(len(values[a]))
	}
	return
}
