package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"math/rand/v2"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const GRID_LINES_TOTAL = uint(pgrid.GRID_LINES_TOTAL)

func genRandomInitialPoint(min, max int) string {
	dirNames := [2]string{"-", "+"}

	ax1 := rand.UintN(GRID_LINES_TOTAL)
	ax2 := (ax1 + 1 + rand.UintN(GRID_LINES_TOTAL-1)) % GRID_LINES_TOTAL

	dir := rand.IntN(2)

	off1 := rand.IntN(max+1-min) + min
	off2 := rand.IntN(max+1-min) + min

	return fmt.Sprintf("%s%d%s%s%d", pgrid.AxisNames[ax1], off1, dirNames[dir], pgrid.AxisNames[ax2], off2)
}

func numToName(num uint64, bitWidth int) string {
	format := fmt.Sprintf("%%0%ds", bitWidth)
	binary := fmt.Sprintf(format, strconv.FormatUint(num, 2))
	return strings.Replace(strings.Replace(binary, "0", "L", -1), "1", "R", -1)
}

type runModeType int

const (
	InitialPointRange runModeType = iota
	NameRange
	NameAndInitialPointRange
)

type Flags struct {
	initialPointMax   int
	initialPointCount int
	antNameRange      string
}

func flagsSetup() *Flags {
	flags := &Flags{}

	flag.IntVar(&flags.initialPointMax, "im", 0, "Initial point max offset")
	flag.IntVar(&flags.initialPointCount, "ic", 0, "Initial point count")
	flag.StringVar(&flags.antNameRange, "nr", "", "Ant name range MIN-MAX")

	flag.Usage = func() {
		flag.PrintDefaults()
	}

	return flags
}

func main() {
	commonFlags := utils.CommonFlagsSetup(pgrid.GRID_LINES_TOTAL)
	flags := flagsSetup()
	flag.Parse()

	var runMode runModeType
	switch {
	case flags.antNameRange != "" && flags.initialPointMax > 0:
		runMode = NameAndInitialPointRange
	case commonFlags.AntName != "":
		runMode = InitialPointRange
	case flags.antNameRange != "":
		runMode = NameRange
	default:
		fmt.Fprintln(os.Stderr, "Ant name or range required")
		os.Exit(1)
	}

	switch runMode {
	case InitialPointRange:
		// -ir initial point range -ic random points count
		minInitialPointOffset := -flags.initialPointMax
		maxInitialPointOffset := +flags.initialPointMax

		for range flags.initialPointCount {
			randomInitialPoint := genRandomInitialPoint(minInitialPointOffset, maxInitialPointOffset)

			fmt.Printf(
				"-d %s -j %s.%s.%d\n",
				path.Clean(commonFlags.Dir), commonFlags.AntName, randomInitialPoint, commonFlags.MaxSteps,
			)
		}

	case NameRange:
		result := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(flags.antNameRange)
		minBitWidth, _ := strconv.Atoi(result[1])
		maxBitWidth, _ := strconv.Atoi(result[2])

		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			maxNum := uint64(1<<bitWidth) - 1
			for num := uint64(1); num < maxNum; num++ {
				name := numToName(num, bitWidth)

				fmt.Printf(
					"-d %s -j %s.%s.%d\n",
					path.Clean(commonFlags.Dir), name, commonFlags.InitialPoint, commonFlags.MaxSteps,
				)
			}
		}

	case NameAndInitialPointRange:
		result := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(flags.antNameRange)
		minBitWidth, _ := strconv.Atoi(result[1])
		maxBitWidth, _ := strconv.Atoi(result[2])

		// -ir initial point range -ic random points count
		minInitialPointOffset := -flags.initialPointMax
		maxInitialPointOffset := +flags.initialPointMax

		for bitWidth := minBitWidth; bitWidth <= maxBitWidth; bitWidth++ {
			for num := uint64(1); num < 1<<bitWidth-1; num++ {
				name := numToName(num, bitWidth)

				for range flags.initialPointCount {
					randomInitialPoint := genRandomInitialPoint(minInitialPointOffset, maxInitialPointOffset)

					fmt.Printf(
						"-d %s -j %s.%s.%d\n",
						path.Clean(commonFlags.Dir), name, randomInitialPoint, commonFlags.MaxSteps,
					)
				}
			}
		}
	}
}
