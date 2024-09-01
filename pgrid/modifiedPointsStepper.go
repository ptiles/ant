package pgrid

import (
	"image"
	"image/color"
	"math"
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

func floor(v int) int {
	return int(math.Floor(float64(v)/256.0)) * 256
}

func ceil(v int) int {
	return int(math.Ceil(float64(v)/256.0)) * 256
}

func modifiedPointsToImages(modifiedPointsCh <-chan []pointColor, modifiedImagesCh chan<- *image.RGBA, palette []color.RGBA) {
	for points := range modifiedPointsCh {
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}

		for i := range points {
			points[i].centerPoint = points[i].gridPoint.getCenterPoint()
			rect = rect.Union(image.Rectangle{
				Min: points[i].centerPoint,
				Max: points[i].centerPoint.Add(pixelRect),
			})
		}

		currentImage := image.NewRGBA(image.Rectangle{
			Min: image.Point{X: floor(rect.Min.X), Y: floor(rect.Min.Y)},
			Max: image.Point{X: ceil(rect.Max.X), Y: ceil(rect.Max.Y)},
		})
		for i := range points {
			currentImage.Set(
				points[i].centerPoint.X,
				points[i].centerPoint.Y,
				palette[points[i].color],
			)
		}
		modifiedImagesCh <- currentImage
	}
	close(modifiedImagesCh)
}
