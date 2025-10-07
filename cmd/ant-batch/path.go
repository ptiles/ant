package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/seq"
	"github.com/ptiles/ant/utils"
	"iter"
	"slices"
	"strconv"
	"strings"
)

type path struct {
	flags    utils.CommonFlags
	gridSize int
	present  bool
}

func (p *path) pathParser() flagParser {
	return func(path string) error {
		if path == "" {
			return nil
		}

		p.flags.ParseShorthand(path)
		p.present = true

		return nil
	}
}

func (p *path) gridSizeParser() flagParser {
	return func(gridSizeStr string) error {
		if gridSizeStr == "" {
			return nil
		}

		gridSize, err := strconv.Atoi(gridSizeStr)
		if err != nil {
			return err
		}

		p.gridSize = gridSize

		return nil
	}
}

func (p *path) skip() bool {
	return !p.present
}

func (p *path) seq(debug *strings.Builder) iter.Seq[string] {
	return func(yield func(string) bool) {
		if p.skip() {
			return
		}

		if debug != nil {
			debug.WriteString("\nPathPoints: ")
		}

		pathFlags := p.flags
		field := pgrid.New(pathFlags.Pattern, pathFlags.AntRules, pathFlags.InitialPoint)

		if p.gridSize > 0 {
			gi := newGridInfo(p.gridSize, field.InitialOffsets())

			for gridAxes := range field.RunAxes(pathFlags.Steps.Max) {
				for _, turnString := range gi.turnsOnGrid(gridAxes) {
					if !yield(fmt.Sprintf(" -i %s", turnString)) {
						return
					}
					if debug != nil {
						debug.WriteString(" ")
						debug.WriteString(turnString)
					}
				}
			}

			return
		}

		steps := pathFlags.Steps
		for i, turn := range field.RunTurns(pathFlags.Steps.Max) {
			if i >= steps.Min && steps.Inc > 0 && i%steps.Inc == 0 {
				turnString := turn.String()
				if !yield(fmt.Sprintf(" -i %s", turnString)) {
					return
				}
				if debug != nil {
					debug.WriteString(" ")
					debug.WriteString(turnString)
				}
			}
		}
	}
}

type void struct{}

type position struct {
	index  int
	offset pgrid.OffsetInt
	min    pgrid.OffsetInt
	max    pgrid.OffsetInt
	prev   pgrid.OffsetInt
	next   pgrid.OffsetInt
	near   bool
}

type gridInfo struct {
	gridOffsets []pgrid.OffsetInt
	curr        [pgrid.GridLinesTotal]position
	visited     map[pgrid.GridAxes]void
	gap         pgrid.OffsetInt
}

func newGridInfo(gridSize int, offsets pgrid.GridOffsets) gridInfo {
	gi := gridInfo{
		gridOffsets: make([]pgrid.OffsetInt, 0, len(seq.WythoffReverse)),
		visited:     make(map[pgrid.GridAxes]void),
		//gap:         12,
		gap: 128,
	}

	for i, rowCol := range seq.WythoffReverse {
		if rowCol.Col >= gridSize {
			gi.gridOffsets = append(gi.gridOffsets, pgrid.OffsetInt(i))
		}
	}

	slices.Sort(gi.gridOffsets)

	for ax := range pgrid.GridLinesTotal {
		off := offsets[ax]

		index, found := slices.BinarySearch(gi.gridOffsets, off)

		if !found {
			left := off - gi.gridOffsets[index-1]
			right := gi.gridOffsets[index] - off

			if left < right {
				index -= 1
			}
		}

		gi.curr[ax] = gi.fromIndex(index, off)
		gi.curr[ax].near = false
	}

	return gi
}

func (gi *gridInfo) fromIndex(index int, off pgrid.OffsetInt) position {
	offset := gi.gridOffsets[index]

	return position{
		index:  index,
		offset: offset,
		min:    offset - gi.gap,
		max:    offset + gi.gap,
		prev:   gi.gridOffsets[index-1] + gi.gap,
		next:   gi.gridOffsets[index+1] - gi.gap,
		near:   (offset-gi.gap < off) && (off < offset+gi.gap),
	}
}

const triangular = pgrid.GridLinesTotal * (pgrid.GridLinesTotal - 1) / 2

func (gi *gridInfo) turnsOnGrid(gridAxes pgrid.GridAxes) []string {
	turns := make([]string, 0, triangular)

	ax0 := gridAxes.Axis0
	off0 := gridAxes.Coords.Offset0
	ax0NearNow := gi.curr[ax0].min < off0 && off0 < gi.curr[ax0].max

	ax1 := gridAxes.Axis1
	off1 := gridAxes.Coords.Offset1
	ax1NearNow := gi.curr[ax1].min < off1 && off1 < gi.curr[ax1].max

	//changed := (ax0NearNow != gi.curr[ax0].near) || (ax1NearNow != gi.curr[ax1].near)
	changed := (ax0NearNow && !gi.curr[ax0].near) || (ax1NearNow && !gi.curr[ax1].near)

	gi.curr[ax0].near = ax0NearNow
	gi.curr[ax1].near = ax1NearNow

	if changed {
		for axA, axB := range pgrid.AxesCanon() {
			if !gi.curr[axA].near || !gi.curr[axB].near {
				continue
			}

			ga := pgrid.GridAxes{
				Axis0: axA, Axis1: axB,
				Coords: pgrid.GridCoords{Offset0: gi.curr[axA].offset, Offset1: gi.curr[axB].offset},
			}
			if _, ok := gi.visited[ga]; !ok {
				gi.visited[ga] = void{}
				turns = append(turns, ga.TurnString("+"), ga.TurnString("-"))

				ga.Axis0, ga.Axis1 = ga.Axis1, ga.Axis0
				ga.Coords.Offset0, ga.Coords.Offset1 = ga.Coords.Offset1, ga.Coords.Offset0
				turns = append(turns, ga.TurnString("+"), ga.TurnString("-"))
			}
		}
	}

	if off0 < gi.curr[ax0].prev {
		gi.curr[ax0] = gi.fromIndex(gi.curr[ax0].index-1, off0)
	} else if off0 > gi.curr[ax0].next {
		gi.curr[ax0] = gi.fromIndex(gi.curr[ax0].index+1, off0)
	}
	if off0 < gi.gridOffsets[gi.curr[ax0].index-1] || gi.gridOffsets[gi.curr[ax0].index+1] < off0 {
		panic("jump by more than 1")
	}

	if off1 < gi.curr[ax1].prev {
		gi.curr[ax1] = gi.fromIndex(gi.curr[ax1].index-1, off1)
	} else if off1 > gi.curr[ax1].next {
		gi.curr[ax1] = gi.fromIndex(gi.curr[ax1].index+1, off1)
	}
	if off1 < gi.gridOffsets[gi.curr[ax1].index-1] || gi.gridOffsets[gi.curr[ax1].index+1] < off1 {
		panic("jump by more than 1")
	}

	return turns
}
