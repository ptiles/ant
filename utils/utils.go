package utils

import (
	"image/color"
	"log"
	"math"
	"os"
	"runtime/pprof"
)

func FromDegrees(deg int) float64 {
	return float64(deg) * math.Pi / 180.0
}

func GetPalette(steps int, whiteBackground bool) color.Palette {
	var palette = make(color.Palette, steps+1)
	if whiteBackground {
		palette[0] = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	} else {
		palette[0] = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
	}

	for c := 0; c < steps; c++ {
		step := c * 360 / steps

		ra := step + 0*120 + 90
		ga := step + 1*120 + 90
		ba := step + 2*120 + 90

		rs := (1 + math.Sin(FromDegrees(ra))) / 2
		gs := (1 + math.Sin(FromDegrees(ga))) / 2
		bs := (1 + math.Sin(FromDegrees(ba))) / 2

		r := uint8(math.Round(rs*0xd0 + 0x0f))
		g := uint8(math.Round(gs * 0xb0))
		b := uint8(math.Round(bs*0xb0 + 0x4f))

		palette[c+1] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return palette
}

func StartCPUProfile(cpuprofile *string) {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func StopCPUProfile() {
	pprof.StopCPUProfile()
}
