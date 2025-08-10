package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"iter"
	"math/rand/v2"
	"strconv"
	"strings"
)

type near struct {
	count    int
	distance int

	axis1        int
	offset1      int
	axis2        int
	offset2      int
	pointPresent bool
}

func (n *near) pointParser() flagParser {
	return func(point string) error {
		if point == "" {
			return nil
		}

		ax1, off1, _, ax2, off2 := utils.ParseInitialPoint(point)

		n.axis1 = ax1
		n.offset1 = off1
		n.axis2 = ax2
		n.offset2 = off2
		n.pointPresent = true

		return nil
	}
}

func (n *near) countParser() flagParser {
	return func(countStr string) error {
		if countStr == "" {
			return nil
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return err
		}

		n.count = count

		return nil
	}
}

func (n *near) distanceParser() flagParser {
	return func(nearDistance string) error {
		if nearDistance == "" {
			return nil
		}

		distance, err := strconv.Atoi(nearDistance)
		if err != nil {
			return err
		}

		n.distance = distance

		return nil
	}
}

func (n *near) skip() bool {
	return n.count == 0 || n.distance == 0 || !n.pointPresent
}

func (n *near) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if n.skip() {
			return
		}

		min1 := n.offset1 - n.distance
		max1 := n.offset1 + n.distance
		min2 := n.offset2 - n.distance
		max2 := n.offset2 + n.distance

		if debug != nil {
			debug.WriteString("\nInitialPoints near:")
		}
		for range n.count {
			point := genRandomPointAround(n.axis1, min1, max1, n.axis2, min2, max2)
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

func genRandomPointAround(ax1, min1, max1, ax2, min2, max2 int) string {
	dir := [2]string{"-", "+"}[rand.IntN(2)]

	off1 := rand.IntN(max1+1-min1) + min1
	off2 := rand.IntN(max2+1-min2) + min2

	return fmt.Sprintf("%s%d%s%s%d", pgrid.AxisNames[ax1], off1, dir, pgrid.AxisNames[ax2], off2)
}
