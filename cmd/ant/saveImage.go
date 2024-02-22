package main

import (
	"fmt"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/store"
	"github.com/ptiles/ant/utils"
	"image"
	"image/png"
	"os"
)

func saveImage(field *pgrid.Field, fileName string, limit uint8, steps, minWidth, minHeight int, whiteBackground bool) {
	if steps == 0 {
		return
	}

	maxX, maxY := getMinMax(field, minWidth, minHeight)

	points := 0
	img := image.NewPaletted(image.Rect(0, 0, maxX*2, maxY*2), utils.GetPalette(int(limit), whiteBackground))
	store.ForEach(func(axis0, axis1 uint8, off0, off1 int16, color uint8) {
		line0 := pgrid.GridLine{Axis: axis0, Offset: off0}
		line1 := pgrid.GridLine{Axis: axis1, Offset: off1}
		gp := field.MakeGridPoint(line0, line1)
		point := field.GetCenterPoint(&gp)
		img.SetColorIndex(int(point[0])+maxX, int(point[1])+maxY, color+1)
		points += 1
	})

	fmt.Printf("%s Steps: %d; Points: %d; Size: %dx%d\n", fileName, steps, points, maxX*2, maxY*2)

	// Create a new file to save the PNG image
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image as a PNG and save it to the file
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
