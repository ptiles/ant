package pgrid

import (
	"github.com/emirpasic/gods/v2/trees/btree"
	"os"
	"strconv"
)

type tAxisValues map[gridCoordsDown]uint8

var tValues = [GridLinesTotal * GridLinesTotal]*btree.Tree[gridCoordsUp, tAxisValues]{}

type upInt uint32
type gridCoordsUp struct {
	Offset0 upInt
	Offset1 upInt
}

type downInt uint16
type gridCoordsDown struct {
	Offset0 downInt
	Offset1 downInt
}

var bits = 8
var upMask = upInt(0b11111111_11111111_11111111_00000000)
var downMask = downInt(0b00000000_11111111)
var order = 4

func init() {
	ResetValues()

	bitsStr, bitsPresent := os.LookupEnv("BITS")
	if bitsPresent {
		bits, _ = strconv.Atoi(bitsStr)
		upMask = 0b11111111_11111111_11111111_11111111 << bits
		downMask = 0b11111111_11111111 >> (16 - bits)
	}

	orderStr, orderPresent := os.LookupEnv("ORDER")
	if orderPresent {
		order, _ = strconv.Atoi(orderStr)
	}
}

func DivCoords(c GridCoords) (up gridCoordsUp, down gridCoordsDown) {
	up = gridCoordsUp{
		Offset0: upInt(c.Offset0) & upMask,
		Offset1: upInt(c.Offset1) & upMask,
	}
	down = gridCoordsDown{
		Offset0: downInt(c.Offset0) & downMask,
		Offset1: downInt(c.Offset1) & downMask,
	}
	return
}

func Get(axes GridAxes) uint8 {
	up, down := DivCoords(axes.Coords)
	val, _ := tValues[axes.Axis0*GridLinesTotal+axes.Axis1].Get(up)
	return val[down]
}

func Set(axes GridAxes, value uint8) {
	ax := axes.Axis0*GridLinesTotal + axes.Axis1
	up, down := DivCoords(axes.Coords)
	val, found := tValues[ax].Get(up)
	if !found {
		val = tAxisValues{}
		tValues[ax].Put(up, val)
	}
	val[down] = value
}

func Inc(axes GridAxes, limit uint8) (uint8, uint8) {
	ax := axes.Axis0*GridLinesTotal + axes.Axis1
	up, down := DivCoords(axes.Coords)
	val, found := tValues[ax].Get(up)
	if !found {
		tValues[ax].Put(up, tAxisValues{down: 1})
		return 0, 1
	}
	value := val[down]
	newValue := (value + 1) % limit
	val[down] = newValue
	return value, newValue
}

func Set0(axes GridAxes, value uint8) {
	ax := axes.Axis0*GridLinesTotal + axes.Axis1
	up, down := DivCoords(axes.Coords)
	val, found := tValues[ax].Get(up)
	if !found {
		tValues[ax].Put(up, tAxisValues{down: 1})
		return
	}
	if value == 0 {
		delete(val, down)
	} else {
		val[down] = value
	}
}

func gridCoordsUpCmp(a, b gridCoordsUp) int {
	switch {
	case a.Offset0 > b.Offset0:
		return 1
	case a.Offset0 < b.Offset0:
		return -1
	case a.Offset1 > b.Offset1:
		return 1
	case a.Offset1 < b.Offset1:
		return -1
	default:
		return 0
	}
}
func ResetValues() {
	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			tValues[ax0*GridLinesTotal+ax1] = btree.NewWith[gridCoordsUp, tAxisValues](order, gridCoordsUpCmp)
		}
	}
}
