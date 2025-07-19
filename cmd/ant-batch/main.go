package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/utils"
	"image"
	"os"
	"slices"
	"strings"
)

type Flags struct {
	antNameRange  string
	patternsCount int

	initialAxis1        string
	initialDirection    string
	initialAxis2        string
	initialLines        string
	initialOffsets      string
	initialPointWythoff string

	initialPointCount int
	initialPointRange string
	initialPointNear  string
	initialPointPath  string
	kaleidoscope      bool

	Rectangle   image.Rectangle
	ScaleFactor int

	debug bool
}

func main() {
	flags := &Flags{}

	flag.StringVar(&flags.antNameRange, "nr", "", "Ant name range MIN-MAX")
	flag.IntVar(&flags.patternsCount, "pc", 0, "Patterns random count")

	flag.Func("ia", "Initial axes and direction (ex: A+C), use with -io or -iw", func(axisPairStr string) (err error) {
		flags.initialAxis1, flags.initialDirection, flags.initialAxis2 = utils.ParseInitialAxes(axisPairStr)
		return
	})
	flag.StringVar(&flags.initialLines, "il", "", "Initial point lines (comma separated)")
	flag.StringVar(&flags.initialOffsets, "io", "", "Initial point offsets (comma separated), use with -ia")
	flag.StringVar(&flags.initialPointWythoff, "iw", "", "Initial point offsets from wythoff array 'min-max%delta', use with -ia")

	flag.IntVar(&flags.initialPointCount, "ic", 0, "Initial point random count")
	flag.StringVar(&flags.initialPointRange, "ir", "0-8192", "Initial point offsets range")
	flag.StringVar(&flags.initialPointNear, "in", "", "Initial point near point")
	flag.StringVar(&flags.initialPointPath, "ip", "", "Initial points from ant path")
	flag.BoolVar(&flags.kaleidoscope, "ik", false, "Initial point kaleidoscope style")

	flag.Func("rs", "Output image rectangle size", func(rectangleStr string) (err error) {
		flags.Rectangle, flags.ScaleFactor, err = utils.ParseRectangleStr(rectangleStr)
		return
	})

	flag.BoolVar(&flags.debug, "d", false, "Print values")

	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	var debug strings.Builder

	antNames := slices.Collect(flags.AntNames(&debug))
	patterns := slices.Collect(flags.Patterns(&debug))

	var initialPoints []string
	initialPoints = slices.AppendSeq(initialPoints, flags.ListLines(&debug))
	initialPoints = slices.AppendSeq(initialPoints, flags.ListOffsets(&debug))
	initialPoints = slices.AppendSeq(initialPoints, flags.WythoffOffsets(&debug))
	initialPoints = slices.AppendSeq(initialPoints, flags.PathPoints(&debug))
	initialPoints = slices.AppendSeq(initialPoints, flags.InitialPoints(&debug))
	if len(initialPoints) == 0 {
		initialPoints = []string{""}
	}

	for _, antName := range antNames {
		for _, pattern := range patterns {
			for _, initialPoint := range initialPoints {
				printArgs(antName, pattern, initialPoint)
			}
		}
	}

	if flags.debug {
		fmt.Fprintln(os.Stderr, "Values used:", debug.String())
	}
}

func printArgs(all ...string) {
	var sb strings.Builder
	for _, one := range all {
		sb.WriteString(one)
	}
	fmt.Println(sb.String())
}
