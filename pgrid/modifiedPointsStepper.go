package pgrid

import (
	"fmt"
	"github.com/StephaneBunel/bresenham"
	"github.com/ptiles/ant/utils"
	"image"
	"image/color"
	"math"
	"os"
)

type gridPointColor struct {
	gridPoint   GridPoint
	centerPoint image.Point
	color       uint8
}

func (gpc *gridPointColor) String() string {
	return fmt.Sprintf("%s %d", gpc.gridPoint.Axes.String(), gpc.color)
}

const MaxModifiedPoints = 32 * 1024
const noiseMin = 512
const noiseMax = 32 * 1024
const noiseClear = max(MaxModifiedPoints, noiseMax)

func getDotSize(maxSteps uint64) uint64 {
	dr := uint64(100)
	de := maxSteps / 1000
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

func (f *Field) ModifiedPointsStepper(
	modifiedImagesCh chan<- ModifiedImage,
	palette []color.RGBA,
	maxSteps, partialSteps, minCleanStreak, maxNoisyDots uint64,
) {
	modifiedPointsCh := make(chan []gridPointColor, 64)

	go modifiedPointsToImages(modifiedPointsCh, modifiedImagesCh, palette, partialSteps)

	points := make([]gridPointColor, MaxModifiedPoints)
	modifiedCount := uint64(0)

	stepLen := 1 + len(utils.WithUnderscores(maxSteps))
	stepFormat := fmt.Sprintf("\n%%%ds", stepLen)

	dotSize := getDotSize(maxSteps)
	dotFormat := fmt.Sprintf("%%%ds", stepLen)
	dotValue := fmt.Sprintf(". = %s", utils.WithUnderscores(dotSize))
	fmt.Printf(dotFormat, dotValue)

	var visited [GridLinesTotal * GridLinesTotal]map[GridCoords]uint64
	for a := range visited {
		visited[a] = make(map[GridCoords]uint64, noiseClear)
	}

	stepNumber := uint64(0)
	dotNumber := 0
	noise := uint64(0)

	shouldStop := false
	cleanStreak := uint64(0)
	noisyCount := uint64(0)

	for gridPoint, color := range f.Run(maxSteps) {
		visitedStep, ok := visited[gridPoint.Axes.Axis0*GridLinesTotal+gridPoint.Axes.Axis1][gridPoint.Axes.Coords]
		if ok {
			stepDiff := stepNumber - visitedStep
			if noiseMin < stepDiff && stepDiff < noiseMax {
				noise += 1
			}
		}
		visited[gridPoint.Axes.Axis0*GridLinesTotal+gridPoint.Axes.Axis1][gridPoint.Axes.Coords] = stepNumber

		if stepNumber%dotSize == 0 {
			if dotNumber%50 == 0 {
				fmt.Printf(stepFormat, utils.WithUnderscores(stepNumber))
				ClearZeros()
			}
			if dotNumber%10 == 0 {
				fmt.Print(" ")
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
			if stepNumber >= noiseClear {
				clearStep := stepNumber - noiseClear
				for a := range visited {
					for k, v := range visited[a] {
						if v < clearStep {
							delete(visited[a], k)
						}
					}
				}
			}

			modifiedPointsCh <- points
			modifiedCount = 0
			points = make([]gridPointColor, MaxModifiedPoints)
		}
		points[modifiedCount] = gridPointColor{gridPoint: gridPoint, color: color}

		stepNumber += 1
		modifiedCount += 1

		if shouldStop {
			break
		}
	}
	modifiedPointsCh <- points[:modifiedCount]
	close(modifiedPointsCh)
	fmt.Printf("\n")
}

func floorSnap(v int) int {
	return int(math.Floor(float64(v)/256.0)) * 256
}

func ceilSnap(v int) int {
	return int(math.Ceil(float64(v)/256.0)) * 256
}

func snapRect(rect image.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: floorSnap(rect.Min.X - padding), Y: floorSnap(rect.Min.Y - padding)},
		Max: image.Point{X: ceilSnap(rect.Max.X + padding), Y: ceilSnap(rect.Max.Y + padding)},
	}
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
	img := image.NewRGBA(snapRect(rect))

	if drawTilesAndPoints {
		for i := range points {
			x, y := points[i].centerPoint.X, points[i].centerPoint.Y
			img.Set(x, y, palette[points[i].color])
			if i%9 > 0 {
				img.Set(x+2, y, palette[points[i].color])
			}
			if i%9 > 1 {
				img.Set(x+4, y, palette[points[i].color])
			}
			if i%9 > 2 {
				img.Set(x, y+2, palette[points[i].color])
			}
			if i%9 > 3 {
				img.Set(x+2, y+2, palette[points[i].color])
			}
			if i%9 > 4 {
				img.Set(x+4, y+2, palette[points[i].color])
			}
			if i%9 > 5 {
				img.Set(x, y+4, palette[points[i].color])
			}
			if i%9 > 6 {
				img.Set(x+2, y+4, palette[points[i].color])
			}
			if i%9 > 7 {
				img.Set(x+4, y+4, palette[points[i].color])
			}
			drawTile(img, points[i].gridPoint, palette[points[i].color])
		}
	} else {
		for i := range points {
			x, y := points[i].centerPoint.X, points[i].centerPoint.Y
			img.Set(x, y, palette[points[i].color])
		}
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
	palette []color.RGBA,
	partialSteps uint64,
) {
	stepsCount := uint64(0)
	prevPoint := image.Point{}

	for points := range modifiedPointsCh {
		start := 0
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}

		for i := range points {
			gridPoint := points[i].gridPoint

			centerPoint := gridPoint.getCenterPoint()
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
				mImage.Img = drawPoints(rect, points[start:i], palette)
				modifiedImagesCh <- mImage

				rect = image.Rectangle{}
				start = i + 1
			}
			stepsCount += 1
		}

		mImage := ModifiedImage{Steps: stepsCount}
		mImage.Img = drawPoints(rect, points[start:], palette)
		modifiedImagesCh <- mImage
	}
	close(modifiedImagesCh)
}

func drawTile(currentImage *image.RGBA, gridPoint GridPoint, color color.RGBA) {
	cornerPoints := gridPoint.getCornerPoints()
	p0, p1, p2, p3 := cornerPoints[0], cornerPoints[1], cornerPoints[2], cornerPoints[3]

	currentImage.Set(p0.X, p0.Y, color)
	currentImage.Set(p1.X, p1.Y, color)
	currentImage.Set(p2.X, p2.Y, color)
	currentImage.Set(p3.X, p3.Y, color)

	bresenham.DrawLine(currentImage, p0.X, p0.Y, p1.X, p1.Y, color)
	bresenham.DrawLine(currentImage, p1.X, p1.Y, p2.X, p2.Y, color)
	bresenham.DrawLine(currentImage, p2.X, p2.Y, p3.X, p3.Y, color)
	bresenham.DrawLine(currentImage, p3.X, p3.Y, p0.X, p0.Y, color)

}
