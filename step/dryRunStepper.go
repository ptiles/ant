package step

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"maps"
	"time"
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

	var visited [pgrid.GridLinesTotal][pgrid.GridLinesTotal]map[pgrid.GridCoords]uint64
	for ax0 := range pgrid.GridLinesTotal {
		for ax1 := range pgrid.GridLinesTotal {
			visited[ax0][ax1] = make(map[pgrid.GridCoords]uint64, visitedMapSize)
		}
	}

	stepNumber := uint64(0)
	dotNumber := uint64(0)
	noise := uint64(0)

	noisyCount := uint64(0)

	start := time.Now()
	lineSize := float64(dotSize * 50)

	for gridAxes := range f.RunAxes(maxSteps) {
		if visitedStep, ok := visited[gridAxes.Axis0][gridAxes.Axis1][gridAxes.Coords]; ok {
			stepDiff := stepNumber - visitedStep
			if noiseMin < stepDiff && stepDiff < noiseMax {
				noise += 1
			}
		}

		visited[gridAxes.Axis0][gridAxes.Axis1][gridAxes.Coords] = stepNumber

		if stepNumber%dotSize == 0 {
			showProgress(
				&dotNumber, &noise, &noisyCount, &start,
				dotSize, stepNumber, maxSteps, lineSize,
			)

			if noisyCount >= maxNoisyDots {
				break
			}
		}

		if modifiedCount == MaxModifiedPoints {
			if stepNumber >= noiseClear {
				clearStep := stepNumber - noiseClear
				for ax0 := range pgrid.GridLinesTotal {
					for ax1 := range pgrid.GridLinesTotal {
						maps.DeleteFunc(visited[ax0][ax1], func(_ pgrid.GridCoords, v uint64) bool {
							return v < clearStep
						})
					}
				}
			}

			modifiedCount = 0
		}

		stepNumber += 1
		modifiedCount += 1
	}

	rect := f.Rect()
	fmt.Printf("\n%s", utils.RectCenteredString(rect, 0))
	//fmt.Printf("\n")
}
