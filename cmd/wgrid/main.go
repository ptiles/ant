package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"

	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/pgrid/parse"
	"github.com/ptiles/ant/utils"
	"github.com/ptiles/ant/utils/ximage"
	"github.com/ptiles/ant/wgrid"
)

type Flags struct {
	alpha       bool
	dir         string
	gridSize    int
	in          bool
	open        bool
	mark        bool
	rectangle   image.Rectangle
	scaleFactor int
	verbose     bool
}

func main() {
	flags := &Flags{}

	flag.BoolVar(&flags.alpha, "alpha", false, "Save transparent image with alpha channel")
	flag.StringVar(&flags.dir, "d", fmt.Sprintf("results%d", pgrid.GridLinesTotal), "Results directory")
	flag.IntVar(&flags.gridSize, "g", 0, "Grid size")
	flag.BoolVar(&flags.in, "i", false, "Read intersections from stdin")
	flag.BoolVar(&flags.open, "o", false, "Open resulting image")
	flag.BoolVar(&flags.mark, "m", false, "Mark intersections")
	flag.Func("r", "Overall output image rectangle, ex: [0,0]#[5120,5120]*48", func(rectangleStr string) (err error) {
		flags.rectangle, flags.scaleFactor, err = utils.ParseRectangleStr(rectangleStr)
		return
	})
	flag.BoolVar(&flags.verbose, "v", false, "Verbose output")

	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	exit := false

	if flags.gridSize == 0 {
		fmt.Fprintf(os.Stderr, "gridSize -g is required\n")
		exit = true
	}

	if flags.rectangle.Empty() {
		fmt.Fprintf(os.Stderr, "rectangle -r is required\n")
		exit = true
	}

	if exit {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rectStr := fmt.Sprintf("%d", flags.scaleFactor)

	alphaStr := ""
	if flags.alpha {
		alphaStr = "_alpha"
	}

	fileName := fmt.Sprintf("%s/grid_%s_%d%s.png", flags.dir, rectStr, flags.gridSize, alphaStr)
	halfSize, thickness, red := 10, 5, color.NRGBA{R: 255, G: 40, B: 40, A: 128}

	wg := wgrid.New(flags.rectangle)
	rectS := ximage.RectDiv(flags.rectangle, flags.scaleFactor)
	gridImage := image.NewNRGBA(rectS)
	DrawMultiGrid(&wg, gridImage, flags.gridSize)
	intersections := wg.IntersectionsMap(flags.gridSize+3, math.MaxInt)

	if flags.mark {
		for point, axes := range intersections {
			if flags.verbose {
				fmt.Fprintf(os.Stderr, "[%6d,%6d]x%d  |%s\n", point.X, point.Y, len(axes), axes)
			}
			ximage.DrawSquareThick(gridImage, point, halfSize, thickness, red)
		}
		if flags.verbose {
			fmt.Fprintf(os.Stderr, "%3d uniq intersections\t%s\n", len(intersections), rectS)
		}
	}

	if flags.in {
		green := color.NRGBA{R: 40, G: 255, B: 40, A: 128}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ax1, off1, _, ax2, off2 := parse.InitialPoint(scanner.Text())
			point := wgrid.Intersection(ax1, off1, ax2, off2)
			ximage.DrawSquareThick(gridImage, point, halfSize, thickness, green)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}

	if !flags.alpha {
		ximage.RemoveAlpha(gridImage)
	}

	txt := map[string]string{
		"rectangle":   flags.rectangle.String(),
		"scaleFactor": fmt.Sprintf("%d", flags.scaleFactor),

		"-g": strconv.Itoa(flags.gridSize),
		"-r": utils.RectCenteredString(flags.rectangle, flags.scaleFactor),
	}

	ximage.SavePNG(gridImage, fileName, txt)

	fmt.Println("Saved grid image to:", fileName)

	if flags.open {
		utils.Open(fileName)
	}
}
