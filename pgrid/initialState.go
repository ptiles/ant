package pgrid

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"slices"
	"strconv"
)

func (f *Field) InitialState() (GridLine, GridLine, bool) {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := ParseInitialPoint(f.InitialPoint)

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

func ParseInitialPoint(initialPoint string) (int, int, bool, int, int) {
	re := regexp.MustCompile(`([A-X])(-?\d+)([+-]?)([A-X])(-?\d+)`)
	result := re.FindStringSubmatch(initialPoint)

	currAx, currOff, dir, prevAx, prevOff := result[1], result[2], result[3], result[4], result[5]

	currAxis := slices.Index(AxisNames, currAx)
	currOffset, _ := strconv.Atoi(currOff)

	currAxIncreasing := dir != "-"

	prevAxis := slices.Index(AxisNames, prevAx)
	prevOffset, _ := strconv.Atoi(prevOff)

	return currAxis, currOffset, currAxIncreasing, prevAxis, prevOffset
}

const seedDropBits = 8

func InitialPointSeed(initialPoint string) *rand.Rand {
	currAxis, currOffset, prevPointSign, prevAxis, prevOffset := ParseInitialPoint(initialPoint)

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

	var seed [32]byte
	copy(seed[:], seedString)

	return rand.New(rand.NewChaCha8(seed))
}
