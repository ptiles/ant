package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/pgrid/axis"
	"github.com/ptiles/ant/utils"
	"iter"
	"math/rand/v2"
	"strconv"
	"strings"
)

type interval struct {
	count int

	rangeMin     int
	rangeMax     int
	rangePresent bool

	kaleidoscope bool
}

func (i *interval) intervalParser() flagParser {
	return func(interval string) error {
		if interval == "" {
			return nil
		}

		rangeMin, rangeMax, err := utils.ParseRangeStr(interval)
		if err != nil {
			return err
		}

		i.rangeMin = rangeMin
		i.rangeMax = rangeMax
		i.rangePresent = true

		return nil
	}
}

func (i *interval) countParser() flagParser {
	return func(countStr string) error {
		if countStr == "" {
			return nil
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return err
		}

		i.count = count

		return nil
	}
}

func (i *interval) kaleidoscopeParser() flagParser {
	return func(_ string) error {
		i.kaleidoscope = true

		return nil
	}
}

func (i *interval) skip() bool {
	return i.count == 0 || !i.rangePresent
}

func (i *interval) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if i.skip() {
			return
		}

		if i.kaleidoscope {
			if debug != nil {
				debug.WriteString("\nInitialPoints kaleidoscope:")
			}
			for range i.count {
				points := genRandomPointKaleidoscope(i.rangeMin, i.rangeMax)
				for _, point := range points {
					if !yield(fmt.Sprintf(" -i %s", point)) {
						return
					}
					if debug != nil {
						debug.WriteString(" ")
						debug.WriteString(point)
					}
				}
			}
		} else {
			if debug != nil {
				debug.WriteString("\nInitialPoints count:")
			}
			for range i.count {
				point := genRandomPointString(i.rangeMin, i.rangeMax)
				if !yield(fmt.Sprintf(" -i %s", point)) {
					return
				}
				if debug != nil {
					debug.WriteString(" ")
					debug.WriteString(point)
				}
			}
		}
	}
}

const GridLinesTotal = uint(pgrid.GridLinesTotal)

func genRandomAxesDirection() (uint, uint, string) {
	ax := rand.Perm(int(GridLinesTotal))
	ax1, ax2 := uint(ax[0]), uint(ax[1])
	dir := [2]string{"-", "+"}[rand.IntN(2)]

	return ax1, ax2, dir
}

func genRandomPoint(min, max int) (uint, int, string, uint, int) {
	ax1, ax2, dir := genRandomAxesDirection()

	off1 := rand.IntN(max+1-min) + min
	off2 := rand.IntN(max+1-min) + min

	if rand.IntN(2) == 0 {
		off1 = -off1
	}
	if rand.IntN(2) == 0 {
		off2 = -off2
	}

	return ax1, off1, dir, ax2, off2
}

func genRandomPointString(min, max int) string {
	ax1, off1, dir, ax2, off2 := genRandomPoint(min, max)

	ax1s := axis.Name[ax1]
	ax2s := axis.Name[ax2]

	return fmt.Sprintf("%s%d%s%s%d", ax1s, off1, dir, ax2s, off2)
}

func genRandomPointKaleidoscope(min, max int) [GridLinesTotal]string {
	ax1, off1, dir, ax2, off2 := genRandomPoint(min, max)

	var result [GridLinesTotal]string
	for i := range GridLinesTotal {
		ax1s := axis.Name[(ax1+i)%GridLinesTotal]
		ax2s := axis.Name[(ax2+i)%GridLinesTotal]
		point := fmt.Sprintf("%s%d%s%s%d", ax1s, off1, dir, ax2s, off2)
		result[i] = point
	}
	return result
}
