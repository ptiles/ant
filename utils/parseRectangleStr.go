package utils

import (
	"fmt"
	"image"
	"math"
	"regexp"
)

func ParseRectangleStr(rectangleStr string) (rect image.Rectangle, scaleFactor int, err error) {
	if rectangleStr == "" {
		return
	}

	exprMinMaxDiv := regexp.MustCompile(
		`\((?P<minX>-?\d+),(?P<minY>-?\d+)\)-\((?P<maxX>-?\d+),(?P<maxY>-?\d+)\)/(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprMinMaxDiv, rectangleStr); matches != nil {
		minPoint := image.Point{X: matches["minX"], Y: matches["minY"]}
		maxPoint := image.Point{X: matches["maxX"], Y: matches["maxY"]}
		return image.Rectangle{Min: minPoint, Max: maxPoint}, matches["scale"], nil
	}

	exprCenterSizeDiv := regexp.MustCompile(
		`\((?P<centerX>-?\d+),(?P<centerY>-?\d+)\)#\((?P<sizeX>-?\d+),(?P<sizeY>-?\d+)\)/(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprCenterSizeDiv, rectangleStr); matches != nil {
		centerPoint := image.Point{X: matches["centerX"], Y: matches["centerY"]}
		halfSizePoint := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}.Div(2)
		minPoint := centerPoint.Sub(halfSizePoint)
		maxPoint := centerPoint.Add(halfSizePoint)
		return image.Rectangle{Min: minPoint, Max: maxPoint}, matches["scale"], nil
	}

	exprCenterSizeMul := regexp.MustCompile(
		`\((?P<centerX>-?\d+),(?P<centerY>-?\d+)\)#\((?P<sizeX>-?\d+),(?P<sizeY>-?\d+)\)\*(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprCenterSizeMul, rectangleStr); matches != nil {
		centerPoint := image.Point{X: matches["centerX"], Y: matches["centerY"]}
		halfSizePoint := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}.Mul(matches["scale"]).Div(2)
		minPoint := centerPoint.Sub(halfSizePoint)
		maxPoint := centerPoint.Add(halfSizePoint)
		return image.Rectangle{Min: minPoint, Max: maxPoint}, matches["scale"], nil
	}

	exprSizeMul := regexp.MustCompile(
		`\((?P<sizeX>-?\d+),(?P<sizeY>-?\d+)\)\*(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprSizeMul, rectangleStr); matches != nil {
		sizePoint := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}.Mul(matches["scale"])
		return image.Rectangle{Min: image.Point{}, Max: sizePoint}, matches["scale"], nil
	}
	return image.Rectangle{}, 0, nil
}

func RectCenteredString(rect image.Rectangle, scaleFactor int) string {
	centerPoint := rect.Min.Add(rect.Max).Div(2)
	sizePoint := rect.Size()

	if scaleFactor == 0 {
		size := max(sizePoint.X, sizePoint.Y)
		m := float64(size / (16 * 1024))
		scaleFactor = 2 << uint(math.Log2(m+1))
	}

	return fmt.Sprintf("%s#%s/%d", centerPoint, sizePoint, scaleFactor)
}
