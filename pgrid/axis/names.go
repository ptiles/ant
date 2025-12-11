package axis

import "slices"

var Name = []string{
	"A", "B", "C", "D", "E",
	"F", "G", "H", "I", "J",
	"K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T",
	"U", "V", "W", "X", "Y",
}

func Index(ax string) uint8 {
	return uint8(slices.Index(Name, ax))
}
