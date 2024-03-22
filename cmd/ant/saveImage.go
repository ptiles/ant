package main

import (
	"fmt"
	"github.com/ptiles/ant/store"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"os"
)

func newBoundFromRect(r image.Rectangle, maxDimension int) image.Rectangle {
	halfImage := image.Point{X: maxDimension / 2, Y: maxDimension / 2}
	centerPoint := r.Min.Add(r.Max).Div(2)
	return image.Rectangle{Min: centerPoint.Sub(halfImage), Max: centerPoint.Add(halfImage)}
}

func rectDiv(r image.Rectangle, scaleFactor int) image.Rectangle {
	if scaleFactor == 1 {
		return r
	}
	return image.Rectangle{Min: r.Min.Div(scaleFactor), Max: r.Max.Div(scaleFactor)}
}

func croppedImage(activeImage *image.RGBA, r image.Rectangle) *image.RGBA {
	resultRect := image.Rectangle{Min: image.Point{}, Max: image.Point{X: r.Dx(), Y: r.Dy()}}
	resultImage := image.NewRGBA(resultRect)
	draw.Draw(resultImage, resultRect, activeImage, r.Min, draw.Over)
	return resultImage
}

func drawImg(activeImageS, img *image.RGBA, scaleFactor int) {
	draw.CatmullRom.Scale(activeImageS, rectDiv(img.Rect, scaleFactor), img, img.Rect, draw.Over, nil)
}

func saveImageFromModifiedImages(modifiedImagesCh <-chan *image.RGBA, fileNameFmt string, maxDimension, steps, partialImages int) {
	imagesCount := 0
	scaleFactor := 1

	img0 := <-modifiedImagesCh

	activeRectN := img0.Rect
	boundRectN := newBoundFromRect(img0.Rect, maxDimension)
	activeImageS := image.NewRGBA(boundRectN)
	drawImg(activeImageS, img0, scaleFactor)

	for img := range modifiedImagesCh {
		imagesCount += 1
		if partialImages > 0 && imagesCount%partialImages == 0 {
			saveImage(activeImageS, activeRectN, scaleFactor, fileNameFmt, -imagesCount)
		}
		activeRectN = activeRectN.Union(img.Rect)

		if !activeRectN.In(boundRectN) {
			scaleFactor *= 2
			maxDimension *= 2
			boundRectN = newBoundFromRect(activeRectN, maxDimension)
			newActiveImageS := image.NewRGBA(rectDiv(boundRectN, scaleFactor))
			drawImg(newActiveImageS, activeImageS, 2)
			activeImageS = newActiveImageS
		}

		drawImg(activeImageS, img, scaleFactor)
	}

	saveImage(activeImageS, activeRectN, scaleFactor, fileNameFmt, steps)

	fileName := fmt.Sprintf(fileNameFmt, steps)
	uniqPct := 100 * store.Uniq() / steps
	fmt.Printf(
		"%s %d steps; %d%% uniq; %d images; %dx%d (%dx%d)\n",
		fileName, steps, uniqPct, imagesCount, activeRectN.Dx()/scaleFactor, activeRectN.Dy()/scaleFactor, activeRectN.Dx(), activeRectN.Dy(),
	)
}

func saveImage(activeImageS *image.RGBA, activeRectN image.Rectangle, scaleFactor int, fileNameFmt string, steps int) {
	fileName := fmt.Sprintf(fileNameFmt, steps)

	if steps < 0 {
		fmt.Println(fileName)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	activeRectS := rectDiv(activeRectN, scaleFactor)
	resultImageS := croppedImage(activeImageS, activeRectS)
	err = png.Encode(file, resultImageS)
	if err != nil {
		panic(err)
	}
}
