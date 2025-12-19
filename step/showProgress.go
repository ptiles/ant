package step

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"time"
)

func showProgress(
	dotNumber, noise, noisyCount *uint64, start *time.Time,
	dotSize, stepNumber, stepsMax uint64, lineSize float64,
) {
	// new row
	if *dotNumber%50 == 0 {
		fmt.Print(" ")

		fmt.Printf("%s MiB", utils.WithSeparators(utils.MemStatsMB()))

		if stepNumber != 0 {
			seconds := time.Since(*start).Seconds()
			fmt.Printf("; %s st/s", utils.WithSeparators(uint64(lineSize/seconds)))
			*start = time.Now()
		}

		fmt.Printf("\n%s", utils.WithSeparatorsSpacePadded(stepNumber, stepsMax))
	}

	// new block
	if *dotNumber%10 == 0 {
		fmt.Print(" ")
	}
	*dotNumber += 1

	// new dot
	dotNoise := noiseCharsLen * *noise / dotSize
	fmt.Printf("%c", noiseChars[dotNoise])
	*noise = 0

	if dotNoise > 3 {
		*noisyCount += 1
	}
}
