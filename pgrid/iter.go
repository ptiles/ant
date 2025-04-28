package pgrid

import "iter"

func AxesCanon() iter.Seq2[uint8, uint8] {
	return func(yield func(uint8, uint8) bool) {
		for dax := range GridLinesTotal / 2 {
			for ax0, ax1 := range AxesPairsSymmetric(dax) {
				if !yield(ax0, ax1) {
					return
				}
			}
		}
	}
}

func AxesPairsSymmetric(dax uint8) iter.Seq2[uint8, uint8] {
	return func(yield func(uint8, uint8) bool) {
		for ax0 := range GridLinesTotal {
			ax1 := ax0 + dax + 1
			if ax1 >= GridLinesTotal {
				ax0, ax1 = ax1%GridLinesTotal, ax0
			}
			if !yield(ax0, ax1) {
				return
			}
		}
	}
}

func AxesAll() iter.Seq2[uint8, uint8] {
	return func(yield func(uint8, uint8) bool) {
		for ax0 := range GridLinesTotal {
			for ax1 := range GridLinesTotal {
				if ax0 == ax1 {
					continue
				}
				if !yield(ax0, ax1) {
					return
				}
			}
		}
	}
}

func otherAxes(ax0, ax1 uint8) iter.Seq2[uint8, uint8] {
	return func(yield func(uint8, uint8) bool) {
		i := uint8(0)
		for ax := range GridLinesTotal {
			if ax == ax0 || ax == ax1 {
				continue
			}
			if !yield(i, ax) {
				return
			}
			i += 1
		}
	}
}
