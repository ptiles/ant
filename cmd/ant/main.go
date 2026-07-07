package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/utils/palette"
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
)

type Flags struct {
	jsonStats    bool
	maxDimension int
	openResult   bool
	openResults  bool
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.BoolVar(&flags.jsonStats, "j", false, "Save stats in json file")
	flag.IntVar(&flags.maxDimension, "w", 16*1024, "Max image size")
	flag.BoolVar(&flags.openResult, "o", false, "Open resulting file")
	flag.BoolVar(&flags.openResults, "oo", false, "Open intermediate resulting files")

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

	field := pgrid.New(commonFlags.Pattern, commonFlags.AntRules, commonFlags.InitialPoint)

	if commonFlags.QuitOutside && field.InitialPointOutside(commonFlags.Rectangle) {
		fmt.Fprintln(os.Stderr, "Initial point is outside of rectangle.")
		os.Exit(0)
	}

	var pal palette.Palette
	if commonFlags.Monochrome {
		seedString := field.SeedString(8)
		pal = palette.GetPaletteMonochromatic(int(field.Limit), seedString)
	} else if commonFlags.Monochrome0 {
		seedString := field.SeedString(0)
		pal = palette.GetPaletteMonochromatic(int(field.Limit), seedString)
	} else {
		pal = palette.GetPaletteRainbow(int(field.Limit))
	}

	modifiedImagesCh := make(chan step.ModifiedImage, 64)

	go step.ModifiedPointsStepper(
		field,
		modifiedImagesCh, pal,
		commonFlags.Steps,
		commonFlags.MaxNoisyDots,
		commonFlags.MinTailSteps,
		commonFlags.MinTailSize,
	)

	fileNameFmt := fmt.Sprintf(
		"%s/%%s%s__%v__%s__%%s.%%s",
		commonFlags.Dir, commonFlags.AntName, commonFlags.Pattern, commonFlags.InitialPoint,
	)

	stepsTotal := saveImageFromModifiedImages(modifiedImagesCh, fileNameFmt, flags, commonFlags)

	if flags.openResult || flags.openResults {
		fileName := fmt.Sprintf(fileNameFmt, "", utils.WithSeparators(stepsTotal), "png")
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
