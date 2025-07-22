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
	center := RectCenter(rect)
	size := image.Point{X: maxDimension, Y: maxDimension}
	return RectCenterSize(center, size)
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

func RectCenter(rect image.Rectangle) image.Point {
	return rect.Min.Add(rect.Max).Div(2)
}

func RectCenterSize(center, size image.Point) image.Rectangle {
	halfSize := size.Div(2)
	return image.Rectangle{
		Min: center.Sub(halfSize),
		Max: center.Add(halfSize),
	}
}
