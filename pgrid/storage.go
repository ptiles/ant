package pgrid

import (
	"os"
	"strconv"
)

type upArray struct {
	Maps     []downMap
	Stride   offsetInt
	Min, Max GridCoords
}

const snapBits = 8

func floorSnap(v offsetInt) offsetInt {
	return (v >> snapBits) << snapBits
}

func ceilSnap(v offsetInt) offsetInt {
	return (v>>snapBits + 1) << snapBits
}

func newUpArray(Min, Max GridCoords) upArray {
	size := GridCoords{Max.Offset0 - Min.Offset0, Max.Offset1 - Min.Offset1}

	maps := make([]downMap, size.Offset0*size.Offset1)

	return upArray{
		Maps:   maps,
		Stride: size.Offset0,
		Min:    Min,
		Max:    Max,
	}
}

func (ua *upArray) Initialize(p GridCoords) {
	pad := offsetInt(127)
	ua.Min = GridCoords{Offset0: floorSnap(p.Offset0 - pad), Offset1: floorSnap(p.Offset1 - pad)}
	ua.Max = GridCoords{Offset0: ceilSnap(p.Offset0 + pad), Offset1: ceilSnap(p.Offset1 + pad)}

	sizeOffset0 := ua.Max.Offset0 - ua.Min.Offset0
	sizeOffset1 := ua.Max.Offset1 - ua.Min.Offset1
	ua.Stride = sizeOffset0
	ua.Maps = make([]downMap, sizeOffset0*sizeOffset1)
}

func (ua *upArray) ResizeIfNeeded(p GridCoords) {
	if p.Offset0 >= ua.Min.Offset0 && p.Offset0 < ua.Max.Offset0 && p.Offset1 >= ua.Min.Offset1 && p.Offset1 < ua.Max.Offset1 {
		return
	}

	pad := offsetInt(127)
	newMin := GridCoords{Offset0: floorSnap(min(ua.Min.Offset0, p.Offset0-pad)), Offset1: floorSnap(min(ua.Min.Offset1, p.Offset1-pad))}
	newMax := GridCoords{Offset0: ceilSnap(max(ua.Max.Offset0-1, p.Offset0+pad)), Offset1: ceilSnap(max(ua.Max.Offset1-1, p.Offset1+pad))}

	ua.Maps, ua.Stride = ua.Copy(newMin, newMax)
	ua.Min = newMin
	ua.Max = newMax
}

func (ua *upArray) Copy(newMin, newMax GridCoords) ([]downMap, offsetInt) {
	newSize0 := newMax.Offset0 - newMin.Offset0
	newSize1 := newMax.Offset1 - newMin.Offset1
	newMaps := make([]downMap, newSize0*newSize1)

	rectOffset0 := ua.Min.Offset0 - newMin.Offset0
	rectOffset1 := ua.Min.Offset1 - newMin.Offset1
	newStride := newSize0

	oldIndex := offsetInt(0)
	newIndex := rectOffset1*newStride + rectOffset0
	for range ua.Max.Offset1 - ua.Min.Offset1 {
		copy(
			newMaps[newIndex:newIndex+ua.Stride],
			ua.Maps[oldIndex:oldIndex+ua.Stride],
		)
		oldIndex += ua.Stride
		newIndex += newStride
	}

	return newMaps, newStride
}

func (ua *upArray) Get(p GridCoords) downMap {
	if ua.Maps == nil {
		ua.Initialize(p)
	} else {
		ua.ResizeIfNeeded(p)
	}
	off0 := p.Offset0 - ua.Min.Offset0
	off1 := p.Offset1 - ua.Min.Offset1
	i := off1*ua.Stride + off0

	if ua.Maps[i] == nil {
		ua.Maps[i] = make(downMap, 3*1024)
	}
	return ua.Maps[i]
}

type downMap map[gridCoordsDown]uint8

var aValues [GridLinesTotal][GridLinesTotal]upArray

type downInt uint16
type gridCoordsDown struct {
	Offset0 downInt
	Offset1 downInt
}

var bits = 8
var downMask = downInt(0b00000000_11111111)

func init() {
	ResetValues()

	bitsStr, bitsPresent := os.LookupEnv("BITS")
	if bitsPresent {
		bits, _ = strconv.Atoi(bitsStr)
		downMask = 0b11111111_11111111 >> (16 - bits)
	}
}

func (gc *GridCoords) equals(oth GridCoords) bool {
	return gc.Offset0 == oth.Offset0 && gc.Offset1 == oth.Offset1
}

func DivCoords(c GridCoords) (GridCoords, gridCoordsDown) {
	return GridCoords{
			Offset0: c.Offset0 >> bits,
			Offset1: c.Offset1 >> bits,
		}, gridCoordsDown{
			Offset0: downInt(c.Offset0) & downMask,
			Offset1: downInt(c.Offset1) & downMask,
		}
}

func Get(axes GridAxes) uint8 {
	up, down := DivCoords(axes.Coords)
	val := aValues[axes.Axis0][axes.Axis1].Get(up)
	return val[down]
}

func Set(axes GridAxes, value uint8) {
	up, down := DivCoords(axes.Coords)
	val := aValues[axes.Axis0][axes.Axis1].Get(up)
	val[down] = value
}

func Inc(axes GridAxes, limit uint8) (uint8, uint8) {
	up, down := DivCoords(axes.Coords)
	val := aValues[axes.Axis0][axes.Axis1].Get(up)
	value := val[down]
	newValue := (value + 1) % limit
	val[down] = newValue
	return value, newValue
}

func Set0(axes GridAxes, value uint8) {
	up, down := DivCoords(axes.Coords)
	val := aValues[axes.Axis0][axes.Axis1].Get(up)
	if value == 0 {
		delete(val, down)
	} else {
		val[down] = value
	}
}

func ResetValues() {
	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			aValues[ax0][ax1] = upArray{}
		}
	}
}

func Uniq() (uint64, int) {
	uPoints := uint64(0)
	uMaps := 0

	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			uArr := aValues[ax0][ax1]
			for _, dMap := range uArr.Maps {
				if dMap != nil {
					uPoints += uint64(len(dMap))
					uMaps += 1
				}
			}
		}
	}

	return uPoints, uMaps
}
