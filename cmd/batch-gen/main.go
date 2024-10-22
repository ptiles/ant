package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const GridLinesTotal = uint(pgrid.GridLinesTotal)

func genRandomPoint(min, max int) (uint, int, string, uint, int) {
	dirNames := [2]string{"-", "+"}

	ax := rand.Perm(int(GridLinesTotal))
	ax1, ax2 := uint(ax[0]), uint(ax[1])

	dir := dirNames[rand.IntN(2)]

	off1 := rand.IntN(max+1-min) + min
	off2 := rand.IntN(max+1-min) + min

	if rand.IntN(2) == 0 {
		off1 = -off1
	}
	if rand.IntN(2) == 0 {
		off2 = -off2
	}

	return ax1, off1, dir, ax2, off2
}

func genRandomPointString(min, max int) string {
	ax1, off1, dir, ax2, off2 := genRandomPoint(min, max)

	ax1s := pgrid.AxisNames[ax1]
	ax2s := pgrid.AxisNames[ax2]

	return fmt.Sprintf("%s%d%s%s%d", ax1s, off1, dir, ax2s, off2)
}

func genRandomPointKaleidoscope(min, max int) [GridLinesTotal]string {
	ax1, off1, dir, ax2, off2 := genRandomPoint(min, max)

	var result [GridLinesTotal]string
	for i := range GridLinesTotal {
		ax1s := pgrid.AxisNames[(ax1+i)%GridLinesTotal]
		ax2s := pgrid.AxisNames[(ax2+i)%GridLinesTotal]
		point := fmt.Sprintf("%s%d%s%s%d", ax1s, off1, dir, ax2s, off2)
		result[i] = point
	}
	return result
}

func genRandomPointAround(ax1, min1, max1, ax2, min2, max2 int) string {
	dirNames := [2]string{"-", "+"}

	dir := dirNames[rand.IntN(2)]

	off1 := rand.IntN(max1+1-min1) + min1
	off2 := rand.IntN(max2+1-min2) + min2

	return fmt.Sprintf("%s%d%s%s%d", pgrid.AxisNames[ax1], off1, dir, pgrid.AxisNames[ax2], off2)
}

func numToName(num uint64, bitWidth int) string {
	format := fmt.Sprintf("%%0%ds", bitWidth)
	binary := fmt.Sprintf(format, strconv.FormatUint(num, 2))
	return strings.Replace(strings.Replace(binary, "0", "L", -1), "1", "R", -1)
}

type Flags struct {
	antNameRange string

	initialPointCount int
	initialPointRange string
	initialPointNear  bool
	kaleidoscope      bool

	radiusCount int
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.antNameRange, "nr", "", "Ant name range MIN-MAX")

	flag.IntVar(&flags.initialPointCount, "ic", 0, "Initial point count")
	flag.BoolVar(&flags.kaleidoscope, "ik", false, "Initial point kaleidoscope style")
	flag.StringVar(&flags.initialPointRange, "ir", "0-8192", "Initial point offset range")
	flag.BoolVar(&flags.initialPointNear, "in", false, "Initial point near -i value")

	flag.IntVar(&flags.radiusCount, "rc", 0, "Random radius count")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	return flags
}

func parseNameRange(antNameRange string) (int, int) {
	result := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(antNameRange)
	minBitWidth, _ := strconv.Atoi(result[1])
	maxBitWidth, _ := strconv.Atoi(result[2])

	return minBitWidth, maxBitWidth
}

func getAntNames(flags *Flags, commonFlags *utils.CommonFlags) []string {
	if flags.antNameRange != "" {
		minBitWidth, maxBitWidth := parseNameRange(flags.antNameRange)

		// TODO: calculate antNames length without loop
		antNamesLength := 0
		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			antNamesLength += 1<<bitWidth - 2
		}

		antNames := make([]string, antNamesLength)
		i := 0
		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			maxNum := uint64(1<<bitWidth) - 1
			for num := uint64(1); num < maxNum; num++ {
				antNames[i] = numToName(num, bitWidth)
				i += 1
			}
		}
		return antNames
	} else if commonFlags.AntName != "" {
		return strings.Split(commonFlags.AntName, ",")
	}

	fmt.Fprintln(os.Stderr, "Ant name or range required")
	return []string{"LLR"}
}

func getInitialPoints(flags *Flags, commonFlags *utils.CommonFlags) []string {
	initialPointRangeMin, initialPointRangeMax, _ := utils.ParseRangeStr(flags.initialPointRange)
	if flags.initialPointNear && flags.initialPointCount > 0 {
		ax1, off1, _, ax2, off2 := utils.ParseInitialPoint(commonFlags.InitialPoint)

		min1 := off1 - initialPointRangeMax
		max1 := off1 + initialPointRangeMax
		min2 := off2 - initialPointRangeMax
		max2 := off2 + initialPointRangeMax

		initialPoints := make([]string, flags.initialPointCount)
		for i := range flags.initialPointCount {
			initialPoints[i] = genRandomPointAround(ax1, min1, max1, ax2, min2, max2)
		}
		return initialPoints
	} else if flags.kaleidoscope && flags.initialPointCount > 0 {
		initialPoints := make([]string, flags.initialPointCount*int(GridLinesTotal))
		for i := range flags.initialPointCount {
			points := genRandomPointKaleidoscope(initialPointRangeMin, initialPointRangeMax)
			for j, point := range points {
				initialPoints[i*int(GridLinesTotal)+j] = point
			}
		}
		return initialPoints
	} else if flags.initialPointCount > 0 {
		initialPoints := make([]string, flags.initialPointCount)
		for i := range flags.initialPointCount {
			initialPoints[i] = genRandomPointString(initialPointRangeMin, initialPointRangeMax)
		}
		return initialPoints
	}
	return strings.Split(commonFlags.InitialPoint, ",")
}

func getRadii(flags *Flags, commonFlags *utils.CommonFlags) []float64 {
	precision := uint(10_000)
	if flags.radiusCount > 0 {
		radii := make([]float64, flags.radiusCount)
		for i := range flags.radiusCount {
			radii[i] = float64(precision-rand.UintN(precision)) / float64(precision)
		}
		return radii
	}

	return []float64{commonFlags.Radius}
}

func main() {
	commonFlags := &utils.CommonFlags{}
	commonFlags.CommonFlagsSetup(pgrid.GridLinesTotal)
	flags := flagsSetup()
	flag.Parse()
	commonFlags.ParseArgs()

	antNames := getAntNames(flags, commonFlags)
	initialPoints := getInitialPoints(flags, commonFlags)
	radii := getRadii(flags, commonFlags)

	for _, radius := range radii {
		for _, initialPoint := range initialPoints {
			for _, antName := range antNames {
				rFlag := ""
				if !commonFlags.Rectangle.Empty() {
					rFlag = fmt.Sprintf("-r \\'%s/%d\\'",
						commonFlags.Rectangle.String(), commonFlags.ScaleFactor,
					)
				}
				fmt.Printf(
					"-d %s -sc %d -sn %d -sm %d -su %d %s %s__%f__%s__%d\n",
					commonFlags.Dir,
					commonFlags.MinCleanStreak, commonFlags.MaxNoisyDots,
					commonFlags.MinSteps, commonFlags.MinUniq,
					rFlag,
					antName, radius, initialPoint, commonFlags.MaxSteps,
				)
			}
		}
	}
}
