package main

import (
	"encoding/json"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"os"
	"path"
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

func saveImageFromModifiedImages(modifiedImagesCh <-chan pgrid.ModifiedImage, fileNameFmt string, flags *Flags, commonFlags *utils.CommonFlags) {
	maxDimension := flags.maxDimension
	steps := commonFlags.MaxSteps

	imagesCount := 1
	scaleFactor := 1

	mImg0 := <-modifiedImagesCh
	img0 := mImg0.Img

	activeRectN := img0.Rect
	boundRectN := newBoundFromRect(img0.Rect, maxDimension)
	activeImageS := image.NewRGBA(boundRectN)
	drawImg(activeImageS, img0, scaleFactor)

	for mImg := range modifiedImagesCh {
		img := mImg.Img
		imagesCount += 1
		if mImg.Save {
			saveImage(activeImageS, activeRectN, scaleFactor, fileNameFmt, mImg.Steps)
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

	fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(steps), "png")
	uniqPct := 100 * len(commonFlags.AntName) * pgrid.Uniq() / int(steps)
	dimensions := fmt.Sprintf("%dx%d", activeRectN.Dx(), activeRectN.Dy())
	dimensionsScaled := fmt.Sprintf("%dx%d", activeRectN.Dx()/scaleFactor, activeRectN.Dy()/scaleFactor)
	fmt.Printf(
		"%s %d steps; %s (%s); %d%% uniq\n",
		fileName, steps, dimensionsScaled, dimensions, uniqPct,
	)

	maxSide := activeRectN.Dx()
	if activeRectN.Dx() < activeRectN.Dy() {
		maxSide = activeRectN.Dy()
	}

	if flags.jsonStats {
		writeStats(fileNameFmt, statsType{
			AntName:          commonFlags.AntName,
			FileName:         fileName,
			Steps:            steps,
			UniqPct:          uniqPct,
			ImagesCount:      imagesCount,
			MaxSide:          maxSide,
			Dimensions:       dimensions,
			DimensionsScaled: dimensionsScaled,
		})
	}
}

func saveImage(activeImageS *image.RGBA, activeRectN image.Rectangle, scaleFactor int, fileNameFmt string, steps uint64) {
	fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(steps), "png")

	if steps < 0 {
		fmt.Println(fileName)
	}

	err := os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		panic(err)
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

type statsType struct {
	AntName          string `json:"antName"`
	FileName         string `json:"fileName"`
	Steps            uint64 `json:"steps"`
	UniqPct          int    `json:"uniqPct"`
	ImagesCount      int    `json:"imagesCount"`
	MaxSide          int    `json:"maxSide"`
	Dimensions       string `json:"dimensions"`
	DimensionsScaled string `json:"dimensionsScaled"`
}

func writeStats(fileNameFmt string, stats statsType) {
	fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(stats.Steps), "json")

	err := os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := json.MarshalIndent(&stats, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	file.Write(b)
	file.WriteString("\n")
}
