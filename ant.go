package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/store"
	"github.com/ptiles/ant/utils"
	"image"
	"image/color"
	"image/png"
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
	var minWidth int
	var minHeight int
	var antName string
	var maxSteps int

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.IntVar(&r, "r", 2, "Radius")
	flag.IntVar(&dist, "d", 8, "Distance")
	flag.IntVar(&minWidth, "W", 128, "Canvas min width")
	flag.IntVar(&minHeight, "H", 128, "Canvas min height")

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

	field := pgrid.New(float64(r), float64(dist))

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
	fileName := fmt.Sprintf("results/%s-%d.png", antName, maxSteps)

	maxX := minWidth / 2
	maxY := minHeight / 2

	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		//store.ForEach2(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
		gp := field.MakeGridPoint(line0, line1)
		point := field.GetCenterPoint(&gp)
		pointX := int(math.Abs(point[0]))
		if pointX > maxX {
			maxX = pointX
		}
		pointY := int(math.Abs(point[1]))
		if pointY > maxY {
			maxY = pointY
		}
	})

	maxX = (maxX/128 + 1) * 128
	maxY = (maxY/128 + 1) * 128

	fmt.Printf("%s Name: %s; Steps: %d; Size: %dx%d\n", fileName, antName, maxSteps, maxX*2, maxY*2)
	img := image.NewPaletted(image.Rect(0, 0, maxX*2, maxY*2), getPalette(int(limit)))

	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		//store.ForEach2(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
		gp := field.MakeGridPoint(line0, line1)
		point := field.GetCenterPoint(&gp)
		img.SetColorIndex(int(point[0])+maxX, int(point[1])+maxY, color+1)
	})

	// Create a new file to save the PNG image
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image as a PNG and save it to the file
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func getPalette(steps int) color.Palette {
	var palette = make(color.Palette, steps+1)
	palette[0] = color.RGBA{R: 0, G: 0, B: 0, A: 0xff}

	for c := 0; c < steps; c++ {
		step := c * 360 / steps

		ra := step + 0*120 + 90
		ga := step + 1*120 + 90
		ba := step + 2*120 + 90

		rs := (1 + math.Sin(utils.FromDegrees(ra))) / 2
		gs := (1 + math.Sin(utils.FromDegrees(ga))) / 2
		bs := (1 + math.Sin(utils.FromDegrees(ba))) / 2

		r := uint8(math.Round(rs*0xd0 + 0x0f))
		g := uint8(math.Round(gs * 0xb0))
		b := uint8(math.Round(bs*0xb0 + 0x4f))

		palette[c+1] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return palette
}
