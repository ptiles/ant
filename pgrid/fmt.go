package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/pgrid/axis"
)

func (gl *GridLine) String() string {
	return fmt.Sprintf("%s%d", axis.Name[gl.Axis], gl.Offset)
}

func (gl *GridLine) Print() {
	fmt.Println(gl)
}

func (gp *GridPoint) String() string {
	offsets := gp.Offsets
	ax0, ax1 := gp.Axes.Axis0, gp.Axes.Axis1
	return fmt.Sprintf(
		"[A%d B%d C%d D%d E%d] %s%d:%s%d",
		offsets[0], offsets[1], offsets[2], offsets[3], offsets[4],
		axis.Name[ax0], offsets[ax0], axis.Name[ax1], offsets[ax1],
	)
}

func (gc *GridCoords) String() string {
	return fmt.Sprintf("(%d,%d)", gc.Offset0, gc.Offset1)
}

func (gcu *gridCoordsUp) String() string {
	return fmt.Sprintf("(%d,%d)", gcu.Offset0, gcu.Offset1)
}

func (ga *GridAxes) String() string {
	return fmt.Sprintf(
		"%s%d:%s%d",
		axis.Name[ga.Axis0], ga.Coords.Offset0, axis.Name[ga.Axis1], ga.Coords.Offset1,
	)
}
