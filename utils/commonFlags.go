package utils

import (
	"flag"
	"fmt"
	"image"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type StepCounts struct {
	Min uint64
	Max uint64
	Inc uint64
}

func (sc StepCounts) String() string {
	return fmt.Sprintf("%d-%d%%%d", sc.Min, sc.Max, sc.Inc)
}

type CommonFlags struct {
	Cpuprofile string
	Memprofile string

	Dir          string
	InitialPoint string
	AntName      string
	AntRules     []bool
	Pattern      float64

	Steps StepCounts

	MaxNoisyDots uint64
	MinStepsPct  uint64
	MinUniqPct   uint64
	MinTailSteps uint64
	MinTailSize  uint64

	Rectangle   image.Rectangle
	ScaleFactor int
	QuitOutside bool
	Monochrome  bool
	Monochrome0 bool
	Alpha       bool
}

func (cf *CommonFlags) String() string {
	return fmt.Sprintf(
		"%s__%v__%s__%s\n",
		cf.AntName, cf.Pattern, cf.InitialPoint, WithSeparators(cf.Steps.Max),
	)
}

func parseStepsStr(stepsStr string) (StepCounts, error) {
	re := regexp.MustCompile(`((\d+)-)?(\d+)(%(\d+))?`)
	result := re.FindStringSubmatch(strings.Replace(stepsStr, "_", "", -1))

	minSt, _ := strconv.ParseUint(result[2], 0, 64)
	maxSt, _ := strconv.ParseUint(result[3], 0, 64)
	incSt, _ := strconv.ParseUint(result[5], 0, 64)

	return StepCounts{minSt, maxSt, incSt}, nil
}

func (cf *CommonFlags) CommonFlagsSetup(gridLinesTotal uint8) {
	flag.StringVar(&cf.Cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&cf.Memprofile, "memprofile", "", "Write mem profile to file")
	flag.StringVar(&cf.Dir, "d", fmt.Sprintf("results%d", gridLinesTotal), "Results directory")
	flag.StringVar(&cf.InitialPoint, "i", "A0+B0", "Initial axes and direction")
	flag.BoolVar(&cf.Monochrome, "m", false, "Monochromatic palette")
	flag.BoolVar(&cf.Monochrome0, "m0", false, "Monochromatic palette exact point")
	flag.BoolVar(&cf.Alpha, "alpha", false, "Save transparent image with alpha channel")
	flag.Func("n", "Ant name", func(antName string) (err error) {
		cf.AntName = antName
		cf.AntRules, err = GetRules(antName)
		return
	})
	flag.Func("r", "Output image rectangle", func(rectangleStr string) (err error) {
		cf.Rectangle, cf.ScaleFactor, err = ParseRectangleStr(rectangleStr)
		return
	})
	flag.BoolVar(&cf.QuitOutside, "rq", false, "Quit if initial point is outside of rectangle")
	flag.Func("s", "Steps", func(stepsStr string) (err error) {
		cf.Steps, err = parseStepsStr(stepsStr)
		return
	})
	flag.Uint64Var(&cf.MaxNoisyDots, "sn", 0, "Max noisy dots")
	flag.Uint64Var(&cf.MinStepsPct, "sm", 0, "Min steps percent")
	flag.Uint64Var(&cf.MinUniqPct, "su", 0, "Min uniq points percent")
	flag.Uint64Var(&cf.MinTailSteps, "st", 0, "Min steps in the loop to stop when the tail is reached")
	flag.Uint64Var(&cf.MinTailSize, "sts", 32768, "Tail size")
	flag.Float64Var(&cf.Pattern, "p", 7e-06, "Pattern radius")
}

func (cf *CommonFlags) ParseArgs() {
	shorthand := flag.Arg(0)
	if shorthand != "" {
		cf.ParseShorthand(shorthand)
	}

	if cf.AntName == "" {
		cf.AntName = "RLL"
		cf.AntRules = []bool{true, false, false}
	}

	if cf.Steps.Max == 0 {
		cf.Steps = StepCounts{Min: 0, Max: 5_000_000, Inc: 0}
	}

	if cf.MaxNoisyDots == 0 {
		cf.MaxNoisyDots = math.MaxUint64
	}

	cf.Dir = path.Clean(cf.Dir)
}

func (cf *CommonFlags) ParseShorthand(shorthand string) {
	antNameR := `[RL]+`
	radiusR := `([01]\.\d+)|(\d+e-\d+)`
	initialPointR := `[A-X]-?\d+[+-]?[A-X]-?\d+`
	stepsR := `(([0-9_]+)-)?([0-9_]+)(%([0-9_]+))?`

	expr := regexp.MustCompile(fmt.Sprintf(
		"(?P<antName>%s)__(?P<radius>%s)__(?P<initialPoint>%s)__(?P<steps>%s)",
		antNameR, radiusR, initialPointR, stepsR,
	))

	matches := NamedStringMatches(expr, shorthand)

	cf.AntName = matches["antName"]
	cf.AntRules, _ = GetRules(cf.AntName)

	radius, radiusErr := strconv.ParseFloat(matches["radius"], 64)
	if radiusErr == nil {
		cf.Pattern = radius
	}
	if radius == 0 {
		fmt.Println("Pattern cannot be zero")
		os.Exit(1)
	}
	if radius < 1e-10 {
		fmt.Println("Pattern too small, can be inaccurate")
	}

	cf.InitialPoint = matches["initialPoint"]

	cf.Steps, _ = parseStepsStr(matches["steps"])
}
