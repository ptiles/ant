package pgrid

type Point [2]float64

type Line [2]Point

func intersection(line1, line2 Line) Point {
	line1point1, line1point2 := line1[0], line1[1]
	line2point1, line2point2 := line2[0], line2[1]

	x1, y1 := line1point1[X], line1point1[Y]
	x2, y2 := line1point2[X], line1point2[Y]
	x3, y3 := line2point1[X], line2point1[Y]
	x4, y4 := line2point2[X], line2point2[Y]

	dx12 := x1 - x2
	dy12 := y1 - y2
	dx34 := x3 - x4
	dy34 := y3 - y4

	den := dx12*dy34 - dy12*dx34

	return Point{
		((x1*y2-y1*x2)*dx34 - dx12*(x3*y4-y3*x4)) / den,
		((x1*y2-y1*x2)*dy34 - dy12*(x3*y4-y3*x4)) / den,
	}
}

func distance(line Line, point Point) float64 {
	x1, y1 := line[0][X], line[0][Y]
	x2, y2 := line[1][X], line[1][Y]

	x0, y0 := point[X], point[Y]

	x10 := x1 - x0
	y10 := y1 - y0
	x21 := x2 - x1
	y21 := y2 - y1

	return x21*y10 - x10*y21
}
