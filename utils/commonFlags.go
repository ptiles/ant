package utils

import (
	"flag"
	"fmt"
	"image"
	"math"
	"path"
	"strconv"
	"strings"
)

type CommonFlags struct {
	Cpuprofile string
	Memprofile string

	Dir          string
	InitialPoint string
	AntName      string
	Radius       float64

	MaxSteps     uint64
	MaxNoisyDots uint64
	MinStepsPct  uint64
	MinUniqPct   uint64

	Rectangle   image.Rectangle
	ScaleFactor int
	Monochrome  bool
	Alpha       bool
}

func (cf *CommonFlags) String() string {
	return fmt.Sprintf(
		"%s__%f__%s__%s\n",
		cf.AntName, cf.Radius, cf.InitialPoint, WithUnderscores(cf.MaxSteps),
	)
}

func (cf *CommonFlags) CommonFlagsSetup(gridLinesTotal uint8) {
	flag.StringVar(&cf.Cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&cf.Memprofile, "memprofile", "", "Write mem profile to file")
	flag.StringVar(&cf.Dir, "d", fmt.Sprintf("results%d", gridLinesTotal), "Results directory")
	flag.StringVar(&cf.InitialPoint, "i", "A0+B0", "Initial axes and direction")
	flag.BoolVar(&cf.Monochrome, "m", false, "Monochromatic palette")
	flag.BoolVar(&cf.Alpha, "alpha", false, "Save transparent image with alpha channel")
	flag.StringVar(&cf.AntName, "n", "RLL", "Ant name")
	flag.Func("r", "Output image rectangle", func(rectangleStr string) (err error) {
		cf.Rectangle, cf.ScaleFactor, err = ParseRectangleStr(rectangleStr)
		return
	})
	flag.Uint64Var(&cf.MaxSteps, "s", 1_000_000, "Steps")
	flag.Uint64Var(&cf.MaxNoisyDots, "sn", 0, "Max noisy dots")
	flag.Uint64Var(&cf.MinStepsPct, "sm", 0, "Min steps percent")
	flag.Uint64Var(&cf.MinUniqPct, "su", 0, "Min uniq points percent")
	flag.Float64Var(&cf.Radius, "tr", 0.000007, "Tiles config - radius")
}

func (cf *CommonFlags) ParseArgs() {
	shorthand := flag.Arg(0)
	if shorthand != "" {
		cf.ParseShorthand(shorthand)
	}

	if cf.MaxNoisyDots == 0 {
		cf.MaxNoisyDots = math.MaxUint64
	}

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
