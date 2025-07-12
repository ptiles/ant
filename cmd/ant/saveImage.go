package main

import (
	"encoding/json"
	"fmt"
	"github.com/ptiles/ant/output"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"image"
	"os"
	"path"
)

func saveImageFromModifiedImages(modifiedImagesCh <-chan step.ModifiedImage, fileNameFmt string, flags *Flags, commonFlags *utils.CommonFlags) uint64 {
	out := output.NewImage(commonFlags.Rectangle, commonFlags.ScaleFactor, flags.maxDimension)

	stepsTotal := uint64(0)
	minSteps := commonFlags.Steps.Max * commonFlags.MinStepsPct / 100

	imagesCount := 0

	for mImg := range modifiedImagesCh {
		out.Merge(mImg.Img)
		if mImg.Save {
			img, resultRectS := out.Draw(commonFlags.Alpha)
			saveImage(img, resultRectS, out.ResultRectN, out.ScaleFactor, fileNameFmt, mImg.Steps, commonFlags.Steps.Max)
		}
		imagesCount += 1
		stepsTotal = mImg.Steps
	}

	uniq, uMaps := pgrid.Uniq()
	uniqPct := uint64(len(commonFlags.AntName)) * uniq * 100 / stepsTotal
	fmt.Printf("%s steps;  %s unique points  (%d%%) in %s maps\n",
		utils.WithSeparatorsSpacePadded(stepsTotal, commonFlags.Steps.Max),
		utils.WithSeparators(uniq), uniqPct, utils.WithSeparators(uint64(uMaps)),
	)

	img, resultRectS := out.Draw(commonFlags.Alpha)
	if stepsTotal >= minSteps && uniqPct >= commonFlags.MinUniqPct {
		fmt.Print(saveImage(img, resultRectS, out.ResultRectN, out.ScaleFactor, fileNameFmt, stepsTotal, commonFlags.Steps.Max))
		if flags.jsonStats {
			fileName := fmt.Sprintf(fileNameFmt, utils.WithSeparators(stepsTotal), "png")
			bounds, sizes, sizeMin, sizeMax := pgrid.GetBounds()

			writeStats(fileNameFmt, statsType{
				AntName:          commonFlags.AntName,
				FileName:         fileName,
				Steps:            stepsTotal,
				UniqPct:          uniqPct,
				ImagesCount:      imagesCount,
				MaxSide:          max(out.ResultRectN.Dx(), out.ResultRectN.Dy()),
				Dimensions:       out.ResultRectN.Size().String(),
				DimensionsScaled: out.ResultRectN.Size().Div(out.ScaleFactor).String(),
				Rect:             out.ResultRectN.String(),
				RectMinX:         out.ResultRectN.Min.X,
				RectMinY:         out.ResultRectN.Min.Y,
				RectMaxX:         out.ResultRectN.Max.X,
				RectMaxY:         out.ResultRectN.Max.Y,
				ScaleFactor:      out.ScaleFactor,
				Bounds:           bounds,
				BoundsSizes:      sizes,
				BoundsSizeMin:    int32(sizeMin),
				BoundsSizeMax:    int32(sizeMax),
			})
		}
	}

	return stepsTotal
}

func saveImage(resultImageS *image.NRGBA, resultRectS, resultRectN image.Rectangle, scaleFactor int, fileNameFmt string, steps, max uint64) string {
	fileName := fmt.Sprintf(fileNameFmt, utils.WithSeparatorsZeroPadded(steps, max), "png")
	utils.SaveImage(fileName, resultImageS)
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
	Rect             string `json:"rect"`
	RectMinX         int    `json:"rectMinX"`
	RectMinY         int    `json:"rectMinY"`
	RectMaxX         int    `json:"rectMaxX"`
	RectMaxY         int    `json:"rectMaxY"`
	ScaleFactor      int    `json:"scaleFactor"`

	Bounds        pgrid.Bounds     `json:"bounds"`
	BoundsSizes   pgrid.BoundsSize `json:"boundsSizes"`
	BoundsSizeMin int32            `json:"boundsSizeMin"`
	BoundsSizeMax int32            `json:"boundsSizeMax"`
}

func writeStats(fileNameFmt string, stats statsType) {
	fileName := fmt.Sprintf(fileNameFmt, utils.WithSeparators(stats.Steps), "json")

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
	fmt.Println(fileName)
}
