package main

import (
	"encoding/json"
	"fmt"
	"github.com/ptiles/ant/output"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/step"
	"github.com/ptiles/ant/utils"
	"os"
	"path"
	"strings"
)

func saveImageFromModifiedImages(modifiedImagesCh <-chan step.ModifiedImage, fileNameFmt string, flags *Flags, commonFlags *utils.CommonFlags) uint64 {
	out := output.NewImage(commonFlags.Rectangle, commonFlags.ScaleFactor, flags.maxDimension)

	if flags.gridSize > 0 && flags.gridEmpty && !commonFlags.Rectangle.Empty() {
		gridPrefix := fmt.Sprintf("grid_%d_", flags.gridSize)
		fileName := fmt.Sprintf(fileNameFmt, gridPrefix, utils.WithSeparators(0), "png")
		out.SaveGridOnly(fileName, flags.gridSize, commonFlags.Alpha)
	}

	stepsTotal := uint64(0)
	minSteps := commonFlags.Steps.Max * commonFlags.MinStepsPct / 100

	imagesCount := 0

	for mImg := range modifiedImagesCh {
		out.Merge(mImg.Img)

		if mImg.Save {
			saveImages(
				out, commonFlags.Alpha, flags.gridBoth, flags.gridSize,
				fileNameFmt, mImg.Steps, commonFlags.Steps.Max,
			)
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

	if stepsTotal >= minSteps && uniqPct >= commonFlags.MinUniqPct {
		fmt.Print(saveImages(
			out, commonFlags.Alpha, flags.gridBoth, flags.gridSize,
			fileNameFmt, stepsTotal, commonFlags.Steps.Max,
		))

		if flags.jsonStats {
			fileName := fmt.Sprintf(fileNameFmt, "", utils.WithSeparators(stepsTotal), "png")
			bounds, sizes, sizeMin, sizeMax := pgrid.GetBounds(32)
			resultRectN := out.ResultRectN

			writeStats(fileNameFmt, statsType{
				AntName:          commonFlags.AntName,
				FileName:         fileName,
				Steps:            stepsTotal,
				UniqPct:          uniqPct,
				ImagesCount:      imagesCount,
				MaxSide:          max(resultRectN.Dx(), resultRectN.Dy()),
				Dimensions:       resultRectN.Size().String(),
				DimensionsScaled: resultRectN.Size().Div(out.ScaleFactor).String(),
				Rect:             resultRectN.String(),
				RectMinX:         resultRectN.Min.X,
				RectMinY:         resultRectN.Min.Y,
				RectMaxX:         resultRectN.Max.X,
				RectMaxY:         resultRectN.Max.Y,
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

func saveImages(out *output.Image, keepAlpha, gridBoth bool, gridSize int, fileNameFmt string, steps, max uint64) string {
	var fileName, gridFileName string
	var result strings.Builder
	result.WriteString("\n")

	if gridBoth || gridSize == 0 {
		fileName = fmt.Sprintf(fileNameFmt, "", utils.WithSeparatorsZeroPadded(steps, max), "png")
		fmt.Fprintf(&result, "%s\n", fileName)
	}
	if gridSize > 0 {
		gridPrefix := fmt.Sprintf("grid_%d_", gridSize)
		gridFileName = fmt.Sprintf(fileNameFmt, gridPrefix, utils.WithSeparatorsZeroPadded(steps, max), "png")
		fmt.Fprintf(&result, "%s\n", gridFileName)
	}

	resultRectS := out.SaveImages(fileName, gridFileName, gridSize, keepAlpha)
	resultRectN := out.ResultRectN
	resultRectNFormatted := out.RectCenteredString()

	fmt.Fprintf(&result, "\n%s %s %s\n", resultRectS.Size(), resultRectN.Size(), resultRectNFormatted)

	return result.String()
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
	fileName := fmt.Sprintf(fileNameFmt, "", utils.WithSeparators(stats.Steps), "json")

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
