package seq

import (
	"cmp"
	"iter"
	"maps"
	"math"
	"slices"
)

type RowCol struct{ Row, Col int }
type NumCol struct{ Num, Col int }
type NumColSlice []NumCol

// recalculate capacities on rowsMax or colsMax change

var Wythoff = make(map[RowCol]int, 16450)
var WythoffReverse = make(map[int]RowCol, 32901)
var WythoffReverseSorted NumColSlice

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

	WythoffReverseSorted = make([]NumCol, len(WythoffReverse))
	for i, num := range slices.Sorted(maps.Keys(WythoffReverse)) {
		WythoffReverseSorted[i] = NumCol{
			Num: num,
			Col: WythoffReverse[num].Col,
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
			if col >= minColumn {
				if !yield(off) {
					return
				}
			}
		}
	}
}

func compare(e NumCol, t int) int {
	return cmp.Compare(e.Num, t)
}

func (ncs NumColSlice) Slice(a, b int) NumColSlice {
	aIndex, _ := slices.BinarySearchFunc(WythoffReverseSorted, a, compare)
	bIndex, _ := slices.BinarySearchFunc(WythoffReverseSorted, b, compare)

	return WythoffReverseSorted[aIndex:bIndex]
}

func (ncs NumColSlice) MinMaxColumn(minColumn, maxColumn int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, numCol := range ncs {
			col := numCol.Col
			if col >= minColumn && col < maxColumn {
				if !yield(numCol.Num) {
					return
				}
			}
		}
	}
}
