package utils

import (
	"fmt"
	"github.com/ptiles/ant/pgrid/axis"
	"github.com/ptiles/ant/utils/ximage"
	"github.com/ptiles/ant/wgrid"
	"image"
	"regexp"
)

func ParseRectangleStr(rectangleStr string) (rect image.Rectangle, scaleFactor int, err error) {
	if rectangleStr == "" {
		return
	}

	exprCenterPointSizeMul := regexp.MustCompile(
		`[\[(](?P<ax0>[A-Z])(?P<off0>-?\d+)[:+-](?P<ax1>[A-Z])(?P<off1>-?\d+)[)\]]#[\[(](?P<sizeX>-?\d+),(?P<sizeY>-?\d+)[)\]]\*(?P<scale>\d+)`,
	)
	stringMatches := NamedStringMatches(exprCenterPointSizeMul, rectangleStr)
	intMatches := NamedIntMatches(exprCenterPointSizeMul, rectangleStr)
	if stringMatches != nil && intMatches != nil {
		center := wgrid.Intersection(
			axis.Index(stringMatches["ax0"]), intMatches["off0"],
			axis.Index(stringMatches["ax1"]), intMatches["off1"], 1)
		size := image.Point{X: intMatches["sizeX"], Y: intMatches["sizeY"]}.Mul(intMatches["scale"])
		return ximage.RectCenterSize(center, size), intMatches["scale"], nil
	}

	exprMinMaxDiv := regexp.MustCompile(
		`[\[(](?P<minX>-?\d+),(?P<minY>-?\d+)[)\]]-[\[(](?P<maxX>-?\d+),(?P<maxY>-?\d+)[)\]]/(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprMinMaxDiv, rectangleStr); matches != nil {
		minPoint := image.Point{X: matches["minX"], Y: matches["minY"]}
		maxPoint := image.Point{X: matches["maxX"], Y: matches["maxY"]}
		return image.Rectangle{Min: minPoint, Max: maxPoint}, matches["scale"], nil
	}

	exprCenterSizeDiv := regexp.MustCompile(
		`[\[(](?P<centerX>-?\d+),(?P<centerY>-?\d+)[)\]]#[\[(](?P<sizeX>-?\d+),(?P<sizeY>-?\d+)[)\]]/(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprCenterSizeDiv, rectangleStr); matches != nil {
		center := image.Point{X: matches["centerX"], Y: matches["centerY"]}
		size := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}
		return ximage.RectCenterSize(center, size), matches["scale"], nil
	}

	exprCenterSizeMul := regexp.MustCompile(
		`[\[(](?P<centerX>-?\d+),(?P<centerY>-?\d+)[)\]]#[\[(](?P<sizeX>-?\d+),(?P<sizeY>-?\d+)[)\]]\*(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprCenterSizeMul, rectangleStr); matches != nil {
		center := image.Point{X: matches["centerX"], Y: matches["centerY"]}
		size := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}.Mul(matches["scale"])
		return ximage.RectCenterSize(center, size), matches["scale"], nil
	}

	exprSizeMul := regexp.MustCompile(
		`[\[(](?P<sizeX>-?\d+),(?P<sizeY>-?\d+)[)\]]\*(?P<scale>\d+)`,
	)
	if matches := NamedIntMatches(exprSizeMul, rectangleStr); matches != nil {
		size := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}.Mul(matches["scale"])
		return ximage.RectCenterSize(image.Point{}, size), matches["scale"], nil
	}
	return image.Rectangle{}, 0, nil
}

func ParseCropStr(rectangleStr string) (image.Point, image.Point, error) {
	if rectangleStr == "" {
		return image.Point{}, image.Point{}, nil
	}

	exprSizeCenter := regexp.MustCompile(
		`[\[(](?P<centerX>-?\d+),(?P<centerY>-?\d+)[)\]]#[\[(](?P<sizeX>-?\d+),(?P<sizeY>-?\d+)[)\]]`,
	)
	if matches := NamedIntMatches(exprSizeCenter, rectangleStr); matches != nil {
		sizePoint := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}
		centerPoint := image.Point{X: matches["centerX"], Y: matches["centerY"]}
		return sizePoint, centerPoint, nil
	}

	exprSize := regexp.MustCompile(
		`[\[(](?P<sizeX>-?\d+),(?P<sizeY>-?\d+)[)\]]`,
	)
	if matches := NamedIntMatches(exprSize, rectangleStr); matches != nil {
		sizePoint := image.Point{X: matches["sizeX"], Y: matches["sizeY"]}
		centerPoint := sizePoint.Div(2)
		return sizePoint, centerPoint, nil
	}

	return image.Point{}, image.Point{}, nil
}

func RectCenteredString(rect image.Rectangle, scaleFactor int) string {
	if scaleFactor == 0 {
		maxSide := max(rect.Size().X, rect.Size().Y)
		for scaleFactor = 1; maxSide/scaleFactor > 16*1024; scaleFactor *= 2 {
		}
	}

	center := ximage.RectCenter(rect)
	sizeS := rect.Size().Div(scaleFactor)
	return fmt.Sprintf("[%d,%d]#[%d,%d]*%d",
		center.X, center.Y, sizeS.X, sizeS.Y, scaleFactor,
	)
}
