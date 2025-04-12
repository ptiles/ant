package utils

import (
	"image"
	"regexp"
	"strconv"
)

func ParseRectangleStr(rectangleStr string) (rect image.Rectangle, scaleFactor int, err error) {
	if rectangleStr == "" {
		return
	}

	re := regexp.MustCompile(`\((-?\d+),(-?\d+)\)-\((-?\d+),(-?\d+)\)/(\d+)`)
	result := re.FindStringSubmatch(rectangleStr)

	x1s, y1s, x2s, y2s, scs := result[1], result[2], result[3], result[4], result[5]

	x1, err := strconv.Atoi(x1s)
	if err != nil {
		return
	}

	y1, err := strconv.Atoi(y1s)
	if err != nil {
		return
	}

	x2, err := strconv.Atoi(x2s)
	if err != nil {
		return
	}

	y2, err := strconv.Atoi(y2s)
	if err != nil {
		return
	}

	sc, err := strconv.Atoi(scs)
	if err != nil {
		return
	}

	return image.Rectangle{Min: image.Point{X: x1, Y: y1}, Max: image.Point{X: x2, Y: y2}}, sc, nil
}
