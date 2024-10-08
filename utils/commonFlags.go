package utils

import (
	"flag"
	"fmt"
	"image"
	"path"
	"strconv"
	"strings"
)

type CommonFlags struct {
	Cpuprofile   string
	Dir          string
	InitialPoint string
	AntName      string
	Radius       float64
	MaxSteps     uint64
	Rectangle    image.Rectangle
	ScaleFactor  int
	Monochrome   bool
}

func (cf *CommonFlags) String() string {
	return fmt.Sprintf(
		"%s__%f__%s__%d\n",
		cf.AntName, cf.Radius, cf.InitialPoint, cf.MaxSteps,
	)
}

func (cf *CommonFlags) CommonFlagsSetup(gridLinesTotal uint8) {
	flag.StringVar(&cf.Cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&cf.Dir, "d", fmt.Sprintf("results%d", gridLinesTotal), "Results directory")
	flag.StringVar(&cf.InitialPoint, "i", "A0+B0", "Initial axes and direction")
	flag.BoolVar(&cf.Monochrome, "m", false, "Monochromatic palette")
	flag.StringVar(&cf.AntName, "n", "RLLLL", "Ant name")
	flag.Func("r", "Output image rectangle", func(rectangleStr string) (err error) {
		cf.Rectangle, cf.ScaleFactor, err = ParseRectangleStr(rectangleStr)
		return
	})
	flag.Uint64Var(&cf.MaxSteps, "s", 1000000, "Steps")
	flag.Float64Var(&cf.Radius, "tr", 0.5, "Tiles config - radius")
}

func (cf *CommonFlags) ParseArgs() {
	shorthand := flag.Arg(0)
	if shorthand == "" {
		return
	}
	cf.ParseShorthand(shorthand)
	cf.Dir = path.Clean(cf.Dir)
}

func (cf *CommonFlags) ParseShorthand(shorthand string) {
	antNameR := `[RL]+`
	radiusR := `[01]\.\d+`
	initialPointR := `[A-X]-?\d+[+-]?[A-X]-?\d+`
	maxStepsR := `[0-9_]+`

	expr := fmt.Sprintf(
		"(?P<antName>%s)__(?P<radius>%s)__(?P<initialPoint>%s)__(?P<maxSteps>%s)",
		antNameR, radiusR, initialPointR, maxStepsR,
	)

	matches := NamedMatches(expr, shorthand)

	cf.AntName = matches["antName"]

	radius, radiusErr := strconv.ParseFloat(matches["radius"], 64)
	if radiusErr == nil {
		cf.Radius = radius
	}

	cf.InitialPoint = matches["initialPoint"]

	maxStepsS := strings.Replace(matches["maxSteps"], "_", "", -1)
	maxSteps, maxStepsErr := strconv.ParseUint(maxStepsS, 10, 0)
	if maxStepsErr == nil {
		cf.MaxSteps = maxSteps
	}
}

func (cf *CommonFlags) parseShorthandOld(shorthand string) {
	shortSplit := strings.Split(shorthand, "__")

	switch len(shortSplit) {
	case 1:
		println(shortSplit[0])
		cf.AntName = shortSplit[0]
	case 2:
		println(shortSplit[0], shortSplit[1])
		cf.AntName = shortSplit[0]
		maxStepsFromShort, err := strconv.ParseUint(shortSplit[1], 10, 64)
		if err == nil {
			cf.MaxSteps = maxStepsFromShort
		}
	case 3:
		println(shortSplit[0], shortSplit[1], shortSplit[2])
		cf.AntName = shortSplit[0]
		cf.InitialPoint = shortSplit[1]
		maxStepsFromShort, err := strconv.ParseUint(shortSplit[2], 10, 64)
		if err == nil {
			cf.MaxSteps = maxStepsFromShort
		}
	}
}
