package utils

import (
	"image"
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

func RectGrow(rect image.Rectangle, maxDimension int) image.Rectangle {
	centerPoint := rect.Min.Add(rect.Max).Div(2)
	halfSize := image.Point{X: maxDimension / 2, Y: maxDimension / 2}
	return image.Rectangle{Min: centerPoint.Sub(halfSize), Max: centerPoint.Add(halfSize)}
}

func RectDiv(rect image.Rectangle, scaleFactor int) image.Rectangle {
	if scaleFactor == 0 || rect.Empty() {
		return image.Rectangle{}
	}
	if scaleFactor == 1 {
		return rect
	}
	return image.Rectangle{Min: rect.Min.Div(scaleFactor), Max: rect.Max.Div(scaleFactor)}
}
