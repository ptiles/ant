package store

import (
	"encoding/binary"
	"fmt"
	"log"
)

const A = 0
const B = 1
const C = 2
const D = 3
const E = 4

var MinOffset0 = int16(0)
var MaxOffset0 = int16(0)

var MinOffset1 = int16(0)
var MaxOffset1 = int16(0)

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

func unpackAxes(packedAxes uint8) (uint8, uint8) {
	return axis0ByIndex[packedAxes], axis1ByIndex[packedAxes]
}

type segmentData [1 << 8][1 << 8][10]uint8

type segment struct {
	offset0Min int16
	offset0Max int16
	offset1Min int16
	offset1Max int16
	values     *segmentData
}

var cs segment
var segments = make([]segment, 0, 4)

func init() {
	//fmt.Println("Allocating new segment")
	values := segmentData{}
	cs = segment{0, 255, 0, 255, &values}
	segments = append(segments, cs)
}

func getCurrentSegment(offset0, offset1 int16) segment {
	if offset0 < cs.offset0Min || offset0 > cs.offset0Max || offset1 < cs.offset1Min || offset1 > cs.offset1Max {
		//fmt.Println("Changing current segment")
		//fmt.Print(".")
		for _, s := range segments {
			if offset0 >= s.offset0Min &&
				offset0 <= s.offset0Max &&
				offset1 >= s.offset1Min &&
				offset1 <= s.offset1Max {
				cs = s
				return cs
			}
		}
		//fmt.Println("Allocating new segment")
		values := segmentData{}
		offset0Min := offset0 >> 8 << 8
		offset0Max := offset0Min | 0xff
		offset1Min := offset1 >> 8 << 8
		offset1Max := offset1Min | 0xff
		cs = segment{offset0Min, offset0Max, offset1Min, offset1Max, &values}
		segments = append(segments, cs)
	}
	return cs
}

func Get(coords PackedCoordinates) uint8 {
	return get(coords.Offset0, coords.Offset1, coords.PackedAxes)
}

func get(offset0, offset1 int16, packedAxes uint8) uint8 {
	//fmt.Printf("Get(%d, %d)\n", offset0, offset1)
	segment := getCurrentSegment(offset0, offset1)
	result := segment.values[uint8(offset0&0xff)][uint8(offset1&0xff)][packedAxes]
	//fmt.Printf(" => %d\n", result)
	return result
}

func Set(coords PackedCoordinates, value uint8) {
	offset0, offset1, packedAxes := coords.Offset0, coords.Offset1, coords.PackedAxes
	setMinMax(offset0, offset1)
	//fmt.Printf("Set(%d, %d)\n", offset0, offset1)
	segment := getCurrentSegment(offset0, offset1)
	//fmt.Printf(" <= %d\n", value)
	segment.values[uint8(offset0&0xff)][uint8(offset1&0xff)][packedAxes] = value
}

func setMinMax(offset0, offset1 int16) {
	if offset0 > MaxOffset0 {
		MaxOffset0 = offset0
	}
	if offset0 < MinOffset0 {
		MinOffset0 = offset0
	}

	if offset1 > MaxOffset1 {
		MaxOffset1 = offset1
	}
	if offset1 < MinOffset1 {
		MinOffset1 = offset1
	}
}

func ForEach(callback func(axis0, axis1 uint8, off0, off1 int16, color uint8)) {
	for off0 := MinOffset0; off0 <= MaxOffset0; off0++ {
		for off1 := MinOffset1; off1 <= MaxOffset1; off1++ {
			for packedAxes := uint8(0); packedAxes < 10; packedAxes++ {
				color := get(off0, off1, packedAxes)
				if color > 0 {
					axis0, axis1 := unpackAxes(packedAxes)
					callback(axis0, axis1, off0, off1, color)
				}
			}
		}
	}
}

type axisValues []uint8

var values [10]axisValues
var bits uint8
var minO int16
var maxO int16
var deltaI int16

func Allocate(b uint8) {
	bits = b
	for a := 0; a < 10; a++ {
		values[a] = make(axisValues, 2<<bits<<bits)
	}
	significantBits := bits - 1
	minO = -(2 << significantBits)
	maxO = 2<<significantBits - 1
	deltaI = 2 << significantBits

	fmt.Printf("Allocated %.fMB\n", float32(binary.Size(values[0]))*10/1024/1024)
}

func Get2(coords PackedCoordinates) uint8 {
	return get2(coords.Offset0, coords.Offset1, coords.PackedAxes)
}

func get2(offset0, offset1 int16, packedAxes uint8) uint8 {
	i := uint16(offset0+deltaI)<<bits + uint16(offset1+deltaI)
	//fmt.Printf("Get2(%d, %d)\n", offset0, offset1)
	result := values[packedAxes][i]
	//fmt.Printf(" => %d\n", result)
	return result
}

func Set2(coords PackedCoordinates, value uint8) {
	setMinMax2(coords.Offset0, coords.Offset1)
	i := uint16(coords.Offset0+deltaI)<<bits + uint16(coords.Offset1+deltaI)
	//fmt.Printf("Set2(%d, %d)\n", offset0, offset1)
	values[coords.PackedAxes][i] = value
	//fmt.Printf(" <= %d\n", value)
}

func ForEach2(callback func(axis0, axis1 uint8, off0, off1 int16, color uint8)) {
	for off0 := MinOffset0; off0 <= MaxOffset0; off0++ {
		for off1 := MinOffset1; off1 <= MaxOffset1; off1++ {
			for packedAxes := uint8(0); packedAxes < 10; packedAxes++ {
				color := get2(off0, off1, packedAxes)
				if color > 0 {
					axis0, axis1 := unpackAxes(packedAxes)
					callback(axis0, axis1, off0, off1, color)
				}
			}
		}
	}
}

func setMinMax2(offset0, offset1 int16) {
	if offset0 < minO || offset1 < minO || offset0 > maxO || offset1 > maxO {
		log.Fatal("Offset out of bounds")
	}
	if offset0 > MaxOffset0 {
		MaxOffset0 = offset0
	}
	if offset0 < MinOffset0 {
		MinOffset0 = offset0
	}

	if offset1 > MaxOffset1 {
		MaxOffset1 = offset1
	}
	if offset1 < MinOffset1 {
		MinOffset1 = offset1
	}
}
