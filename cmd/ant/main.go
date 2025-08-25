package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/output"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/utils/palette"
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
)

type Flags struct {
	jsonStats    bool
	maxDimension int
	gridSize     int
	gridBoth     bool
	gridEmpty    bool
	gridOnly     int
	openResults  bool
	openResult   bool
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.BoolVar(&flags.jsonStats, "j", false, "Save stats in json file")
	flag.IntVar(&flags.gridSize, "g", 0, "Save files with grid")
	flag.BoolVar(&flags.gridBoth, "gb", false, "Save both files with and without grid")
	flag.BoolVar(&flags.gridEmpty, "ge", false, "Save empty grid file")
	flag.IntVar(&flags.gridOnly, "go", 0, "Save only empty grid file")
	flag.BoolVar(&flags.openResult, "o", false, "Open resulting file")
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

	if flags.gridOnly > 0 && !commonFlags.Rectangle.Empty() {
		out := output.NewImage(commonFlags.Rectangle, commonFlags.ScaleFactor, flags.maxDimension)
		fileName := fmt.Sprintf("%s/grid_%d_%s.png", commonFlags.Dir, flags.gridOnly, out.RectCenteredString())
		out.SaveGridOnly(fileName, flags.gridOnly, commonFlags.Alpha)
		fmt.Println(fileName)
		os.Exit(0)
	}

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
	)

	fileNameFmt := fmt.Sprintf(
		"%s/%%s%s__%v__%s__%%s.%%s",
		commonFlags.Dir, commonFlags.AntName, commonFlags.Pattern, commonFlags.InitialPoint,
	)

	stepsTotal := saveImageFromModifiedImages(modifiedImagesCh, fileNameFmt, flags, commonFlags)

	if flags.openResult && flags.gridSize > 0 {
		gridPrefix := fmt.Sprintf("grid_%d_", flags.gridSize)
		fileName := fmt.Sprintf(fileNameFmt, gridPrefix, utils.WithSeparators(stepsTotal), "png")
		utils.Open(fileName)
	}
	//if flags.openResult || flags.openResults {
	if flags.openResult && flags.gridSize == 0 {
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
