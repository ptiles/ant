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

func saveImageFromModifiedImages(modifiedImagesCh <-chan pgrid.ModifiedImage, fileNameFmt string, flags *Flags, commonFlags *utils.CommonFlags) uint64 {
	maxDimension := flags.maxDimension
	dynamic := commonFlags.Rectangle.Empty()
	stepsTotal := uint64(0)
	minSteps := commonFlags.MaxSteps * commonFlags.MinStepsPct / 100

	imagesCount := 0
	scaleFactor := 0

	var resultRectN image.Rectangle
	var paddingRectN image.Rectangle
	var resultImageS *image.RGBA

	if dynamic {
		resultImageS = image.NewRGBA(image.Rectangle{})
	} else {
		resultRectN = commonFlags.Rectangle
		scaleFactor = commonFlags.ScaleFactor
		resultImageS = image.NewRGBA(utils.RectDiv(commonFlags.Rectangle, scaleFactor))
	}

	for mImg := range modifiedImagesCh {
		if dynamic {
			resultRectN = resultRectN.Union(mImg.Img.Rect)
			if !resultRectN.In(paddingRectN) {
				if scaleFactor == 0 {
					scaleFactor = 1
				} else {
					scaleFactor *= 2
					maxDimension *= 2
				}
				paddingRectN = utils.RectGrow(resultRectN, maxDimension)
				newResultImageS := image.NewRGBA(utils.RectDiv(paddingRectN, scaleFactor))
				mergeImage(newResultImageS, resultImageS, 2)
				resultImageS = newResultImageS
			}
		}

		mergeImage(resultImageS, mImg.Img, scaleFactor)
		if mImg.Save {
			saveImage(resultImageS, resultRectN, scaleFactor, commonFlags.Alpha, fileNameFmt, mImg.Steps)
		}
		imagesCount += 1
		stepsTotal = mImg.Steps
	}

	uniq, uMaps := pgrid.Uniq()
	uniqPct := uint64(len(commonFlags.AntName)) * uniq * 100 / stepsTotal
	fmt.Printf("%s steps;  %s unique points  (%d%%) in %s maps\n",
		utils.WithUnderscoresPadded(stepsTotal, commonFlags.MaxSteps),
		utils.WithUnderscores(uniq), uniqPct, utils.WithUnderscores(uint64(uMaps)),
	)

	if stepsTotal >= minSteps && uniqPct >= commonFlags.MinUniqPct {
		fmt.Printf(saveImage(resultImageS, resultRectN, scaleFactor, commonFlags.Alpha, fileNameFmt, stepsTotal))
	}

	if flags.jsonStats {
		fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(stepsTotal), "png")
		writeStats(fileNameFmt, statsType{
			AntName:          commonFlags.AntName,
			FileName:         fileName,
			Steps:            stepsTotal,
			UniqPct:          uniqPct,
			ImagesCount:      imagesCount,
			MaxSide:          max(resultRectN.Dx(), resultRectN.Dy()),
			Dimensions:       resultRectN.Size().String(),
			DimensionsScaled: resultRectN.Size().Div(scaleFactor).String(),
		})
	}
	return stepsTotal
}

func cropImage(src *image.RGBA, cropRect image.Rectangle, keepAlpha bool) *image.NRGBA {
	dstRect := image.Rectangle{Min: image.Point{}, Max: cropRect.Size()}
	dstImage := image.NewNRGBA(dstRect)
	draw.Draw(dstImage, dstRect, src, cropRect.Min, draw.Over)

	if !keepAlpha {
		for y := range dstImage.Rect.Dy() {
			yOffset := y * dstImage.Stride
			for x := range dstImage.Rect.Dx() {
				dstImage.Pix[yOffset+x*4+3] = 255
			}
		}
	}

	return dstImage
}

func saveImage(activeImageS *image.RGBA, activeRectN image.Rectangle, scaleFactor int, keepAlpha bool, fileNameFmt string, steps uint64) string {
	fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(steps), "png")
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
	resultImageS := cropImage(activeImageS, activeRectS, keepAlpha)
	err = png.Encode(file, resultImageS)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"%s %s %s %s/%d\n",
		fileName,
		activeRectS.Size().String(),
		activeRectN.Size().String(),
		activeRectN.String(), scaleFactor,
	)
}

type statsType struct {
	AntName          string `json:"antName"`
	FileName         string `json:"fileName"`
	Steps            uint64 `json:"steps"`
	UniqPct          uint64 `json:"uniqPct"`
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
