package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid/axis"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/wgrid"
	"image"
	"iter"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type grid struct {
	rect        image.Rectangle
	scaleFactor int

	gridSize    int
	gridSizeMax int
	minAxes     int
}

func (g *grid) rectParser() flagParser {
	return func(rectStr string) error {
		if rectStr == "" {
			return nil
		}

		rectangle, scaleFactor, err := utils.ParseRectangleStr(rectStr)
		if err != nil {
			return err
		}

		g.rect = rectangle
		g.scaleFactor = scaleFactor

		return nil
	}
}

func (g *grid) gridSizeParser() flagParser {
	return func(gridSizeStr string) error {
		if gridSizeStr == "" {
			return nil
		}

		expr := regexp.MustCompile(`(?P<min>\d+)(-(?P<max>\d+))?`)
		result := utils.NamedIntMatches(expr, gridSizeStr)

		g.gridSize = result["min"]
		g.gridSizeMax = result["max"]

		if g.gridSizeMax == 0 {
			g.gridSizeMax = math.MaxInt
		}

		return nil
	}
}

func (g *grid) minAxesParser() flagParser {
	return func(minAxesStr string) error {
		if minAxesStr == "" {
			return nil
		}

		minAxes, err := strconv.Atoi(minAxesStr)
		if err != nil {
			return err
		}

		g.minAxes = minAxes

		return nil
	}
}

func (g *grid) skip() bool {
	return g.rect.Empty() || g.gridSize == 0
}

func (g *grid) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if g.skip() {
			return
		}

		if debug != nil {
			debug.WriteString("\nGrid:")
		}

		if g.gridSizeMax == 0 {
			g.gridSizeMax = math.MaxInt
		}

		wg := wgrid.New(g.rect, g.scaleFactor)
		intersections := wg.IntersectionsMap(g.gridSize, g.gridSizeMax)
		startingPoints := 0

		for _, axes := range intersections {
			axesCount := len(axes)
			if axesCount >= g.minAxes {
				for ax0, off0 := range axes {
					for ax1, off1 := range axes {
						if ax0 == ax1 {
							continue
						}

						p1 := fmt.Sprintf(" -i %s%d%s%s%d", axis.Name[ax0], off0, "+", axis.Name[ax1], off1)
						if !yield(p1) {
							return
						}

						p2 := fmt.Sprintf(" -i %s%d%s%s%d", axis.Name[ax0], off0, "-", axis.Name[ax1], off1)
						if !yield(p2) {
							return
						}

						startingPoints += 2
					}
				}

				if debug != nil {
					debug.WriteString(fmt.Sprintf("\nx%d  |%s", axesCount, axes))
				}
			}
		}

		if debug != nil {
			debug.WriteString(fmt.Sprintf("\n%d intersections; %d starting points", len(intersections), startingPoints))
		}
	}
}
