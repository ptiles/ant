package utils

import "image"

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

func RectGrow(r image.Rectangle, maxDimension int) image.Rectangle {
	halfImage := image.Point{X: maxDimension / 2, Y: maxDimension / 2}
	centerPoint := r.Min.Add(r.Max).Div(2)
	return image.Rectangle{Min: centerPoint.Sub(halfImage), Max: centerPoint.Add(halfImage)}
}

func RectDiv(r image.Rectangle, scaleFactor int) image.Rectangle {
	if scaleFactor == 1 {
		return r
	}
	return image.Rectangle{Min: r.Min.Div(scaleFactor), Max: r.Max.Div(scaleFactor)}
}

func PointRect(p image.Point, maxDimension int) image.Rectangle {
	gridSize := 16
	snapped := PointSnap(p, gridSize)
	return RectGrow(image.Rectangle{Min: snapped, Max: snapped.Add(image.Point{X: gridSize, Y: gridSize})}, maxDimension)
}

func IsOutside(point image.Point, rect image.Rectangle) bool {
	return point.X < rect.Min.X+16 || point.Y < rect.Min.Y+16 ||
		point.X > rect.Max.X-16 || point.Y > rect.Max.Y-16
}
