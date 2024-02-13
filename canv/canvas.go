package canv

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/ptiles/ant/geom"
	"math"
	"os"
	"strings"
)

var gridPalette = [...]string{
	"#8fce00",
	"#2986cc",
	"#f44336",
	"#6a329f",
	"#c90076",
	"#8fce0040",
	"#2986cc40",
	"#f4433640",
	"#6a329f40",
	"#c9007640",
	"#8fce0010",
	"#2986cc10",
	"#f4433610",
	"#6a329f10",
	"#c9007610",
}

type Canvas struct {
	file        *os.File
	svg         *svg.SVG
	maxX        int
	maxY        int
	paletteSize int
}

func New(fileName string, maxX, maxY, paletteSize int) Canvas {
	result := Canvas{}
	result.paletteSize = paletteSize
	result.maxX = maxX
	result.maxY = maxY

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	result.file = file
	result.svg = svg.New(file)

	result.svg.Start(result.maxX*2, result.maxY*2)
	styles := make([]string, 0, result.paletteSize+1)
	styles = append(styles, "svg { background-color: black; }")
	for c := 0; c < result.paletteSize; c++ {
		format := ".c%d { fill: color-mix(in srgb, #0000ff, hsl(%ddeg, 30%%, 60%%) 80%%); }"
		degrees := 360 * c / result.paletteSize
		styles = append(styles, fmt.Sprintf(format, c, degrees))
	}
	result.svg.Style("text/css", strings.Join(styles, "\n"))
	result.svg.Rect(1, 1, result.maxX*2-2, result.maxY*2-2, "stroke:#444; stroke-width:1px")
	result.svg.Translate(result.maxX, result.maxY)
	return result
}

func (cf *Canvas) Close() {
	cf.svg.Gend()
	cf.svg.End()

	err := cf.file.Close()
	if err != nil {
		panic(err)
	}
}

func (cf *Canvas) DrawOrigin() {
	style := fmt.Sprintf("stroke:%s; stroke-width:1", gridPalette[6])
	cf.svg.Line(-10, 0, 10, 0, style)
	cf.svg.Line(0, -10, 0, 10, style)
	cf.svg.Line(int(cf.maxX)-50, 0, int(cf.maxX)-40, 0, style)
	cf.svg.Line(0, int(cf.maxY)-50, 0, int(cf.maxY)-40, style)
}

func (cf *Canvas) DrawPoint(point geom.Point, color uint8) {
	class := fmt.Sprintf(`class="c%d"`, color%uint8(cf.paletteSize))
	cf.svg.Circle(int(point[0]), int(point[1]), 1, class)
}

func (cf *Canvas) drawLineSegment(line geom.Line, color int) {
	point1, point2 := line[0], line[1]
	style := fmt.Sprintf("stroke:%s; stroke-width:1", gridPalette[color])
	cf.svg.Line(int(point1[0]), int(point1[1]), int(point2[0]), int(point2[1]), style)
}

func borderIntersection(line geom.Line, canvasWidth float64, canvasHeight float64) geom.Line {
	var result geom.Line

	topLeft := geom.Point{-canvasWidth, canvasHeight}
	topRight := geom.Point{canvasWidth, canvasHeight}
	bottomRight := geom.Point{canvasWidth, -canvasHeight}
	bottomLeft := geom.Point{-canvasWidth, -canvasHeight}

	topLine := geom.Line{topLeft, topRight}
	rightLine := geom.Line{topRight, bottomRight}
	bottomLine := geom.Line{bottomRight, bottomLeft}
	leftLine := geom.Line{bottomLeft, topLeft}

	found := 0

	top := geom.Intersection(topLine, line)
	if -canvasWidth < top[0] && top[0] < canvasWidth {
		result[found] = top
		found++
	}

	right := geom.Intersection(rightLine, line)
	if -canvasHeight < right[1] && right[1] < canvasHeight {
		result[found] = right
		found++
	}
	if found == 2 {
		return result
	}

	bottom := geom.Intersection(bottomLine, line)
	if -canvasWidth < bottom[0] && bottom[0] < canvasWidth {
		result[found] = bottom
		found++
	}
	if found == 2 {
		return result
	}

	left := geom.Intersection(leftLine, line)
	if -canvasHeight < left[1] && left[1] < canvasHeight {
		result[found] = left
		found++
	}
	return result
}

func (cf *Canvas) DrawLine(line geom.Line, color int) {
	bi := borderIntersection(line, float64(cf.maxX), float64(cf.maxY))
	cf.drawLineSegment(bi, color)
}

func (cf *Canvas) IsOutside(point geom.Point) bool {
	return math.Abs(point[0]) > float64(cf.maxX) || math.Abs(point[1]) > float64(cf.maxY)
}
