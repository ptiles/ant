package utils

import (
	"errors"
	"github.com/crazy3lf/colorconv"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
)

func FromDegrees(deg int) float64 {
	return float64(deg) * math.Pi / 180.0
}

func GetPalette(steps int) []color.RGBA {
	var palette = make([]color.RGBA, steps)

	for c := range steps {
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

		h, s, _ := colorconv.RGBToHSL(r, g, b)
		r, g, b, _ = colorconv.HSLToRGB(h, s, 0.75)
		palette[c] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return palette
}

func StartCPUProfile(cpuprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

func StopCPUProfile() {
	pprof.StopCPUProfile()
}

func Open(fileName string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", fileName).Run()
	case "windows":
		exec.Command("start", fileName).Run()
	default:
		exec.Command("xdg-open", fileName).Run()
	}
}

func GetRules(antName string) ([]bool, error) {
	limit := uint8(len(antName))
	var nameInvalid = limit < 2
	rules := make([]bool, limit)
	for i, letter := range antName {
		if letter != 'R' && letter != 'r' && letter != 'L' && letter != 'l' {
			nameInvalid = true
			break
		}
		rules[i] = letter == 'R' || letter == 'r'
	}
	if nameInvalid {
		return rules, errors.New("invalid name")
	}
	return rules, nil
}
