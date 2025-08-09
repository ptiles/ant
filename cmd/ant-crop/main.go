package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/geom"
	"github.com/ptiles/ant/output"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/seq"
	"github.com/ptiles/ant/utils"
	"image"
	"math"
	"os"
)

type Flags struct {
	fileName    string
	outPrefix   string
	gridSize    int
	minAxes     int
	rectangle   image.Rectangle
	cropSize    image.Point
	cropCenter  image.Point
	scaleFactor int
}

func main() {
	flags := &Flags{}

	flag.StringVar(&flags.fileName, "f", "", "Input file name")
	flag.StringVar(&flags.outPrefix, "p", "crop", "Output files prefix")
	flag.IntVar(&flags.gridSize, "g", 0, "Grid size")
	flag.IntVar(&flags.minAxes, "m", 0, "Min intersecting axes")
	flag.Func("r", "Overall output image rectangle, ex: [0,0]#[5120,5120]*48", func(rectangleStr string) (err error) {
		flags.rectangle, flags.scaleFactor, err = utils.ParseRectangleStr(rectangleStr)
		return
	})
	flag.Func("c", "Cropped image size and center point, ex: [1280,720]#[2560,1440]", func(rectangleStr string) (err error) {
		flags.cropSize, flags.cropCenter, err = utils.ParseCropStr(rectangleStr)
		return
	})

	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	exit := false

	if flags.fileName == "" {
		fmt.Fprintf(os.Stderr, "fileName -f is required\n")
		exit = true
	}

	if flags.gridSize == 0 {
		fmt.Fprintf(os.Stderr, "gridSize -g is required\n")
		exit = true
	}

	if flags.rectangle.Empty() {
		fmt.Fprintf(os.Stderr, "rectangle -r is required\n")
		exit = true
	}

	if flags.cropSize.X == 0 || flags.cropSize.Y == 0 {
		fmt.Fprintf(os.Stderr, "cropSize -c is required\n")
		exit = true
	}

	if exit {
		flag.PrintDefaults()
		os.Exit(1)
	}

	intersections := map[image.Point]int{}
	rectangleS := utils.RectDiv(flags.rectangle, flags.scaleFactor)

	for ax0 := uint8(0); ax0 < pgrid.GridLinesTotal-1; ax0 += 1 {
		minOffset0, maxOffset0 := output.AxisRange(ax0, flags.rectangle)
		for off0 := range seq.WythoffMinMaxColumn(minOffset0, maxOffset0, flags.gridSize+5, math.MaxInt) {
			line0 := output.AxisLine(ax0, off0, float64(flags.scaleFactor))
			for ax1 := ax0 + 1; ax1 < pgrid.GridLinesTotal; ax1 += 1 {
				minOffset1, maxOffset1 := output.AxisRange(ax1, flags.rectangle)
				for off1 := range seq.WythoffMinMaxColumn(minOffset1, maxOffset1, flags.gridSize+5, math.MaxInt) {
					line1 := output.AxisLine(ax1, off1, float64(flags.scaleFactor))
					intersection := geom.Intersection(line0, line1)
					intersectionPoint := image.Point{
						X: int(math.Round(intersection.X)),
						Y: int(math.Round(intersection.Y)),
					}
					if intersectionPoint.In(rectangleS) {
						intersections[intersectionPoint] += 1
					}
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "%3d uniq intersections\n", len(intersections))

	crops := make([]image.Rectangle, 0, len(intersections))

	for point, count := range intersections {
		axes := triangleRoot(count) + 1
		fmt.Fprintf(os.Stderr, "[%6d,%6d]x%d\n", point.X, point.Y, axes)
		if axes >= flags.minAxes {
			crop := cropRect(point, flags.cropSize, flags.cropCenter)
			if crop.In(rectangleS) {
				crops = append(crops, cropRect(point, flags.cropSize, flags.cropCenter))
			}

		}
	}

	fmt.Fprintf(os.Stderr, "\n%3d uniq fully contained crops\n", len(crops))
	for _, crop := range crops {
		fmt.Fprintf(os.Stderr, "[%6d,%6d]-[%6d,%6d]\n", crop.Min.X, crop.Min.Y, crop.Max.X, crop.Max.Y)
	}

	fmt.Fprintf(os.Stderr, "\n")
	for i, crop := range crops {
		cropStr := fmt.Sprintf("%dx%d+%d+%d!",
			flags.cropSize.X, flags.cropSize.Y,
			crop.Min.X-rectangleS.Min.X, crop.Min.Y-rectangleS.Min.Y,
		)
		fmt.Printf("magick %s -crop %-24s %s%d.png\n",
			flags.fileName, cropStr, flags.outPrefix, i,
		)
	}
}

func cropRect(p, size, center image.Point) image.Rectangle {
	minSub, maxAdd := center, size.Sub(center)
	return image.Rectangle{Min: p.Sub(minSub), Max: p.Add(maxAdd)}
}

func triangleRoot(x int) int {
	return int((math.Sqrt(8*float64(x)+1) - 1) / 2)
}
