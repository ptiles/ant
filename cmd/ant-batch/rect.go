package main

import (
	"fmt"
	"image"
	"iter"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/pgrid/axis"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/wgrid"
)

type rect struct {
	rect        image.Rectangle
	scaleFactor int

	count int
}

func (r *rect) rectParser() flagParser {
	return func(rectStr string) error {
		if rectStr == "" {
			return nil
		}

		rectangle, scaleFactor, err := utils.ParseRectangleStr(rectStr)
		if err != nil {
			return err
		}

		r.rect = rectangle
		r.scaleFactor = scaleFactor

		return nil
	}
}

func (r *rect) countParser() flagParser {
	return func(countStr string) error {
		if countStr == "" {
			return nil
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return err
		}

		r.count = count

		return nil
	}
}

func (r *rect) skip() bool {
	return r.rect.Empty() || r.count == 0
}

func (r *rect) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if r.skip() {
			return
		}

		if debug != nil {
			debug.WriteString("\nRect:")
		}

		wg := wgrid.New(r.rect)

		count := 0

		for count < r.count {
			ax0, ax1, dir := genRandomAxesDirection()

			minOffset0, maxOffset0 := wg.Ranges[ax0].Min, wg.Ranges[ax0].Max
			minOffset1, maxOffset1 := wg.Ranges[ax1].Min, wg.Ranges[ax1].Max

			off0 := rand.IntN(maxOffset0+1-minOffset0) + minOffset0
			off1 := rand.IntN(maxOffset1+1-minOffset1) + minOffset1

			gridAxes := pgrid.GridAxes{
				Axis0: uint8(ax0), Axis1: uint8(ax1),
				Coords: pgrid.GridCoords{
					Offset0: pgrid.OffsetInt(off0), Offset1: pgrid.OffsetInt(off1),
				},
			}

			if wg.Contains(gridAxes.GetCenterPoint()) {
				ax0s := axis.Name[ax0%GridLinesTotal]
				ax1s := axis.Name[ax1%GridLinesTotal]
				point := fmt.Sprintf("%s%d%s%s%d", ax0s, off0, dir, ax1s, off1)
				if !yield(fmt.Sprintf(" -i %s", point)) {
					return
				}
				if debug != nil {
					debug.WriteString(" ")
					debug.WriteString(point)
				}

				count += 1
			}
		}
	}
}
