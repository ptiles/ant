package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"iter"
	"math/rand/v2"
	"strings"
)

func (fl *Flags) InitialPoints(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if fl.initialPointCount == 0 {
			return
		}

		rangeMin, rangeMax, _, _ := utils.ParseRangeStr(fl.initialPointRange)

		if fl.initialPointNear != "" {
			ax1, off1, _, ax2, off2 := utils.ParseInitialPoint(fl.initialPointNear)

			min1 := off1 - rangeMax
			max1 := off1 + rangeMax
			min2 := off2 - rangeMax
			max2 := off2 + rangeMax

			if fl.debug {
				debug.WriteString("\nInitialPoints near:")
			}
			for range fl.initialPointCount {
				point := genRandomPointAround(ax1, min1, max1, ax2, min2, max2)
				if !yield(fmt.Sprintf(" -i %s", point)) {
					return
				}
				if fl.debug {
					debug.WriteString(" ")
					debug.WriteString(point)
				}
			}
		} else if fl.kaleidoscope {
			if fl.debug {
				debug.WriteString("\nInitialPoints kaleidoscope:")
			}
			for range fl.initialPointCount {
				points := genRandomPointKaleidoscope(rangeMin, rangeMax)
				for _, point := range points {
					if !yield(fmt.Sprintf(" -i %s", point)) {
						return
					}
					if fl.debug {
						debug.WriteString(" ")
						debug.WriteString(point)
					}
				}
			}
		} else {
			if fl.debug {
				debug.WriteString("\nInitialPoints count:")
			}
			for range fl.initialPointCount {
				point := genRandomPointString(rangeMin, rangeMax)
				if !yield(fmt.Sprintf(" -i %s", point)) {
					return
				}
				if fl.debug {
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

	ax1s := pgrid.AxisNames[ax1]
	ax2s := pgrid.AxisNames[ax2]

	return fmt.Sprintf("%s%d%s%s%d", ax1s, off1, dir, ax2s, off2)
}

func genRandomPointKaleidoscope(min, max int) [GridLinesTotal]string {
	ax1, off1, dir, ax2, off2 := genRandomPoint(min, max)

	var result [GridLinesTotal]string
	for i := range GridLinesTotal {
		ax1s := pgrid.AxisNames[(ax1+i)%GridLinesTotal]
		ax2s := pgrid.AxisNames[(ax2+i)%GridLinesTotal]
		point := fmt.Sprintf("%s%d%s%s%d", ax1s, off1, dir, ax2s, off2)
		result[i] = point
	}
	return result
}

func genRandomPointAround(ax1, min1, max1, ax2, min2, max2 int) string {
	dir := [2]string{"-", "+"}[rand.IntN(2)]

	off1 := rand.IntN(max1+1-min1) + min1
	off2 := rand.IntN(max2+1-min2) + min2

	return fmt.Sprintf("%s%d%s%s%d", pgrid.AxisNames[ax1], off1, dir, pgrid.AxisNames[ax2], off2)
}
