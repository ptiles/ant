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
	for _, delta := range f.geometry[ga.Axis0][ga.Axis1].deltas {
		result[delta.targetAx] = OffsetInt(math.Ceil(
			delta.zeroZero + off0*delta.ax0Delta + off1*delta.ax1Delta))
	}

	return result
}
