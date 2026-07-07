package palette

import (
	"image/color"
	"math/rand/v2"

	"github.com/lucasb-eyer/go-colorful"
)

type Palette []color.RGBA

func GetPaletteRainbow(steps int) Palette {
	return GetPaletteRainbowCL(steps, .95, .95)
}

func GetPaletteRainbowCL(steps int, c, l float64) Palette {
	var palette = make(Palette, steps)

	for si := range steps {
		h := float64((si*360/steps + 180) % 360)
		palette[si] = hcl2rgb(h, c, l)
	}

	return palette
}

func GetPaletteMonochromatic(steps int, seedString string) Palette {
	rng := rngFromString(seedString)
	h := float64(rng.IntN(360))

	return GetPaletteHCL(steps, h, .95, .65)
}

func GetPaletteHCL(steps int, h, c_, l_ float64) Palette {
	var palette = make(Palette, steps)

	for si := range steps {
		sf := float64(si) / float64(steps)
		c := c_ - 0.3*sf
		l := l_ + 0.3*sf
		palette[si] = hcl2rgb(h, c, l)
	}

	return palette
}

func hcl2rgb(h, c, l float64) color.RGBA {
	r, g, b := colorful.Hcl(h, c, l).Clamped().RGB255()
	return color.RGBA{R: r, G: g, B: b, A: 0xff}
}

func rngFromString(seedString string) *rand.Rand {
	var seed [32]byte
	copy(seed[:], seedString)

	return rand.New(rand.NewChaCha8(seed))
}
