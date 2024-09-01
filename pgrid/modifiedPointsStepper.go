package pgrid

import (
	"fmt"
	"github.com/StephaneBunel/bresenham"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/time/rate"
	"image"
	"image/color"
	"math"
	"os"
	"time"
)

type gridPointColor struct {
	gridPoint   GridPoint
	centerPoint image.Point
	color       uint8
}

const MaxModifiedPoints = 32 * 1024

func (f *Field) ModifiedPointsStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int64, palette []color.RGBA) {
	bar := progressbar.NewOptions64(maxSteps,
		progressbar.OptionSetWidth(75),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionShowIts(),
	)

	currPoint, currLine, prevLine, prevPointSign, pointColor := f.initialState()

	modifiedPointsCh := make(chan []gridPointColor, 1024)

	go modifiedPointsToImages(modifiedPointsCh, modifiedImagesCh, palette, len(f.Rules), bar)

	points := make([]gridPointColor, MaxModifiedPoints)
	modifiedCount := 0

	for range maxSteps {
		if modifiedCount == MaxModifiedPoints {
			bar.Add(modifiedCount)
			modifiedPointsCh <- points
			modifiedCount = 0
			points = make([]gridPointColor, MaxModifiedPoints)
		}

		currPoint, currLine, prevLine, prevPointSign, pointColor = f.next(currPoint, currLine, prevLine, prevPointSign)
		points[modifiedCount] = gridPointColor{gridPoint: currPoint, color: pointColor}
		modifiedCount += 1
	}
	bar.Add(modifiedCount)
	modifiedPointsCh <- points[:modifiedCount]
	close(modifiedPointsCh)
}

func floorSnap(v int) int {
	return int(math.Floor(float64(v)/256.0)) * 256
}

func ceilSnap(v int) int {
	return int(math.Ceil(float64(v)/256.0)) * 256
}

// const padding = deBruijnScale
const padding = deBruijnScale * 4

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
		fmt.Println("Point is way too far (integer overflow)", centerPoint, prevPoint)
		os.Exit(0)
	}
}

const drawTiles = false

func modifiedPointsToImages(
	modifiedPointsCh <-chan []gridPointColor,
	modifiedImagesCh chan<- *image.RGBA,
	palette []color.RGBA, rulesLength int, bar *progressbar.ProgressBar,
) {
	s := rate.Sometimes{Interval: 1 * time.Second}
	pointsCount := 0
	uniqPointsCount := 0

	for points := range modifiedPointsCh {
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}
		prevPoint := image.Point{}
		uniqPoints := make(map[GridAxes]struct{})

		for i := range points {
			gridPoint := points[i].gridPoint
			uniqPoints[gridPoint.Axes] = struct{}{}
			centerPoint := gridPoint.getCenterPoint()
			if !rect.Empty() {
				overflowCheck(centerPoint, prevPoint)
			}
			rect = rect.Union(image.Rectangle{
				Min: centerPoint,
				Max: centerPoint.Add(pixelRect),
			})
			prevPoint = centerPoint
			points[i].centerPoint = centerPoint
		}

		pointsCount += len(points)
		uniqPointsCount += len(uniqPoints)
		s.Do(func() {
			bar.Describe(fmt.Sprintf("Uniq: %d%%", 100*rulesLength*uniqPointsCount/pointsCount))
			pointsCount = 0
			uniqPointsCount = 0
		})
		currentImage := image.NewRGBA(snapRect(rect))

		if drawTiles {
			for i := range points {
				currentImage.Set(
					points[i].centerPoint.X, points[i].centerPoint.Y,
					palette[points[i].color],
				)
				drawTile(currentImage, points[i].gridPoint, palette[points[i].color])
			}
		} else {
			for i := range points {
				currentImage.Set(
					points[i].centerPoint.X, points[i].centerPoint.Y,
					palette[points[i].color],
				)
			}
		}
		modifiedImagesCh <- currentImage
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
