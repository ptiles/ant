package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func saveImageFromModifiedImages(modifiedImages <-chan *image.RGBA, fileName string, steps int) {

	images := make([]*image.RGBA, 0, 1024)

	firstImage := <-modifiedImages
	rect := firstImage.Rect
	images = append(images, firstImage)

	for modifiedImage := range modifiedImages {
		rect = modifiedImage.Rect.Union(rect)
		images = append(images, modifiedImage)
	}

	resultImage := image.NewRGBA(rect)

	for _, img := range images {
		draw.Draw(resultImage, img.Rect, img, img.Rect.Min, draw.Over)
	}

	fmt.Printf("%s Steps: %d; Images: %d; Size: %dx%d\n", fileName, steps, len(images), rect.Dx(), rect.Dy())

	// Create a new file to save the PNG image
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image as a PNG and save it to the file
	err = png.Encode(file, resultImage)
	if err != nil {
		panic(err)
	}
}
