package pgrid

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
)

type pointColor struct {
	gridPoint   GridPoint
	centerPoint image.Point
	color       uint8
}

const MaxModifiedPoints = 32 * 1024

func (f *Field) ModifiedPointsStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	prevPointPoint, currPoint, prevLine, currLine, currPointColor := f.initialState()

	modifiedPointsCh := make(chan []pointColor, 1024)

	go modifiedPointsToImages(modifiedPointsCh, modifiedImagesCh, palette)

	points := make([]pointColor, MaxModifiedPoints)
	points[0] = pointColor{gridPoint: currPoint, color: currPointColor}
	modifiedCount := 1

	for range maxSteps {
		prevPointPoint, currPoint, prevLine, currLine, currPointColor = f.next(prevPointPoint, currPoint, prevLine, currLine)

		if modifiedCount == MaxModifiedPoints {
			modifiedPointsCh <- points
			modifiedCount = 0
			points = make([]pointColor, MaxModifiedPoints)
		}
		points[modifiedCount] = pointColor{gridPoint: currPoint, color: currPointColor}
		modifiedCount += 1
	}
	modifiedPointsCh <- points[:modifiedCount]
	close(modifiedPointsCh)
}

func floorSnap(v int) int {
	return int(math.Floor(float64(v)/256.0)) * 256
}

func ceilSnap(v int) int {
	return int(math.Ceil(float64(v)/256.0)) * 256
}

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

func modifiedPointsToImages(modifiedPointsCh <-chan []pointColor, modifiedImagesCh chan<- *image.RGBA, palette []color.RGBA) {
	for points := range modifiedPointsCh {
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}
		prevPoint := image.Point{}

		for i := range points {
			centerPoint := points[i].gridPoint.getCenterPoint()
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

		currentImage := image.NewRGBA(snapRect(rect))
		for i := range points {
			currentImage.Set(
				points[i].centerPoint.X, points[i].centerPoint.Y,
				palette[points[i].color],
			)
			if drawTiles {
				for _, cornerPoint := range points[i].gridPoint.getCornerPoints() {
					currentImage.Set(
						cornerPoint.X, cornerPoint.Y,
						palette[points[i].color],
					)
				}
			}
		}
		modifiedImagesCh <- currentImage
	}
	close(modifiedImagesCh)
}
