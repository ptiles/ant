package main

import (
	"flag"
	"fmt"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/wgrid"
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
	short := false

	if flags.fileName == "" {
		fmt.Fprintf(os.Stderr, "fileName -f is required\n")
		short = true
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
		short = true
	}

	if exit {
		flag.PrintDefaults()
		os.Exit(1)
	}

	wg := wgrid.New(flags.rectangle, flags.scaleFactor)
	intersections := wg.IntersectionsMap(flags.gridSize+5, math.MaxInt)

	fmt.Fprintf(os.Stderr, "%3d uniq intersections\n", len(intersections))

	crops := make([]image.Rectangle, 0, len(intersections))

	for point, axes := range intersections {
		axesCount := len(axes)
		fmt.Fprintf(os.Stderr, "[%6d,%6d]x%d  |", point.X, point.Y, axesCount)
		if axesCount >= flags.minAxes {
			crop := cropRect(point, flags.cropSize, flags.cropCenter)
			if crop.In(wg.RectS) {
				crops = append(crops, crop)

				fmt.Fprintf(os.Stderr, "%s\n", axes)
			}
		}
	}

	if short {
		return
	}

	fmt.Fprintf(os.Stderr, "\n%3d uniq fully contained crops\n", len(crops))
	for _, crop := range crops {
		fmt.Fprintf(os.Stderr, "[%6d,%6d]-[%6d,%6d]\n", crop.Min.X, crop.Min.Y, crop.Max.X, crop.Max.Y)
	}

	fmt.Fprintf(os.Stderr, "\n")
	for i, crop := range crops {
		cropStr := fmt.Sprintf("%dx%d+%d+%d!",
			flags.cropSize.X, flags.cropSize.Y,
			crop.Min.X-wg.RectS.Min.X, crop.Min.Y-wg.RectS.Min.Y,
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
