package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/canv"
	"github.com/ptiles/ant/pgrid"
	"os"
	"path/filepath"
	"strconv"
)

var (
	programName = filepath.Base(os.Args[0])
	usageText   = `Run Langton's ant on Penrose tiling (pentagrid)

Usage of %s:
	%s [flags] [name LLLRLRL...] [steps]

Name should consist of letters R, L, r, l.
Steps (default: 50000) should be a positive integer.

Flags:
`
	usageTextShort = "\nFor usage run: %s -h\n"
)

const (
	maxStepsDefault = 100000
)

func walk(currPoint pgrid.GridPoint, steps []bool, maxValue uint8) (bool, bool) {
	value := (currPoint.Get() + 1) % maxValue
	currPoint.Set(value)
	return steps[value], value == 0
}

func main() {
	var r int
	var dist int
	var phi0 int
	var width int
	var height int
	var antName string
	var maxSteps int

	flag.IntVar(&r, "r", 2, "Radius")
	flag.IntVar(&dist, "d", 8, "Distance")
	flag.IntVar(&phi0, "a", 0, "Angle")
	flag.IntVar(&width, "W", 1024, "Canvas width")
	flag.IntVar(&height, "H", 768, "Canvas height")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, programName, programName)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 0:
		fmt.Fprintln(os.Stderr, "Name is required")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	case 1:
		antName = args[0]
		maxSteps = maxStepsDefault
	case 2:
		antName = args[0]
		var err error
		maxSteps, err = strconv.Atoi(args[1])
		if err != nil {
			maxSteps = maxStepsDefault
		}
	default:
		antName = args[0]
		var err error
		maxSteps, err = strconv.Atoi(args[1])
		if err != nil {
			maxSteps = maxStepsDefault
		}
		fmt.Fprintln(os.Stderr, "Warning: Extra positional arguments ignored")
	}

	limit := uint8(len(antName))
	var nameInvalid = limit < 2
	rules := make([]bool, limit)
	for i, letter := range antName {
		if letter != 'R' && letter != 'r' && letter != 'L' && letter != 'l' {
			nameInvalid = true
		}
		rules[i] = letter == 'R' || letter == 'r'
	}
	if nameInvalid {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")

		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	fmt.Printf("Size: %dx%d; Name: %s; Steps: %d\n", width, height, antName, maxSteps)

	fileName := fmt.Sprintf("results/%s-%d.svg", antName, maxSteps)
	//fmt.Printf("Writing result to %s\n", fileName)
	canvas := canv.New(fileName, width/2, height/2)
	defer canvas.Close()

	//canvas.DrawOrigin()

	field := pgrid.New(float64(r), float64(dist), phi0, &canvas)

	// Draw grid
	//for ax := range field.Axes {
	//	for off := -15; off < 16; off++ {
	//		//color := ax + qi(off%5 == 0, 5, 10)
	//		color := ax + 10
	//		canvas.DrawLine(field.MakeGridLine(uint8(ax), int16(off)).Line, color)
	//	}
	//}

	initialLine := field.MakeGridLine(pgrid.E, 0)
	prevLine := field.MakeGridLine(pgrid.A, 0)
	currLine := field.MakeGridLine(pgrid.B, 0)

	prevPoint := field.MakeGridPoint(initialLine, prevLine, "")
	//field.DrawGridPoint(prevPoint, "")

	currPoint := field.MakeGridPoint(prevLine, currLine, "")
	//field.DrawGridPoint(currPoint, "")

	for st := 0; st < maxSteps; st++ {
		isRightTurn, _ := walk(currPoint, rules, limit)
		prevPoint, currPoint, prevLine, currLine = field.NextPoint(prevPoint, currPoint, prevLine, currLine, isRightTurn)
	}
	fmt.Printf("%s  %dx%d\n", fileName, pgrid.MaxOffset0-pgrid.MinOffset0, pgrid.MaxOffset1-pgrid.MinOffset1)

	for a0 := pgrid.MinOffset0; a0 <= pgrid.MaxOffset0; a0++ {
		for a1 := pgrid.MinOffset1; a1 <= pgrid.MaxOffset1; a1++ {
			points, colors := field.GetPointsByOffsets(a0, a1)
			for i := range points {
				canvas.DrawPoint(points[i], "", colors[i])
				//field.DrawGridPoint(point, "")
			}
		}
	}
}
