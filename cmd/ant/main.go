package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"image"
	"os"
	"path"
	"path/filepath"
)

var (
	programName = filepath.Base(os.Args[0])
	usageText   = `Run Langton's ant on Penrose tiling (pentagrid)

Usage of %s:
	%s [flags] [name RLRRLRR...].[initial point A0+B0].[steps number]

Name should consist of letters R, L, r, l.
Steps (default: 1000000) should be a positive integer.

Flags:
`
	usageTextShort = "\nFor usage run: %s -h\n"
)

type Flags struct {
	jsonStats     bool
	maxDimension  int
	openResults   bool
	openResult    bool
	partialImages int
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.BoolVar(&flags.jsonStats, "j", false, "Save stats in json file")
	flag.IntVar(&flags.maxDimension, "w", 16*1024, "Max image size")
	//flag.BoolVar(&flags.openResults, "oo", false, "Open partial resulting files")
	flag.BoolVar(&flags.openResult, "o", false, "Open resulting file")
	flag.IntVar(&flags.partialImages, "p", 0, "Save partial result every N intermediate images")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, programName, programName)
		flag.PrintDefaults()
	}

	return flags
}

func main() {
	commonFlags := utils.CommonFlagsSetup()
	flags := flagsSetup()
	flag.Parse()
	utils.ParseArgs(commonFlags)

	utils.StartCPUProfile(commonFlags.Cpuprofile)
	defer utils.StopCPUProfile()

	rules, err := utils.GetRules(commonFlags.AntName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	field := pgrid.New(commonFlags.Radius, rules, commonFlags.InitialPoint, commonFlags.Verbose)
	palette := utils.GetPalette(int(field.Limit))

	modifiedImagesCh := make(chan *image.RGBA, 1024)

	go field.ModifiedPointsStepper(modifiedImagesCh, commonFlags.MaxSteps, palette)

	fileNameFmt := fmt.Sprintf("%s/%s.%s.%%d.%%s", path.Clean(commonFlags.Dir), commonFlags.AntName, commonFlags.InitialPoint)

	saveImageFromModifiedImages(modifiedImagesCh, fileNameFmt, flags, commonFlags)

	//if flags.openResult || flags.openResults {
	if flags.openResult {
		fileName := fmt.Sprintf(fileNameFmt, commonFlags.MaxSteps)
		utils.Open(fileName)
	}
}
