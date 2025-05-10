package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"math/rand/v2"
)

func (f *Field) InitialState() (GridLine, GridLine, bool) {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := utils.ParseInitialPoint(f.InitialPoint)

	currLine := GridLine{Axis: uint8(currAxis), Offset: offsetInt(currOffset)}
	prevLine := GridLine{Axis: uint8(prevAxis), Offset: offsetInt(prevOffset)}

	//fmt.Printf(
	//	"Initial step: %s %s %s %t\n",
	//	currPoint.String(), currLine.String(), prevLine.String(), prevPointSign,
	//)

	return currLine, prevLine, prevPointSign
}

func (f *Field) initialStateString(currLine GridLine, prevLine GridLine, prevPointSign bool) string {
	prevPointSignString := "-"
	if prevPointSign {
		prevPointSignString = "+"
	}

	return fmt.Sprintf("%s%s%s", currLine.String(), prevPointSignString, prevLine.String())
}

func rngFromString(seedString string) *rand.Rand {
	var seed [32]byte
	copy(seed[:], seedString)

	return rand.New(rand.NewChaCha8(seed))
}

func InitialPointSeed(initialPoint string, seedDropBits uint8) *rand.Rand {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := utils.ParseInitialPoint(initialPoint)

	// Same seed for five symmetric points
	//seedString := fmt.Sprintf(
	//	"%d%t%d%d",
	//	(int(GridLinesTotal)+currAxis-prevAxis)%int(GridLinesTotal),
	//	prevPointSign, currOffset>>seedDropBits, prevOffset>>seedDropBits,
	//)

	// Different seeds
	seedString := fmt.Sprintf(
		"%s%s%t%d%d",
		AxisNames[currAxis], AxisNames[prevAxis],
		prevPointSign, currOffset>>seedDropBits, prevOffset>>seedDropBits,
	)

	return rngFromString(seedString)
}
