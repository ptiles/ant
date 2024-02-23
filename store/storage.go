package store

const A = 0
const B = 1
const C = 2
const D = 3
const E = 4

// _  AB   AC   AD   AE  |  BC    BD    BE  |  CD    CE  |  DE
// _ a[0]                | a[1]             | a[2]       | else
// _ a[1] a[2] a[3] a[4] | a[2]  a[3]  a[4] | a[3]  a[4] | a[4]
// _   0    1    2    3  |   4     5     6  |   7     8  |   9

func packAxes(axis0, axis1 uint8) uint8 {
	if axis0 == A {
		return axis1 - 1
	}
	if axis0 == B {
		return axis1 + 2
	}
	if axis0 == C {
		return axis1 + 4
	}
	return 9
}

type PackedCoordinates struct {
	PackedAxes uint8
	Offset0    int16
	Offset1    int16
}

func PackCoordinates(axis0, axis1 uint8, offset0, offset1 int16) PackedCoordinates {
	if axis0 < axis1 {
		return PackedCoordinates{packAxes(axis0, axis1), offset0, offset1}
	} else {
		return PackedCoordinates{packAxes(axis1, axis0), offset1, offset0}
	}
}

var axis0ByIndex = [10]uint8{A, A, A, A, B, B, B, C, C, D}
var axis1ByIndex = [10]uint8{B, C, D, E, C, D, E, D, E, E}

func UnpackAxes(packedAxes uint8) (uint8, uint8) {
	return axis0ByIndex[packedAxes], axis1ByIndex[packedAxes]
}

type axisValues map[uint32]uint8

var values [10]axisValues

func init() {
	for pa := 0; pa < 10; pa++ {
		values[pa] = axisValues{}
	}
}

func Get(coords PackedCoordinates) uint8 {
	packedOffsets := uint32(uint16(coords.Offset1))<<16 + uint32(uint16(coords.Offset0))
	result := values[coords.PackedAxes][packedOffsets]
	return result
}

func Set(coords PackedCoordinates, value uint8) {
	offset0, offset1, packedAxes := coords.Offset0, coords.Offset1, coords.PackedAxes
	packedOffsets := uint32(uint16(offset1))<<16 + uint32(uint16(offset0))
	values[packedAxes][packedOffsets] = value
}

func ForEach(callback func(axis0, axis1 uint8, off0, off1 int16, color uint8)) {
	for packedAxes := uint8(0); packedAxes < 10; packedAxes++ {
		for packedOffsets, color := range values[packedAxes] {
			if color > 0 {
				axis0, axis1 := UnpackAxes(packedAxes)
				off1, off0 := int16(uint16(packedOffsets>>16)), int16(uint16(packedOffsets))
				callback(axis0, axis1, off0, off1, color)
			}
		}
	}
}
