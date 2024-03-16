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
	var (
		cpuprofile    string
		dist          int
		startingPoint string
		antName       string
		//openResults   bool
		openResult bool
		//partialSteps  int
		radius   int
		maxSteps int
		verbose  bool
	)

	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.IntVar(&dist, "d", 8, "Distance")
	flag.StringVar(&startingPoint, "i", "A0+B0", "Initial axes and direction")
	flag.StringVar(&antName, "n", "", "Ant name")
	//flag.BoolVar(&openResults, "oo", false, "Open partial resulting files")
	flag.BoolVar(&openResult, "o", false, "Open resulting file")
	//flag.IntVar(&partialSteps, "p", 0, "Save partial result every N steps, 0 to disable")
	flag.IntVar(&radius, "r", 2, "Radius")
	flag.IntVar(&maxSteps, "s", maxStepsDefault, "Steps")
	flag.BoolVar(&verbose, "v", false, "Verbose output")

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
		if antName == "" {
			fmt.Fprintf(os.Stderr, "Name is required. Try to run: %s LLLRLRL", programName)
			fmt.Fprintf(os.Stderr, usageTextShort, programName)
			os.Exit(1)
		}
	case 1:
		antName = args[0]
	case 2:
		antName = args[0]
		maxStepsFromArg, err := strconv.Atoi(args[1])
		if err == nil {
			maxSteps = maxStepsFromArg
		}
	default:
		antName = args[0]
		maxStepsFromArg, err := strconv.Atoi(args[1])
		if err == nil {
			maxSteps = maxStepsFromArg
		}
		fmt.Fprintln(os.Stderr, "Warning: Extra positional arguments ignored")
	}

	rules, err := utils.GetRules(antName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	field := pgrid.New(float64(radius), float64(dist), rules, startingPoint, verbose)
	palette := utils.GetPalette(int(field.Limit))

	modifiedImagesCh := make(chan *image.RGBA, 1024)

	go field.ModifiedImagesStepper(modifiedImagesCh, maxSteps, palette)

	fileName := fmt.Sprintf("results/%s-%s-%d.png", antName, startingPoint, maxSteps)

	saveImageFromModifiedImages(modifiedImagesCh, fileName, maxSteps)

	//if openResult || openResults {
	if openResult {
		utils.Open(fileName)
	}
}
