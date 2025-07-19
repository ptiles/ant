package seq

import (
	"iter"
	"math"
)

type RowCol struct{ Row, Col int }

var Wythoff = make(map[RowCol]int)
var WythoffReverse = make(map[int]RowCol)

const rowsMax = 512
const colsMax = 36

func init() {
	WythoffReverse[0] = RowCol{Row: 0, Col: colsMax}

	for row := 1; row <= rowsMax; row += 1 {
		iPhi := math.Floor(float64(row)*math.Phi) * math.Phi
		prev := int(math.Floor(iPhi))
		curr := int(math.Floor(iPhi * math.Phi))

		Wythoff[RowCol{Row: row, Col: 1}] = prev
		Wythoff[RowCol{Row: row, Col: 2}] = curr

		WythoffReverse[prev] = RowCol{Row: row, Col: 1}
		WythoffReverse[-prev] = RowCol{Row: row, Col: 1}
		WythoffReverse[curr] = RowCol{Row: row, Col: 2}
		WythoffReverse[-curr] = RowCol{Row: row, Col: 2}

		for col := 3; col <= colsMax; col += 1 {
			next := prev + curr
			if next > math.MaxInt32 {
				break
			}

			Wythoff[RowCol{Row: row, Col: col}] = next

			WythoffReverse[next] = RowCol{Row: row, Col: col}
			WythoffReverse[-next] = RowCol{Row: row, Col: col}

			prev, curr = curr, next
		}
	}
}

func WythoffMinColumn(a, b, minColumn int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for off := a; off <= b; off += 1 {
			col := WythoffReverse[off].Col
			if col > minColumn {
				if !yield(col, off) {
					return
				}
			}
		}
	}
}

func logPhi(a float64) float64 {
	return math.Log(a) / math.Log(math.Phi)
}

func prevFibIndex(f int) int {
	return int(math.Floor(logPhi((float64(f) + 0.5) * math.Sqrt(5))))
}

func WythoffDelta(a, b, minDelta int) iter.Seq[int] {
	return func(yield func(int) bool) {
		minColumn := prevFibIndex(minDelta)
		for off := a; off <= b; off += 1 {
			col := WythoffReverse[off].Col
			if col > minColumn {
				if !yield(off) {
					return
				}
			}
		}
	}
}
