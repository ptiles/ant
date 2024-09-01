package pgrid

import (
	"image"
	"image/color"
	"math"
)

type pointColor struct {
	x, y int
	c    uint8
}
type modifiedPoints map[GridPoint]uint8

const MaxModifiedPoints = 32 * 1024

func (f *Field) ModifiedPointsStepper(modifiedImagesCh chan<- *image.RGBA, maxSteps int, palette []color.RGBA) {
	prevPoint, currPoint, prevLine, currLine, prevPointColor := f.initialState()

	modifiedPointsCh := make(chan *modifiedPoints, 1024)

	go modifiedPointsToImages(modifiedPointsCh, modifiedImagesCh, palette)

	modifiedPointsCh <- &modifiedPoints{prevPoint: prevPointColor}
	modifiedCount := 0
	points := &modifiedPoints{}

	for range maxSteps {
		prevPoint, currPoint, prevLine, currLine, prevPointColor = f.next(prevPoint, currPoint, prevLine, currLine)

		(*points)[prevPoint] = prevPointColor
		if modifiedCount == MaxModifiedPoints {
			modifiedPointsCh <- points
			modifiedCount = 0
			points = &modifiedPoints{}
		}
		modifiedCount += 1
	}
	modifiedPointsCh <- points
	close(modifiedPointsCh)
}

func floor(v int) int {
	return int(math.Floor(float64(v)/256.0)) * 256
}

func ceil(v int) int {
	return int(math.Ceil(float64(v)/256.0)) * 256
}

func modifiedPointsToImages(modifiedPointsCh <-chan *modifiedPoints, modifiedImagesCh chan<- *image.RGBA, palette []color.RGBA) {
	for pointsMap := range modifiedPointsCh {
		i := 0
		rect := image.Rectangle{}
		pixelRect := image.Point{X: 1, Y: 1}
		points := make([]pointColor, len(*pointsMap))

		for prevPoint, c := range *pointsMap {
			centerPoint := prevPoint.getCenterPoint()
			rect = rect.Union(image.Rectangle{
				Min: centerPoint,
				Max: centerPoint.Add(pixelRect),
			})
			points[i] = pointColor{centerPoint.X, centerPoint.Y, c}
			i += 1
		}

		currentImage := image.NewRGBA(image.Rectangle{
			Min: image.Point{X: floor(rect.Min.X), Y: floor(rect.Min.Y)},
			Max: image.Point{X: ceil(rect.Max.X), Y: ceil(rect.Max.Y)},
		})
		for _, point := range points {
			currentImage.Set(point.x, point.y, palette[point.c])
		}
		modifiedImagesCh <- currentImage
	}
	close(modifiedImagesCh)
}
