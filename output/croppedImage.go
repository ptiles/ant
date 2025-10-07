package output

import (
	"github.com/ptiles/ant/utils/ximage"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"os"
	"path"
)

type croppedImage struct {
	dst         *image.NRGBA
	sourcePoint image.Point
}

func newCropped(i *Image) *croppedImage {
	outputRectS := ximage.RectDiv(i.outputRectN(), i.ScaleFactor)
	croppedRect := image.Rectangle{Min: image.Point{}, Max: outputRectS.Size()}
	dst := image.NewNRGBA(croppedRect)

	return &croppedImage{
		dst:         dst,
		sourcePoint: outputRectS.Min,
	}
}

func (ci *croppedImage) draw(src image.Image) *croppedImage {
	draw.Draw(ci.dst, ci.dst.Rect, src, ci.sourcePoint, draw.Over)
	return ci
}

func (ci *croppedImage) savePNG(fileName string) {
	err := os.MkdirAll(path.Dir(fileName), 0755)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, ci.dst)
	if err != nil {
		panic(err)
	}
}

func (ci *croppedImage) removeAlpha() {
	for y := range ci.dst.Rect.Dy() {
		yOffset := y * ci.dst.Stride
		for x := range ci.dst.Rect.Dx() {
			ci.dst.Pix[yOffset+x*4+3] = 255
		}
	}
}
