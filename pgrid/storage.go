package pgrid

import (
	"os"
	"strconv"
)

type upInt OffsetInt // top 24 bits significant
type gridCoordsUp struct {
	Offset0 upInt
	Offset1 upInt
}

type downInt uint8 // 8 bits
type gridCoordsDown struct {
	Offset0 downInt
	Offset1 downInt
}

type downMap map[gridCoordsDown]uint8

type upArray struct {
	Maps     []downMap
	Stride   upInt
	Min, Max gridCoordsUp
}

var aValues [GridLinesTotal][GridLinesTotal]upArray

const snapBits = 8

func floorSnap(v upInt) upInt {
	return (v >> snapBits) << snapBits
}

func ceilSnap(v upInt) upInt {
	return (v>>snapBits + 1) << snapBits
}

func newUpArray(Min, Max gridCoordsUp) upArray {
	size := gridCoordsUp{Max.Offset0 - Min.Offset0, Max.Offset1 - Min.Offset1}

	maps := make([]downMap, size.Offset0*size.Offset1)

	return upArray{
		Maps:   maps,
		Stride: size.Offset0,
		Min:    Min,
		Max:    Max,
	}
}

const pad = upInt(127)

func (ua *upArray) Initialize(p gridCoordsUp) {
	ua.Min = gridCoordsUp{Offset0: floorSnap(p.Offset0 - pad), Offset1: floorSnap(p.Offset1 - pad)}
	ua.Max = gridCoordsUp{Offset0: ceilSnap(p.Offset0 + pad), Offset1: ceilSnap(p.Offset1 + pad)}

	sizeOffset0 := ua.Max.Offset0 - ua.Min.Offset0
	sizeOffset1 := ua.Max.Offset1 - ua.Min.Offset1
	ua.Stride = sizeOffset0
	ua.Maps = make([]downMap, sizeOffset0*sizeOffset1)
}

func (ua *upArray) ResizeIfNeeded(p gridCoordsUp) {
	if p.Offset0 < ua.Min.Offset0 || p.Offset0 >= ua.Max.Offset0 || p.Offset1 < ua.Min.Offset1 || p.Offset1 >= ua.Max.Offset1 {
		newMin := gridCoordsUp{Offset0: floorSnap(min(ua.Min.Offset0, p.Offset0-pad)), Offset1: floorSnap(min(ua.Min.Offset1, p.Offset1-pad))}
		newMax := gridCoordsUp{Offset0: ceilSnap(max(ua.Max.Offset0-1, p.Offset0+pad)), Offset1: ceilSnap(max(ua.Max.Offset1-1, p.Offset1+pad))}

		ua.Grow(newMin, newMax)
	}
}

func (ua *upArray) Grow(newMin, newMax gridCoordsUp) {
	newSize0 := newMax.Offset0 - newMin.Offset0
	newSize1 := newMax.Offset1 - newMin.Offset1
	newMaps := make([]downMap, newSize0*newSize1)

	rectOffset0 := ua.Min.Offset0 - newMin.Offset0
	rectOffset1 := ua.Min.Offset1 - newMin.Offset1
	newStride := newSize0

	oldIndex := upInt(0)
	newIndex := rectOffset1*newStride + rectOffset0
	for range ua.Max.Offset1 - ua.Min.Offset1 {
		copy(
			newMaps[newIndex:newIndex+ua.Stride],
			ua.Maps[oldIndex:oldIndex+ua.Stride],
		)
		oldIndex += ua.Stride
		newIndex += newStride
	}

	ua.Maps, ua.Stride, ua.Min, ua.Max = newMaps, newStride, newMin, newMax
}

func (ua *upArray) Get(p GridCoords) (downMap, gridCoordsDown) {
	up := gridCoordsUp{Offset0: upInt(p.Offset0 >> bits), Offset1: upInt(p.Offset1 >> bits)}
	down := gridCoordsDown{Offset0: downInt(p.Offset0) & downMask, Offset1: downInt(p.Offset1) & downMask}
	if ua.Maps == nil {
		ua.Initialize(up)
	} else {
		ua.ResizeIfNeeded(up)
	}
	off0 := up.Offset0 - ua.Min.Offset0
	off1 := up.Offset1 - ua.Min.Offset1
	i := off1*ua.Stride + off0

	if ua.Maps[i] == nil {
		ua.Maps[i] = make(downMap, initialMapSize)
	}
	return ua.Maps[i], down
}

var bits = 8
var downMask = downInt(0b11111111)
var initialMapSize = 512

func ResetValues() {
	for ax0, ax1 := range AxesCanon() {
		aValues[ax0][ax1] = upArray{}
	}
}

func init() {
	ResetValues()

	bitsStr, bitsPresent := os.LookupEnv("BITS")
	if bitsPresent {
		bits, _ = strconv.Atoi(bitsStr)
		downMask = 0b11111111 >> (8 - bits)
	}

	initialMapSizeStr, initialMapSizePresent := os.LookupEnv("INITIAL_MAP_SIZE")
	if initialMapSizePresent {
		initialMapSize, _ = strconv.Atoi(initialMapSizeStr)
	}
}

func divUp(c GridCoords) gridCoordsUp {
	return gridCoordsUp{
		Offset0: upInt(c.Offset0 >> bits),
		Offset1: upInt(c.Offset1 >> bits),
	}
}

func Get(axes GridAxes) uint8 {
	val, down := aValues[axes.Axis0][axes.Axis1].Get(axes.Coords)
	return val[down]
}

func Set(axes GridAxes, value uint8) {
	val, down := aValues[axes.Axis0][axes.Axis1].Get(axes.Coords)
	val[down] = value
}

// StepColor is used by RunAxesColor in cmd/ant
func StepColor(axes GridAxes, limit uint8) (uint8, uint8) {
	val, down := aValues[axes.Axis0][axes.Axis1].Get(axes.Coords)
	rule := val[down]
	color := rule + 1
	if color == limit {
		delete(val, down) // no need to store 0 in map
		return rule, 0
	}
	val[down] = color
	return rule, color
}

func StepColor0(axes GridAxes, limit uint8) (uint8, uint8) {
	val, down := aValues[axes.Axis0][axes.Axis1].Get(axes.Coords)
	rule := val[down]
	color := rule + 1
	if color == limit {
		color = 0
	}
	val[down] = color
	return rule, color
}

// Step is used by RunAxes in cmd/ant-dry
func Step(axes GridAxes, limit uint8) uint8 {
	val, down := aValues[axes.Axis0][axes.Axis1].Get(axes.Coords)
	rule := val[down]
	color := rule + 1
	if color == limit {
		delete(val, down) // no need to store 0 in map
		return rule
	}
	val[down] = color
	return rule
}

func Step0(axes GridAxes, limit uint8) uint8 {
	val, down := aValues[axes.Axis0][axes.Axis1].Get(axes.Coords)
	rule := val[down]
	color := rule + 1
	if color == limit {
		color = 0
	}
	val[down] = color
	return rule
}
