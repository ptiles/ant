package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/canv"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/store"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
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

func walk(coords store.PackedCoordinates, steps []bool, maxValue uint8) bool {
	value := (store.Get(coords) + 1) % maxValue
	store.Set(coords, value)
	return steps[value]
}

func walk2(coords store.PackedCoordinates, steps []bool, maxValue uint8) bool {
	value := (store.Get2(coords) + 1) % maxValue
	store.Set2(coords, value)
	return steps[value]
}

func main() {
	var r int
	var dist int
	var phi0 int
	var minWidth int
	var minHeight int
	var antName string
	var maxSteps int

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.IntVar(&r, "r", 2, "Radius")
	flag.IntVar(&dist, "d", 8, "Distance")
	flag.IntVar(&phi0, "a", 0, "Angle")
	flag.IntVar(&minWidth, "W", 1024, "Canvas min width")
	flag.IntVar(&minHeight, "H", 768, "Canvas min height")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, programName, programName)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 0:
		fmt.Fprintf(os.Stderr, "Name is required. Try to run: %s LLLRLRL", programName)
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

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
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

	field := pgrid.New(float64(r), float64(dist), phi0)

	initialLine := pgrid.GridLine{Axis: pgrid.E, Offset: 0}
	prevLine := pgrid.GridLine{Axis: pgrid.A, Offset: 0}
	currLine := pgrid.GridLine{Axis: pgrid.B, Offset: 0}

	prevPoint := field.MakeGridPoint(initialLine, prevLine)
	currPoint := field.MakeGridPoint(prevLine, currLine)

	//store.Allocate(13)

	for st := 0; st < maxSteps; st++ {
		isRightTurn := walk(currPoint.PackedCoords, rules, limit)
		//isRightTurn := walk2(currPoint.PackedCoords, rules, limit)
		prevPoint, currPoint, prevLine, currLine = field.NextPoint(prevPoint, currPoint, prevLine, currLine, isRightTurn)
	}
	fileName := fmt.Sprintf("results/%s-%d.svg", antName, maxSteps)
	fmt.Printf("%s  %dx%d\n", fileName, store.MaxOffset0-store.MinOffset0, store.MaxOffset1-store.MinOffset1)

	maxX := minWidth/2 - 20
	maxY := minHeight/2 - 20
	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		//store.ForEach2(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
		point := field.MakeGridPoint(line0, line1).Point
		pointX := int(math.Abs(point[0]))
		if pointX > maxX {
			maxX = pointX
		}
		pointY := int(math.Abs(point[1]))
		if pointY > maxY {
			maxY = pointY
		}
	})
	maxX += 20
	maxY += 20

	fmt.Printf("Size: %dx%d; Name: %s; Steps: %d\n", maxX*2, maxY*2, antName, maxSteps)
	canvas := canv.New(fileName, maxX, maxY, int(limit))
	defer canvas.Close()

	//fmt.Printf("Writing result to %s\n", fileName)
	//canvas.DrawOrigin()
	// Draw grid
	//for ax := 0; ax < 5; ax++ {
	//	for off := -15; off < 16; off++ {
	//		//color := ax + qi(off%5 == 0, 5, 10)
	//		color := ax + 10
	//		canvas.DrawLine(field.MakeGridLine(uint8(ax), int16(off)).Line, color)
	//	}
	//}
	//
	//canvas.DrawPoint(prevPoint.Point, 0, "E0 A0")
	//canvas.DrawPoint(currPoint.Point, 0, "A0 B0")

	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		//store.ForEach2(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
		point := field.MakeGridPoint(line0, line1).Point
		//if canvas.IsOutside(point) {
		//	return
		//}
		canvas.DrawPoint(point, color)
	})
}
