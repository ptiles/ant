package utils

import (
	"errors"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
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

func GetPaletteMonochromatic(steps int, rng *rand.Rand) []color.RGBA {
	var palette = make([]color.RGBA, steps)

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

func RngFromString(seedString string) *rand.Rand {
	var seed [32]byte
	copy(seed[:], seedString)

	return rand.New(rand.NewChaCha8(seed))
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
	limit := len(antName)
	if limit < 2 {
		return nil, errors.New("name too short")
	}
	rules := make([]bool, limit)
	for i, letter := range antName {
		if letter != 'R' && letter != 'r' && letter != 'L' && letter != 'l' {
			return nil, errors.New("invalid letters in name")
		}
		rules[i] = letter == 'R' || letter == 'r'
	}
	return rules, nil
}

func WithSeparatorsSpacePadded(num, max uint64) string {
	if max == 0 {
		return WithSeparators(num)
	}

	pad := 1 + len(WithSeparators(max))
	return fmt.Sprintf("%*s", pad, WithSeparators(num))
}

func WithSeparatorsZeroPadded(num, max uint64) string {
	pad := len(fmt.Sprintf("%d", max))
	numStr := fmt.Sprintf("%0*d", pad, num)
	return addSeparators(numStr)
}

func WithSeparators(num uint64) string {
	numStr := fmt.Sprintf("%d", num)
	return addSeparators(numStr)
}

func addSeparators(numStr string) string {
	sl := strings.Split(numStr, "")

	n := 3
	j := (len(sl) + n - 1) / n
	result := make([]string, j)

	for i := len(sl); i > 0; i -= n {
		j -= 1
		result[j] = strings.Join(sl[i-min(n, i):i:i], "")
	}

	return strings.Join(result, "_")
}
