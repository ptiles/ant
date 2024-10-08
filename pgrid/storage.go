package pgrid

type axisValues map[GridAxes]uint8

var values = axisValues{}

func Get(axes GridAxes) uint8 {
	return values[axes]
}

func Set(axes GridAxes, value uint8) {
	values[axes] = value
}

func ResetValues() {
	values = axisValues{}
}
