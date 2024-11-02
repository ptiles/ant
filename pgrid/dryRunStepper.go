package pgrid

import (
	"fmt"
	"github.com/ptiles/ant/utils"
	"image"
)

func (f *Field) DryRunStepper(maxSteps, minCleanStreak, maxNoisyDots uint64) {
	modifiedCount := uint64(0)

	stepLen := 1 + len(utils.WithUnderscores(maxSteps))
	stepFormat := fmt.Sprintf("%%%ds", stepLen)

	dotSize := getDotSize(maxSteps)
	dotFormat := fmt.Sprintf("%%%ds", stepLen)
	dotValue := fmt.Sprintf(". = %s", utils.WithUnderscores(dotSize))
	fmt.Printf(dotFormat, dotValue)

	visited := make(map[GridAxes]uint64, max(MaxModifiedPoints, noiseMax, noiseClear))
	stepNumber := uint64(0)
	dotNumber := 0
	noise := uint64(0)

	shouldStop := false
	cleanStreak := uint64(0)
	noisyCount := uint64(0)

	fmt.Printf("\n")
	r := image.Rectangle{}
	one := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 1, Y: 1}}

	//for gridPointAxes := range f.RunAxes(maxSteps) {
	for gridPointAxes, centerPoint := range f.RunPoint(maxSteps) {
		r = r.Union(one.Add(centerPoint))

		if visitedStep, ok := visited[gridPointAxes]; ok {
			stepDiff := stepNumber - visitedStep
			if noiseMin < stepDiff && stepDiff < noiseMax {
				noise += 1
			}
		}

		visited[gridPointAxes] = stepNumber
		if stepNumber%dotSize == 0 {
			if dotNumber%50 == 0 {
				stepNumberPadded := fmt.Sprintf(stepFormat, utils.WithUnderscores(stepNumber))
				if stepNumber > 0 {
					uniq := Uniq()
					ps := float32(stepNumber) / float32(uniq)
					uniqPadded := fmt.Sprintf(stepFormat, utils.WithUnderscores(uniq))
					fmt.Printf(
						"%s /%s = %.2fst/up\n",
						stepNumberPadded, uniqPadded, ps,
					)
				}
				fmt.Printf(stepNumberPadded)
			}
			if dotNumber%10 == 0 {
				fmt.Printf(" ")
			}
			dotNumber += 1

			dotNoise := noiseCharsLen * noise / dotSize
			fmt.Printf("%s", string(noiseChars[dotNoise]))
			noise = 0

			noisyDot := dotNoise > 3
			if noisyDot && cleanStreak > minCleanStreak || noisyCount > maxNoisyDots {
				shouldStop = true
			}
			if noisyDot {
				cleanStreak = 0
				noisyCount += 1
			} else {
				cleanStreak += 1
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
	fmt.Printf("\n%s", r.String())
	//fmt.Printf("\n")
}
