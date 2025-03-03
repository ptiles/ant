package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
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
	jsonStats    bool
	maxDimension int
	openResults  bool
	openResult   bool
	partialSteps uint64
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.BoolVar(&flags.jsonStats, "j", false, "Save stats in json file")
	flag.BoolVar(&flags.openResult, "o", false, "Open resulting file")
	flag.Uint64Var(&flags.partialSteps, "p", 0, "Save partial result every N steps")
	flag.IntVar(&flags.maxDimension, "w", 16*1024, "Max image size")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, programName, programName)
		flag.PrintDefaults()
	}
	flag.Parse()

	return flags
}

func main() {
	commonFlags := &utils.CommonFlags{}
	commonFlags.CommonFlagsSetup(pgrid.GridLinesTotal)
	flags := flagsSetup()
	commonFlags.ParseArgs()

	utils.StartCPUProfile(commonFlags.Cpuprofile)
	defer utils.StopCPUProfile()

	rules, err := utils.GetRules(commonFlags.AntName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	field := pgrid.New(commonFlags.Radius, rules, commonFlags.InitialPoint)
	var palette []color.RGBA
	if commonFlags.Monochrome {
		rng := pgrid.InitialPointSeed(field.InitialPoint)
		palette = utils.GetPaletteMonochromatic(int(field.Limit), rng)
	} else {
		palette = utils.GetPaletteRainbow(int(field.Limit))
	}

	modifiedImagesCh := make(chan pgrid.ModifiedImage, 64)

	go field.ModifiedPointsStepper(
		modifiedImagesCh, palette,
		commonFlags.MaxSteps, flags.partialSteps,
		commonFlags.MinCleanStreak, commonFlags.MaxNoisyDots,
	)

	fileNameFmt := fmt.Sprintf(
		"%s/%s__%f__%s__%%s.%%s",
		commonFlags.Dir, commonFlags.AntName, commonFlags.Radius, commonFlags.InitialPoint,
	)

	stepsTotal := saveImageFromModifiedImages(modifiedImagesCh, fileNameFmt, flags, commonFlags)

	//if flags.openResult || flags.openResults {
	if flags.openResult {
		fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(stepsTotal), "png")
		utils.Open(fileName)
	}

	if commonFlags.Memprofile != "" {
		f, mpErr := os.Create(commonFlags.Memprofile)
		if mpErr != nil {
			log.Fatal("could not create memory profile: ", mpErr)
		}
		defer f.Close()
		runtime.GC()
		if wrErr := pprof.WriteHeapProfile(f); wrErr != nil {
			log.Fatal("could not write memory profile: ", wrErr)
		}
	}
}
