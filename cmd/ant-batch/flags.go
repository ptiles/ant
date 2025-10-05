package main

import (
	"flag"
)

type flagParser func(string) error

type Flags struct {
	probability int
	debug       bool

	grid     grid
	interval interval
	lines    lines
	list     list
	names    names
	near     near
	path     path
	patterns patterns
	rect     rect
	wythoff  wythoff
}

func parseFlags() *Flags {
	fl := Flags{}
	flag.IntVar(&fl.probability, "p", 0, "Probability 1/p\n")
	flag.BoolVar(&fl.debug, "d", false, "Print values\n")

	flag.Func("grid", "Initial points inside rectangle, snapped to grid", fl.grid.rectParser())
	flag.Func("grid-min-axes", "Minimal number of axes intersecting", fl.grid.minAxesParser())
	flag.Func("grid-size", "Initial points snapped to grid size min[-max]\n", fl.grid.gridSizeParser())

	flag.Func("interval", "Initial point offsets interval", fl.interval.intervalParser())
	flag.Func("interval-count", "Initial point offsets interval count", fl.interval.countParser())
	flag.BoolFunc("interval-kaleidoscope", "Initial point kaleidoscope style\n", fl.interval.kaleidoscopeParser())

	flag.Func("lines", "Initial point from lines (comma separated)\n", fl.lines.parser())

	flag.Func("offsets", "Initial point offsets (comma separated)", fl.list.offsetsParser())
	flag.Func("offsets-axes", "Initial axes and direction (ex: A+C)\n", fl.list.axesParser())

	flag.Func("names", "Ant name range MIN-MAX\n", fl.names.parser())

	flag.Func("near", "Initial point near point", fl.near.pointParser())
	flag.Func("near-count", "Initial point near count", fl.near.countParser())
	flag.Func("near-distance", "Initial point near distance\n", fl.near.distanceParser())

	flag.Func("path", "Initial points from ant path", fl.path.pathParser())
	flag.Func("path-grid-size", "Initial points from path snapped to grid size\n", fl.path.gridSizeParser())

	flag.Func("patterns", "Patterns random count\n", fl.patterns.parser())

	flag.Func("rect", "Initial points inside rectangle", fl.rect.rectParser())
	flag.Func("rect-count", "Initial points inside rectangle count\n", fl.rect.countParser())

	flag.Func("wythoff", "Initial point offsets from wythoff array 'min-max%delta'", fl.wythoff.intervalParser())
	flag.Func("wythoff-axes", "Axes and direction for Wythoff offsets (ex: A+C)", fl.wythoff.axesParser())

	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	return &fl
}
