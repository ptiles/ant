package palette

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math/rand/v2"
)

type Palette []color.RGBA

func GetPaletteRainbow(steps int) Palette {
	var palette = make(Palette, steps)

	for si := range steps {
		h := float64((si*360/steps + 180) % 360)
		c := .95
		l := .95
		r, g, b := colorful.Hcl(h, c, l).Clamped().RGB255()
		palette[si] = color.RGBA{R: r, G: g, B: b, A: 0xff}
	}

	return palette
}

func GetPaletteMonochromatic(steps int, seedString string) Palette {
	var palette = make(Palette, steps)

	rng := rngFromString(seedString)
	h := float64(rng.IntN(360))

	for si := range steps {
		sf := float64(si) / float64(steps)
		c := .95 - 0.3*sf
		l := .65 + 0.3*sf
		r, g, b := colorful.Hcl(h, c, l).Clamped().RGB255()
		palette[si] = color.RGBA{R: r, G: g, B: b, A: 0xff}
	}

	return palette
}

func rngFromString(seedString string) *rand.Rand {
	var seed [32]byte
	copy(seed[:], seedString)

	return rand.New(rand.NewChaCha8(seed))
}
