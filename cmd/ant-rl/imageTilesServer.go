package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/draw"
)

type ImageTile struct {
	updated bool
	img     *image.RGBA
}

type Pair struct {
	X, Y int
}

var imageTiles = map[Pair]ImageTile{}

type PositionedImage struct {
	image    *rl.Image
	position Pair
}

func updateTile(modifiedImage *image.RGBA, point Pair) {
	tileR := image.Rect(0, 0, 256, 256).Add(image.Point{X: point.X << 8, Y: point.Y << 8})

	tile, ok := imageTiles[point]
	if !ok {
		tile.img = image.NewRGBA(tileR)
	}

	draw.Draw(tile.img, tile.img.Rect, modifiedImage, tile.img.Rect.Min, draw.Over)

	imageTiles[point] = ImageTile{updated: true, img: tile.img}
	tile.updated = true
}

func splitViewRect(rect image.Rectangle) []Pair {
	// TODO: prepare constant size array
	pairs := make([]Pair, 0)
	if rect.Min.X < 0 {
		rect.Min.X -= 256
	}
	if rect.Min.Y < 0 {
		rect.Min.Y -= 256
	}
	if rect.Max.X < 0 {
		rect.Max.X -= 256
	}
	if rect.Max.Y < 0 {
		rect.Max.Y -= 256
	}
	for y := rect.Min.Y / 256; y < rect.Max.Y/256+1; y++ {
		for x := rect.Min.X / 256; x < rect.Max.X/256+1; x++ {
			pairs = append(pairs, Pair{x, y})
		}
	}
	return pairs
}

func rlImageFromTileImg(imageTileImg *image.RGBA) *rl.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(img, img.Rect, imageTileImg, imageTileImg.Rect.Min, draw.Over)
	return rl.NewImageFromImage(img)
}

func imageTilesServer(
	modifiedImagesCh <-chan *image.RGBA,
	fullViewRectCh <-chan image.Rectangle,
	positionedImagesCh chan<- *PositionedImage,
) {
	for {
		select {
		case modifiedImage := <-modifiedImagesCh:
			if modifiedImage == nil {
				break
			}

			xx, yy := modifiedImage.Rect.Min.X, modifiedImage.Rect.Min.Y
			x, y := xx>>8, yy>>8

			updateTile(modifiedImage, Pair{X: x, Y: y})
			if xx%256 != 0 {
				updateTile(modifiedImage, Pair{X: x + 1, Y: y})
			}
			if yy%256 != 0 {
				updateTile(modifiedImage, Pair{X: x, Y: y + 1})
			}
			if xx%256 != 0 && yy%256 != 0 {
				updateTile(modifiedImage, Pair{X: x + 1, Y: y + 1})
			}
		default:
		}

		// This waits for next animation frame
		fullViewRect := <-fullViewRectCh

		if fullViewRect.Min.Eq(fullViewRect.Max) {
			imageTiles = map[Pair]ImageTile{}
			continue
		}

		for _, viewRect := range splitViewRect(fullViewRect) {
			imageTile := imageTiles[viewRect]
			if imageTile.updated {
				imageTiles[viewRect] = ImageTile{img: imageTile.img, updated: false}

				rlImage := rlImageFromTileImg(imageTile.img)
				positionedImagesCh <- &PositionedImage{image: rlImage, position: viewRect}
			}
		}
		positionedImagesCh <- nil
	}
}
