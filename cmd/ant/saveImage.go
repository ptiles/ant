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

func mergeImage(dst, src *image.RGBA, scaleFactor int) {
	dstRect := utils.RectDiv(src.Rect, scaleFactor)
	draw.BiLinear.Scale(dst, dstRect, src, src.Rect, draw.Over, nil)
}

func saveImageFromModifiedImages(modifiedImagesCh <-chan pgrid.ModifiedImage, fileNameFmt string, flags *Flags, commonFlags *utils.CommonFlags) {
	maxDimension := flags.maxDimension
	dynamic := commonFlags.Rectangle.Empty()
	steps := commonFlags.MaxSteps

	imagesCount := 0
	scaleFactor := 1

	var activeRectN image.Rectangle
	var boundRectN image.Rectangle
	var activeImageS *image.RGBA

	if dynamic {
		mImg0 := <-modifiedImagesCh
		img0 := mImg0.Img
		imagesCount += 1

		activeRectN = img0.Rect
		boundRectN = utils.RectGrow(img0.Rect, maxDimension)
		activeImageS = image.NewRGBA(boundRectN)
		mergeImage(activeImageS, img0, scaleFactor)
	} else {
		activeRectN = commonFlags.Rectangle
		scaleFactor = commonFlags.ScaleFactor
		activeImageS = image.NewRGBA(utils.RectDiv(commonFlags.Rectangle, scaleFactor))
	}

	for mImg := range modifiedImagesCh {
		img := mImg.Img
		imagesCount += 1
		if mImg.Save {
			saveImage(activeImageS, activeRectN, scaleFactor, fileNameFmt, mImg.Steps)
		}

		if dynamic {
			activeRectN = activeRectN.Union(img.Rect)
			if !activeRectN.In(boundRectN) {
				scaleFactor *= 2
				maxDimension *= 2
				boundRectN = utils.RectGrow(activeRectN, maxDimension)
				newActiveImageS := image.NewRGBA(utils.RectDiv(boundRectN, scaleFactor))
				mergeImage(newActiveImageS, activeImageS, 2)
				activeImageS = newActiveImageS
			}
		}

		mergeImage(activeImageS, img, scaleFactor)
	}

	saveImage(activeImageS, activeRectN, scaleFactor, fileNameFmt, steps)

	fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(steps), "png")
	uniqPct := 100 * len(commonFlags.AntName) * pgrid.Uniq() / int(steps)
	dimensionsScaled := fmt.Sprintf("%dx%d", activeRectN.Dx()/scaleFactor, activeRectN.Dy()/scaleFactor)
	dimensions := fmt.Sprintf("%dx%d", activeRectN.Dx(), activeRectN.Dy())
	activeRect := fmt.Sprintf("%s/%d", activeRectN.String(), scaleFactor)
	fmt.Printf(
		"%s %s %s %s; %d%% uniq\n",
		fileName, dimensionsScaled, dimensions, activeRect, uniqPct,
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

func cropImage(src *image.RGBA, cropRect image.Rectangle) *image.RGBA {
	dstRect := image.Rectangle{Min: image.Point{}, Max: image.Point{X: cropRect.Dx(), Y: cropRect.Dy()}}
	dstImage := image.NewRGBA(dstRect)
	draw.Draw(dstImage, dstRect, src, cropRect.Min, draw.Over)
	return dstImage
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

	activeRectS := utils.RectDiv(activeRectN, scaleFactor)
	resultImageS := cropImage(activeImageS, activeRectS)
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
