package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"image"
	"os"
	"path/filepath"
	"strconv"
)

var (
	programName = filepath.Base(os.Args[0])
	usageText   = `Run Langton's ant on Penrose tiling (pentagrid)

Usage of %s:
	%s [flags] [name RLRRLRR...] [steps]

Name should consist of letters R, L, r, l.
Steps (default: 1000000) should be a positive integer.

Flags:
`
	usageTextShort = "\nFor usage run: %s -h\n"
)

const (
	maxStepsDefault = 1000000
)

func main() {
	var r int
	var dist int
	var minWidth int
	var minHeight int
	var whiteBackground bool
	var openResults bool
	var openResult bool
	var antName string
	var maxSteps int
	var partialSteps int
	var startingPoint string

	var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
	flag.IntVar(&r, "r", 2, "Radius")
	flag.IntVar(&dist, "d", 8, "Distance")
	flag.StringVar(&startingPoint, "s", "A0+B0", "Starting axes and direction")
	flag.IntVar(&partialSteps, "p", 0, "Save partial result every N steps, 0 to disable")
	flag.BoolVar(&whiteBackground, "w", false, "White background")
	flag.BoolVar(&openResults, "oo", false, "Open partial resulting files")
	flag.BoolVar(&openResult, "o", false, "Open resulting file")
	flag.IntVar(&minWidth, "W", 128, "Canvas min width")
	flag.IntVar(&minHeight, "H", 128, "Canvas min height")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, programName, programName)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	utils.StartCPUProfile(cpuprofile)
	defer utils.StopCPUProfile()

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

	rules, err := utils.GetRules(antName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	field := pgrid.New(float64(r), float64(dist), rules, startingPoint)
	palette := utils.GetPalette(int(field.Limit), whiteBackground)

	modifiedImagesCh := make(chan *image.RGBA, 1024)

	go field.ModifiedImagesStepper(modifiedImagesCh, maxSteps, palette)

	fileName := fmt.Sprintf("results/%s-%s-%d.png", antName, startingPoint, maxSteps)

	saveImageFromModifiedImages(modifiedImagesCh, fileName, maxSteps)

	if openResult || openResults {
		utils.Open(fileName)
	}
}
