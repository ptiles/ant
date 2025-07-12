package seq

import (
	"iter"
	"math"
)

func logPhi(a float64) float64 {
	return math.Log(a) / math.Log(math.Phi)
}

func phiPow(a float64) float64 {
	return math.Pow(math.Phi, a)
}

func prevFibIndex(f int) int {
	i := math.Floor(logPhi((float64(f) + 0.5) * math.Sqrt(5)))
	return int(i)
}

func nextFibIndex(f int) int {
	return prevFibIndex(f) + 1
}

func fib(n int) int {
	f := math.Round(phiPow(float64(n)) / math.Sqrt(5))
	return int(f)
}

func wyt(r, c int) int {
	return fib(c+2)*int(math.Floor(float64(r+1)*math.Phi)) + fib(c+1)*(r)
}

func fibAround(a, b int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		if a < 3 {
			a = 3
		}
		aFib := prevFibIndex(a)
		bFib := nextFibIndex(b)
		for n := aFib; n <= bFib; n += 1 {
			if !yield(n, fib(n)) {
				return
			}
		}
	}
}

func fibIter2() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		i, a, b := 0, 1, 1
		if !yield(i, a) {
			return
		}
		for {
			i += 1
			if !yield(i, b) {
				return
			}
			a, b = b, a+b
		}
	}
}

func wythoffShifted(n, minDelta int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		c := n - 4
		if minDelta < 1 {
			minDelta = 1
		}
		r := 1
		for depth, f := range fibIter2() {
			if fib(c-depth+1) < minDelta {
				return
			}
			for range f {
				if !yield(wyt(r, c-depth), depth+1) {
					return
				}
				r += 1
			}
		}

	}
}

func WythoffDelta(a, b, minDelta int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for n, f := range fibAround(a, b) {
			if !yield(f) {
				return
			}
			if f >= b {
				return
			}
			for value := range wythoffShifted(n, minDelta) {
				if !yield(value) {
					return
				}
			}
		}
	}
}

func WythoffDelta2(a, b, minDelta int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for n, f := range fibAround(a, b) {
			if !yield(f, 0) {
				return
			}
			if f >= b {
				return
			}
			for value, depth := range wythoffShifted(n, minDelta) {
				if !yield(value, depth) {
					return
				}
			}
		}
	}
}
