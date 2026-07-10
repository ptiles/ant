package step

import (
	"fmt"
	"image"
	"maps"
	"os"
	"time"

	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/utils/palette"
	"github.com/ptiles/ant/utils/ximage"
)

type gridTileColor struct {
	gridPoint   pgrid.GridPoint
	centerPoint image.Point
	color       uint8
}

type gridPointColor struct {
	gridAxes    pgrid.GridAxes
	centerPoint image.Point
	color       uint8
}

func (gpc *gridPointColor) String() string {
	return fmt.Sprintf("%s %d", gpc.gridAxes.String(), gpc.color)
}

type void struct{}

const MaxModifiedPoints = 32 * 1024
const noiseMin = 512
const noiseMax = 32 * 1024
const noiseClear = max(MaxModifiedPoints, noiseMax)
const visitedMapSize = noiseClear * 1.25 // 8/6.5 ~= 1.23

func getDotSize(maxSteps uint64) uint64 {
	dr := uint64(100)
	de := maxSteps / 500
	if dr > de {
		return dr
	}

	mi := 0
	m := [...]float32{2, 2.5, 2}
	for {
		next := uint64(float32(dr) * m[mi])
		if next >= de {
			break
		}
		dr = next
		mi = (mi + 1) % len(m)
	}

	return dr
}

var noiseChars = []rune("....▁▂▂▃▃▃▄▄▄▄▅▅▅▅▅▆▆▆▆▆▆▇▇▇▇▇▇▇████████")
var noiseCharsLen = uint64(len(noiseChars))

func ModifiedPointsStepper(
	f *pgrid.Field,
	modifiedImagesCh chan<- ModifiedImage,
	pal palette.Palette,
	steps utils.StepCounts,
	maxNoisyDots uint64,
	minTailSteps uint64,
	minTailSize uint64,
) {
	modifiedPointsCh := make(chan []gridPointColor, 64)

	//if pgrid.DrawTilesAndPoints {
	//	go modifiedPointsToImages(f, modifiedPointsCh, modifiedImagesCh, pal, steps, drawTiles)
	//} else {
	//go modifiedPointsToImages(f, modifiedPointsCh, modifiedImagesCh, pal, steps, drawPoints)
	go modifiedPointsToImages(modifiedPointsCh, modifiedImagesCh, pal, steps, drawPoints)
	//}

	points := make([]gridPointColor, MaxModifiedPoints)
	modifiedCount := uint64(0)
	maxSteps := steps.Max

	dotSize := getDotSize(maxSteps)
	fmt.Printf(
		"%*s dot %s;   block %s;   row %s;  ",
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

	var tail map[pgrid.GridAxes]void

	stepNumber := uint64(0)
	dotNumber := uint64(0)
	noise := uint64(0)

	noisyCount := uint64(0)

	start := time.Now()
	lineSize := float64(dotSize * 50)

	for gridAxes, color := range f.RunAxesColor(maxSteps) {
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
				modifiedPointsCh <- points[:modifiedCount]
				break
			}
		}

		if modifiedCount == MaxModifiedPoints {
			if minTailSteps > 0 {
				if stepNumber <= minTailSize || tail == nil {
					if tail == nil {
						tail = make(map[pgrid.GridAxes]void, minTailSize*5/4)
					}

					for _, point := range points {
						tail[point.gridAxes] = void{}
					}
				}

				if stepNumber > minTailSteps {
					shouldStop := false

					for i, point := range points {
						if _, ok := tail[point.gridAxes]; ok {
							shouldStop = true
							modifiedCount = uint64(i)
							break
						}
					}

					if shouldStop {
						modifiedPointsCh <- points[:modifiedCount]
						break
					}
				}
			}

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

			modifiedPointsCh <- points
			modifiedCount = 0
			points = make([]gridPointColor, MaxModifiedPoints)
		}
		points[modifiedCount] = gridPointColor{
			gridAxes:    gridAxes,
			centerPoint: gridAxes.GetCenterPoint(),
			color:       color,
		}

		stepNumber += 1
		modifiedCount += 1
	}
	fmt.Print(" ")
	fmt.Printf("%s MiB", utils.WithSeparators(utils.MemStatsMB()))
	seconds := time.Since(start).Seconds()
	lineSizeLeft := stepNumber % (dotSize * 50)
	if lineSizeLeft > 0 {
		lineSize = float64(lineSizeLeft)
	}
	fmt.Printf("; %s st/s\n", utils.WithSeparators(uint64(lineSize/seconds)))

	modifiedPointsCh <- points[:modifiedCount]
	close(modifiedPointsCh)
}

const OverflowOffset = 1024

func overflowCheck(centerPoint, prevPoint image.Point) {
	diff := image.Rectangle{Min: centerPoint, Max: prevPoint}.Canon()
	if diff.Dx() > OverflowOffset || diff.Dy() > OverflowOffset {
		fmt.Fprint(os.Stderr,
			"\nAnt went too far (integer overflow)\n",
			centerPoint, prevPoint, "\n", os.Args, "\n\n",
		)
		os.Exit(1)
	}
}

func drawPoints(rect image.Rectangle, points []gridPointColor, pal palette.Palette) *image.RGBA {
	img := image.NewRGBA(ximage.SnapRect(rect, pgrid.Padding))

	for i := range points {
		x, y := points[i].centerPoint.X, points[i].centerPoint.Y
		img.Set(x, y, pal[points[i].color])
	}

	return img
}

func drawTiles(rect image.Rectangle, points []gridTileColor, pal palette.Palette) *image.RGBA {
	img := image.NewRGBA(ximage.SnapRect(rect, pgrid.Padding))

	for i := range points {
		gridPoint := points[i].gridPoint
		color := pal[points[i].color]

		ximage.DrawQuad(img, gridPoint.GetCornerPoints(), color)
	}

	return img
}

func rectIsLarge(rect image.Rectangle) bool {
	rectSize := rect.Size()
	return rectSize.X > 1024 || rectSize.Y > 1024
}

type ModifiedImage struct {
	Img   *image.RGBA
	Steps uint64
	Save  bool
}

func modifiedPointsToImages(
	modifiedPointsCh <-chan []gridPointColor,
	modifiedImagesCh chan<- ModifiedImage,
	pal palette.Palette,
	steps utils.StepCounts,
	drawPointsFn func(rect image.Rectangle, points []gridPointColor, pal palette.Palette) *image.RGBA,
) {
	stepsCount := uint64(0)
	prevPoint := image.Point{}

	for points := range modifiedPointsCh {
		start := 0
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}

		for i := range points {
			centerPoint := points[i].centerPoint
			if !rect.Empty() {
				overflowCheck(centerPoint, prevPoint)
			}
			prevPoint = centerPoint

			rect = rect.Union(image.Rectangle{
				Min: centerPoint,
				Max: centerPoint.Add(pixelRect),
			})

			shouldSave := steps.Inc > 0 && stepsCount > 0 && stepsCount >= steps.Min && stepsCount%steps.Inc == 0
			if rectIsLarge(rect) || shouldSave {
				modifiedImagesCh <- ModifiedImage{
					Steps: stepsCount, Save: shouldSave,
					Img: drawPointsFn(rect, points[start:i], pal),
				}

				rect = image.Rectangle{}
				start = i + 1
			}
			stepsCount += 1
		}

		if start < len(points) {
			modifiedImagesCh <- ModifiedImage{
				Steps: stepsCount, Save: false,
				Img: drawPointsFn(rect, points[start:], pal),
			}
		}
	}
	close(modifiedImagesCh)
}
