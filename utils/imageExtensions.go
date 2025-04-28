package utils

import (
	"image"
	"image/png"
	"os"
	"path"
)

func floorSnap(v int) int {
	return (v >> 8) << 8
}

func ceilSnap(v int) int {
	return (v>>8 + 1) << 8
}

func SnapRect(rect image.Rectangle, padding int) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: floorSnap(rect.Min.X - padding), Y: floorSnap(rect.Min.Y - padding)},
		Max: image.Point{X: ceilSnap(rect.Max.X + padding), Y: ceilSnap(rect.Max.Y + padding)},
	}
}

func PointSnap(p image.Point, gridSize int) image.Point {
	return image.Point{X: p.X / gridSize * gridSize, Y: p.Y / gridSize * gridSize}
}

func RectGrow(rect image.Rectangle, maxDimension int) image.Rectangle {
	halfImage := image.Point{X: maxDimension / 2, Y: maxDimension / 2}
	centerPoint := rect.Min.Add(rect.Max).Div(2)
	return image.Rectangle{Min: centerPoint.Sub(halfImage), Max: centerPoint.Add(halfImage)}
}

func RectDiv(rect image.Rectangle, scaleFactor int) image.Rectangle {
	if scaleFactor == 1 {
		return rect
	}
	return image.Rectangle{Min: rect.Min.Div(scaleFactor), Max: rect.Max.Div(scaleFactor)}
}

func PointRect(point image.Point, maxDimension int) image.Rectangle {
	gridSize := 16
	snapped := PointSnap(point, gridSize)
	return RectGrow(image.Rectangle{Min: snapped, Max: snapped.Add(image.Point{X: gridSize, Y: gridSize})}, maxDimension)
}

func IsOutside(point image.Point, rect image.Rectangle) bool {
	return point.X < rect.Min.X+16 || point.Y < rect.Min.Y+16 ||
		point.X > rect.Max.X-16 || point.Y > rect.Max.Y-16
}

func SaveImage(fileName string, img image.Image) {
	err := os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
