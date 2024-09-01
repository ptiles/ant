package utils

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type CommonFlags struct {
	Cpuprofile   string
	Dir          string
	InitialPoint string
	AntName      string
	Radius       float64
	MaxSteps     int
	Verbose      bool
}

func CommonFlagsSetup(gridLinesTotal uint8) *CommonFlags {
	commonFlags := &CommonFlags{}

	flag.StringVar(&commonFlags.Cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&commonFlags.Dir, "d", fmt.Sprintf("results%d", gridLinesTotal), "Results directory")
	flag.StringVar(&commonFlags.InitialPoint, "i", "A0+B0", "Initial axes and direction")
	flag.StringVar(&commonFlags.AntName, "n", "", "Ant name")
	flag.IntVar(&commonFlags.MaxSteps, "s", 1000000, "Steps")
	flag.Float64Var(&commonFlags.Radius, "tr", 0.5, "Tiles config - radius")
	flag.BoolVar(&commonFlags.Verbose, "v", false, "Verbose output")

	return commonFlags
}

func ParseArgs(commonFlags *CommonFlags) {
	args := flag.Args()

	if len(args) < 1 {
		return
	}

	argSplit := strings.Split(args[0], ".")

	switch len(argSplit) {
	case 1:
		println(argSplit[0])
		commonFlags.AntName = argSplit[0]
	case 2:
		println(argSplit[0], argSplit[1])
		commonFlags.AntName = argSplit[0]
		maxStepsFromArg, err := strconv.Atoi(argSplit[1])
		if err == nil {
			commonFlags.MaxSteps = maxStepsFromArg
		}
	case 3:
		println(argSplit[0], argSplit[1], argSplit[2])
		commonFlags.AntName = argSplit[0]
		commonFlags.InitialPoint = argSplit[1]
		maxStepsFromArg, err := strconv.Atoi(argSplit[2])
		if err == nil {
			commonFlags.MaxSteps = maxStepsFromArg
		}
	}
}
