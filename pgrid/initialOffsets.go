package pgrid

import (
	"math"
)

func (f *Field) InitialOffsets() GridOffsets {
	return f.getOffsets(GridAxes{
		Axis0: f.currAxis, Axis1: f.prevAxis,
		Coords: GridCoords{
			Offset0: OffsetInt(f.currOffset), Offset1: OffsetInt(f.prevOffset),
		},
	})
}

func (f *Field) getOffsets(ga GridAxes) GridOffsets {
	result := GridOffsets{}

	result[ga.Axis0] = ga.Coords.Offset0
	result[ga.Axis1] = ga.Coords.Offset1

	off0, off1 := float64(ga.Coords.Offset0), float64(ga.Coords.Offset1)
	for _, otl := range f.offsetsToLast[ga.Axis0][ga.Axis1] {
		off := math.Ceil(otl.zeroZero + off0*otl.ax0Delta + off1*otl.ax1Delta)
		result[otl.targetAx] = OffsetInt(off)
	}

	return result
}
