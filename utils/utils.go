package utils

import (
	"errors"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
)

func GetPaletteRainbow(steps int) []color.RGBA {
	var palette = make([]color.RGBA, steps)

	for si := range steps {
		h := float64((si*360/steps + 180) % 360)
		c := .95
		l := .95
		r, g, b := colorful.Hcl(h, c, l).Clamped().RGB255()
		palette[si] = color.RGBA{R: r, G: g, B: b, A: 0xff}
	}

	return palette
}

func GetPaletteMonochromatic(steps int, initialPoint string) []color.RGBA {
	var palette = make([]color.RGBA, steps)

	seed := [32]byte{}
	copy(seed[:], initialPoint)
	rng := rand.New(rand.NewChaCha8(seed))

	h := float64(rng.IntN(360))
	for si := range steps {
		sf := float64(si) / float64(steps-1)
		c := .95 - 0.3*sf
		l := .65 + 0.3*sf
		r, g, b := colorful.Hcl(h, c, l).Clamped().RGB255()
		palette[si] = color.RGBA{R: r, G: g, B: b, A: 0xff}
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

func WithUnderscores(num uint64) string {
	sl := strings.Split(strconv.FormatUint(num, 10), "")

	n := 3
	j := (len(sl) + n - 1) / n
	result := make([]string, j)

	for i := len(sl); i > 0; i -= n {
		j -= 1
		result[j] = strings.Join(sl[i-min(n, i):i:i], "")
	}

	return strings.Join(result, "_")
}
