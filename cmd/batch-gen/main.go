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

func genRandomPoint(min, max int) string {
	dirNames := [2]string{"-", "+"}

	ax1 := rand.UintN(GridLinesTotal)
	ax2 := rand.UintN(GridLinesTotal - 1)
	if ax2 == ax1 {
		ax2 = GridLinesTotal - 1
	}

	dir := dirNames[rand.IntN(2)]

	off1 := rand.IntN(max+1-min) + min
	off2 := rand.IntN(max+1-min) + min

	return fmt.Sprintf("%s%d%s%s%d", pgrid.AxisNames[ax1], off1, dir, pgrid.AxisNames[ax2], off2)
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

	initialPointCount    int
	initialPointMax      int
	initialPointRelative bool

	radiusCount int
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.antNameRange, "nr", "", "Ant name range MIN-MAX")

	flag.IntVar(&flags.initialPointCount, "ic", 0, "Initial point count")
	flag.IntVar(&flags.initialPointMax, "im", 8*1024, "Initial point max offset")
	flag.BoolVar(&flags.initialPointRelative, "ir", false, "Initial point relative to -i value")

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

func main() {
	commonFlags := &utils.CommonFlags{}
	commonFlags.CommonFlagsSetup(pgrid.GridLinesTotal)
	flags := flagsSetup()
	flag.Parse()
	commonFlags.ParseArgs()

	var antNames []string
	if flags.antNameRange != "" {
		minBitWidth, maxBitWidth := parseNameRange(flags.antNameRange)

		// TODO: calculate antNames length without loop
		antNamesLength := 0
		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			antNamesLength += 1<<bitWidth - 2
		}

		antNames = make([]string, antNamesLength)

		i := 0
		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			maxNum := uint64(1<<bitWidth) - 1
			for num := uint64(1); num < maxNum; num++ {
				antNames[i] = numToName(num, bitWidth)
				i += 1
			}
		}
	} else if commonFlags.AntName != "" {
		antNames = strings.Split(commonFlags.AntName, ",")
	} else {
		fmt.Fprintln(os.Stderr, "Ant name or range required")
		os.Exit(1)
	}

	var initialPoints []string
	if flags.initialPointRelative && flags.initialPointCount > 0 {
		ax1, off1, _, ax2, off2 := utils.ParseInitialPoint(commonFlags.InitialPoint)

		min1 := off1 - flags.initialPointMax
		max1 := off1 + flags.initialPointMax
		min2 := off2 - flags.initialPointMax
		max2 := off2 + flags.initialPointMax

		initialPoints = make([]string, flags.initialPointCount)
		for i := range flags.initialPointCount {
			initialPoints[i] = genRandomPointAround(ax1, min1, max1, ax2, min2, max2)
		}
	} else if flags.initialPointCount > 0 {
		minInitialPointOffset := -flags.initialPointMax
		maxInitialPointOffset := +flags.initialPointMax

		initialPoints = make([]string, flags.initialPointCount)
		for i := range flags.initialPointCount {
			initialPoints[i] = genRandomPoint(minInitialPointOffset, maxInitialPointOffset)
		}
	} else {
		initialPoints = strings.Split(commonFlags.InitialPoint, ",")
	}

	precision := uint(10_000)
	var radii []float64
	if flags.radiusCount > 0 {
		radii = make([]float64, flags.radiusCount)
		for i := range flags.radiusCount {
			radii[i] = float64(precision-rand.UintN(precision)) / float64(precision)
		}
	} else {
		radii = []float64{0.5}
	}

	for _, radius := range radii {
		for _, initialPoint := range initialPoints {
			for _, antName := range antNames {
				if commonFlags.Rectangle.Empty() {
					fmt.Printf(
						//"-d %s -j %s__%f__%s__%d\n",
						"-d %s %s__%f__%s__%d\n",
						commonFlags.Dir, antName, radius, initialPoint, commonFlags.MaxSteps,
					)
				} else {
					fmt.Printf(
						//"-d %s -j %s__%f__%s__%d\n",
						"-d %s -r \\'%s/%d\\' %s__%f__%s__%d\n",
						commonFlags.Dir, commonFlags.Rectangle.String(), commonFlags.ScaleFactor,
						antName, radius, initialPoint, commonFlags.MaxSteps,
					)
				}
			}
		}
	}
}
