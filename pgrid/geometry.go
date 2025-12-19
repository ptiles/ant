package pgrid

import (
	"github.com/ptiles/ant/geom"
	"math"
)

type Geometry [GridLinesTotal][GridLinesTotal]intersection

type intersection struct {
	deltas   intersectionGeometry
	rotation bool
}

type intersectionGeometry [GridLinesTotal - 2]struct {
	zeroZero  float64
	ax0Delta  float64
	ax1Delta  float64
	distDelta float64
	ceilSide  bool
	targetAx  uint8
}

func newGeometry(radius float64) Geometry {
	g := Geometry{}

	g.prepareRotation()
	g.prepareDeltas(radius)
	g.prepareNearestNeighbor(radius)

	return g
}

func (g *Geometry) prepareRotation() {
	for ax0 := range GridLinesTotal {
		for ax1 := range GridLinesTotal {
			if ax0 == ax1 {
				g[ax0][ax1].rotation = true
				continue
			}

			a := axisVector(ax0)
			b := axisVector(ax1)
			perpDotProduct := a.X*b.Y - a.Y*b.X

			g[ax0][ax1].rotation = perpDotProduct > 0
		}
	}
}

func axisVector(ax uint8) geom.Point {
	angle := 360 / float64(GridLinesTotal) * float64(ax)

	return geom.Point{
		X: geom.Cos(angle),
		Y: geom.Sin(angle),
	}
}

func (g *Geometry) prepareDeltas(radius float64) {
	for ax0, ax1 := range AxesAll() {
		for i, axT := range otherAxes(ax0, ax1) {
			deltas := newOffsetDeltas(radius, ax0, ax1, axT)

			g[ax0][ax1].deltas[i].targetAx = axT
			g[ax0][ax1].deltas[i].zeroZero = deltas.zeroZero
			g[ax0][ax1].deltas[i].ax0Delta = deltas.ax0Delta
			g[ax0][ax1].deltas[i].ax1Delta = deltas.ax1Delta
			g[ax0][ax1].deltas[i].distDelta = threeAxesOffset(axT, ax1, ax0)
		}
	}
}

func (g *Geometry) prepareNearestNeighbor(radius float64) {
	for ax0, ax1 := range AxesAll() {
		for i, axT := range otherAxes(ax0, ax1) {
			deltas := newOffsetDeltas(radius, ax1, axT, ax0)

			nextLineOffset := math.Ceil(g[ax0][ax1].deltas[i].zeroZero)
			dist := deltas.zeroZero + deltas.ax1Delta*nextLineOffset
			g[ax0][ax1].deltas[i].ceilSide = dist > 0
		}
	}
}

type offsetDeltas struct {
	targetAx uint8
	zeroZero float64
	ax0Delta float64
	ax1Delta float64
}

func newOffsetDeltas(radius float64, ax0, ax1, axT uint8) offsetDeltas {
	ax0Delta := threeAxesOffset(ax0, ax1, axT)
	ax1Delta := threeAxesOffset(ax1, ax0, axT)

	// when radius is not the same for all axes
	//zeroZero := radii[ax0]*ax0Delta + radii[ax1]*ax1Delta - radii[axT]
	zeroZero := radius * (ax0Delta + ax1Delta - 1)

	return offsetDeltas{zeroZero: zeroZero, ax0Delta: ax0Delta, ax1Delta: ax1Delta}
}

func threeAxesOffset(primaryAx, secondaryAx, targetAx uint8) float64 {
	alpha := 360 / float64(GridLinesTotal)

	primary := float64(primaryAx) * alpha
	secondary := float64(secondaryAx) * alpha
	target := float64(targetAx) * alpha

	if GridLinesTotal == 5 {
		return geom.SinOverSin5(target-secondary, primary-secondary)
	}

	return geom.SinOverSin(target-secondary, primary-secondary)
}
