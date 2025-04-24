package step

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
)

func DryRunStepper(f *pgrid.Field, maxSteps, maxNoisyDots uint64) {
	modifiedCount := uint64(0)

	dotSize := getDotSize(maxSteps)
	fmt.Printf(
		"%*s dot %s   block %s   row %s",
		1+len(utils.WithSeparators(maxSteps)), "",
		utils.WithSeparators(dotSize),
		utils.WithSeparators(dotSize*10),
		utils.WithSeparators(dotSize*50),
	)

	visited := make(map[pgrid.GridAxes]uint64, max(MaxModifiedPoints, noiseMax, noiseClear))
	stepNumber := uint64(0)
	dotNumber := 0
	noise := uint64(0)

	shouldStop := false
	noisyCount := uint64(0)

	fmt.Printf("\n")

	for gridPointAxes := range f.RunAxes(maxSteps) {
		if visitedStep, ok := visited[gridPointAxes]; ok {
			stepDiff := stepNumber - visitedStep
			if noiseMin < stepDiff && stepDiff < noiseMax {
				noise += 1
			}
		}

		visited[gridPointAxes] = stepNumber
		if stepNumber%dotSize == 0 {
			if dotNumber%50 == 0 {
				fmt.Printf("\n%s", utils.WithSeparatorsSpacePadded(stepNumber, maxSteps))
			}
			if dotNumber%10 == 0 {
				fmt.Printf(" ")
			}
			dotNumber += 1

			dotNoise := noiseCharsLen * noise / dotSize
			fmt.Printf("%c", noiseChars[dotNoise])
			noise = 0

			noisyDot := dotNoise > 3
			if noisyDot {
				noisyCount += 1
			}
			if noisyCount >= maxNoisyDots {
				shouldStop = true
			}
		}

		if modifiedCount == MaxModifiedPoints {
			if stepNumber > noiseClear {
				clearStep := stepNumber - noiseClear
				for k, v := range visited {
					if v < clearStep {
						delete(visited, k)
					}
				}
			}

			modifiedCount = 0
		}

		stepNumber += 1
		modifiedCount += 1

		if shouldStop {
			break
		}
	}

	rect := pgrid.Rect(f)
	fmt.Printf("\n%s", rect.String())
	//fmt.Printf("\n")
}
