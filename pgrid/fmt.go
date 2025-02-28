package pgrid

import "fmt"

var AxisNames = []string{
	"A", "B", "C", "D", "E",
	"F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y",
}

func (gl *GridLine) String() string {
	return fmt.Sprintf("%s%d", AxisNames[gl.Axis], gl.Offset)
}

func (gl *GridLine) Print() {
	fmt.Println(gl)
}

func (gp *GridPoint) String() string {
	offsets := gp.Offsets
	ax0, ax1 := gp.Axes.Axis0, gp.Axes.Axis1
	return fmt.Sprintf(
		"[A:%d, B:%d, C:%d, D:%d, E:%d] %s%d:%s%d",
		offsets[0], offsets[1], offsets[2], offsets[3], offsets[4],
		AxisNames[ax0], offsets[ax0], AxisNames[ax1], offsets[ax1],
	)
}

func (ga *GridAxes) String() string {
	return fmt.Sprintf(
		"%s%d:%s%d",
		AxisNames[ga.Axis0], ga.Coords.Offset0, AxisNames[ga.Axis1], ga.Coords.Offset1,
	)
}
