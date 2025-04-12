package main

import (
	"encoding/json"
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/result"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"image"
	"image/png"
	"os"
	"path"
)

func saveImageFromModifiedImages(modifiedImagesCh <-chan step.ModifiedImage, fileNameFmt string, flags *Flags, commonFlags *utils.CommonFlags) uint64 {
	out := result.NewImage(commonFlags.Rectangle, commonFlags.ScaleFactor, flags.maxDimension)

	stepsTotal := uint64(0)
	minSteps := commonFlags.MaxSteps * commonFlags.MinStepsPct / 100

	imagesCount := 0

	for mImg := range modifiedImagesCh {
		out.Merge(mImg.Img)
		if mImg.Save {
			img, resultRectS := out.Draw(commonFlags.Alpha)
			saveImage(img, resultRectS, out.ResultRectN, out.ScaleFactor, fileNameFmt, mImg.Steps)
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

	img, resultRectS := out.Draw(commonFlags.Alpha)
	if stepsTotal >= minSteps && uniqPct >= commonFlags.MinUniqPct {
		fmt.Printf(saveImage(img, resultRectS, out.ResultRectN, out.ScaleFactor, fileNameFmt, stepsTotal))
	}

	if flags.jsonStats {
		fileName := fmt.Sprintf(fileNameFmt, utils.WithUnderscores(stepsTotal), "png")
		writeStats(fileNameFmt, statsType{
			AntName:          commonFlags.AntName,
			FileName:         fileName,
			Steps:            stepsTotal,
			UniqPct:          uniqPct,
			ImagesCount:      imagesCount,
			MaxSide:          max(out.ResultRectN.Dx(), out.ResultRectN.Dy()),
			Dimensions:       out.ResultRectN.Size().String(),
			DimensionsScaled: out.ResultRectN.Size().Div(out.ScaleFactor).String(),
		})
	}
	return stepsTotal
}

func saveImage(resultImageS *image.NRGBA, resultRectS, resultRectN image.Rectangle, scaleFactor int, fileNameFmt string, steps uint64) string {
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

	err = png.Encode(file, resultImageS)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"%s %s %s %s/%d\n",
		fileName,
		resultRectS.Size().String(),
		resultRectN.Size().String(),
		resultRectN.String(), scaleFactor,
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
