package ximage

import (
	"image"
	"os"
	"path"

	"github.com/ptiles/ant/utils/xpng"
)

func SavePNG(img image.Image, fileName string, txt map[string]string) {
	err := os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = xpng.Encode(file, img, txt)
	if err != nil {
		panic(err)
	}
}

func RemoveAlpha(img *image.NRGBA) {
	for y := range img.Rect.Dy() {
		yOffset := y * img.Stride
		for x := range img.Rect.Dx() {
			img.Pix[yOffset+x*4+3] = 255
		}
	}
}
