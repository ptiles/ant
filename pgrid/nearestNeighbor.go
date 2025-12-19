package pgrid

import (
	"math"
)

func (ig *intersectionGeometry) nearestNeighbor(
	prevLineOffset, currentLineOffset float64, positiveSide bool,
) (nextLine GridLine, nextPointSign bool) {
	currentDistance := 1.0

	for i := range ig {
		nextLineOffset := ig[i].zeroZero +
			ig[i].ax0Delta*prevLineOffset +
			ig[i].ax1Delta*currentLineOffset

		nextLineOffsetRounded := math.Floor(nextLineOffset)

		if ig[i].ceilSide == positiveSide {
			nextLineOffsetRounded += 1
		}

		absDist := math.Abs((nextLineOffset - nextLineOffsetRounded) * ig[i].distDelta)

		if absDist < currentDistance {
			currentDistance = absDist
			nextLine = GridLine{ig[i].targetAx, OffsetInt(nextLineOffsetRounded)}
			nextPointSign = ig[i].ceilSide == positiveSide
		}
	}

	return nextLine, nextPointSign
}
