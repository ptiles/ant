package step

import (
	"fmt"
	"github.com/StephaneBunel/bresenham"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"image"
	"image/color"
	"os"
	"time"
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
	palette []color.RGBA,
	maxSteps, partialSteps, maxNoisyDots uint64,
) {
	modifiedPointsCh := make(chan []gridPointColor, 64)

	//if pgrid.DrawTilesAndPoints {
	go modifiedPointsToImages(f, modifiedPointsCh, modifiedImagesCh, palette, partialSteps, drawPoints)
	//} else {
	//	go modifiedPointsToImages(f, modifiedPointsCh, modifiedImagesCh, palette, partialSteps, drawTiles)
	//}

	points := make([]gridPointColor, MaxModifiedPoints)
	modifiedCount := uint64(0)

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

	stepNumber := uint64(0)
	dotNumber := 0
	noise := uint64(0)

	shouldStop := false
	noisyCount := uint64(0)

	start := time.Now()
	lineSize := float64(dotSize * 50)

	for gridAxes, color := range f.RunAxesColor(maxSteps) {
		visitedStep, ok := visited[gridAxes.Axis0][gridAxes.Axis1][gridAxes.Coords]
		if ok {
			stepDiff := stepNumber - visitedStep
			if noiseMin < stepDiff && stepDiff < noiseMax {
				noise += 1
			}
		}
		visited[gridAxes.Axis0][gridAxes.Axis1][gridAxes.Coords] = stepNumber

		if stepNumber%dotSize == 0 {
			// new row
			if dotNumber%50 == 0 {
				fmt.Print(" ")

				fmt.Printf("%s MiB", utils.WithSeparators(utils.MemStatsMB()))

				if stepNumber != 0 {
					seconds := time.Since(start).Seconds()
					fmt.Printf("; %s st/s", utils.WithSeparators(uint64(lineSize/seconds)))
					start = time.Now()
				}

				fmt.Printf("\n%s", utils.WithSeparatorsSpacePadded(stepNumber, maxSteps))
			}

			// new block
			if dotNumber%10 == 0 {
				fmt.Print(" ")
			}
			dotNumber += 1

			// new dot
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
			if stepNumber >= noiseClear {
				clearStep := stepNumber - noiseClear
				for ax0 := range pgrid.GridLinesTotal {
					for ax1 := range pgrid.GridLinesTotal {
						for k, v := range visited[ax0][ax1] {
							if v < clearStep {
								delete(visited[ax0][ax1], k)
							}
						}
					}
				}
			}

			modifiedPointsCh <- points
			modifiedCount = 0
			points = make([]gridPointColor, MaxModifiedPoints)
		}
		points[modifiedCount] = gridPointColor{gridAxes: gridAxes, color: color}

		stepNumber += 1
		modifiedCount += 1

		if shouldStop {
			break
		}
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
		fmt.Println("\nAnt went too far (integer overflow)", centerPoint, prevPoint)
		os.Exit(0)
	}
}

func drawPoints(rect image.Rectangle, points []gridPointColor, palette []color.RGBA) *image.RGBA {
	img := image.NewRGBA(utils.SnapRect(rect, pgrid.Padding))

	for i := range points {
		x, y := points[i].centerPoint.X, points[i].centerPoint.Y
		img.Set(x, y, palette[points[i].color])
	}

	return img
}

func drawTiles(rect image.Rectangle, points []gridTileColor, palette []color.RGBA) *image.RGBA {
	img := image.NewRGBA(utils.SnapRect(rect, pgrid.Padding))

	for i := range points {
		gridPoint := points[i].gridPoint
		color := palette[points[i].color]
		cornerPoints := gridPoint.GetCornerPoints()
		p0, p1, p2, p3 := cornerPoints[0], cornerPoints[1], cornerPoints[2], cornerPoints[3]

		bresenham.DrawLine(img, p0.X, p0.Y, p1.X, p1.Y, color)
		bresenham.DrawLine(img, p1.X, p1.Y, p2.X, p2.Y, color)
		bresenham.DrawLine(img, p2.X, p2.Y, p3.X, p3.Y, color)
		bresenham.DrawLine(img, p3.X, p3.Y, p0.X, p0.Y, color)
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
	f *pgrid.Field,
	modifiedPointsCh <-chan []gridPointColor,
	modifiedImagesCh chan<- ModifiedImage,
	palette []color.RGBA,
	partialSteps uint64,
	drawPointsFn func(rect image.Rectangle, points []gridPointColor, palette []color.RGBA) *image.RGBA,
) {
	stepsCount := uint64(0)
	prevPoint := image.Point{}

	for points := range modifiedPointsCh {
		start := 0
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}

		for i := range points {
			centerPoint := f.GetCenterPoint(points[i].gridAxes)
			if !rect.Empty() {
				overflowCheck(centerPoint, prevPoint)
			}
			prevPoint = centerPoint

			points[i].centerPoint = centerPoint

			rect = rect.Union(image.Rectangle{
				Min: centerPoint,
				Max: centerPoint.Add(pixelRect),
			})

			shouldSave := partialSteps > 0 && stepsCount > 0 && stepsCount%partialSteps == 0
			if rectIsLarge(rect) || shouldSave {
				mImage := ModifiedImage{Steps: stepsCount, Save: shouldSave}
				mImage.Img = drawPointsFn(rect, points[start:i], palette)
				modifiedImagesCh <- mImage

				rect = image.Rectangle{}
				start = i + 1
			}
			stepsCount += 1
		}

		mImage := ModifiedImage{Steps: stepsCount}
		mImage.Img = drawPointsFn(rect, points[start:], palette)
		modifiedImagesCh <- mImage
	}
	close(modifiedImagesCh)
}
