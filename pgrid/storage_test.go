package pgrid

import (
	"testing"
)

func height(ua *upArray) upInt {
	return ua.Max.Offset0 - ua.Min.Offset0
}

func width(ua *upArray) upInt {
	return ua.Max.Offset1 - ua.Min.Offset1
}

func (gcu *gridCoordsUp) equals(oth gridCoordsUp) bool {
	return gcu.Offset0 == oth.Offset0 && gcu.Offset1 == oth.Offset1
}

func TestInitialize(t *testing.T) {
	initial := divUp(GridCoords{Offset0: -50, Offset1: 50})

	subj := upArray{}
	subj.Initialize(initial)

	//expectedMin := gridCoordsUp{Offset0: -256, Offset1: 0}
	//expectedMax := gridCoordsUp{Offset0: 0, Offset1: 256}
	//expectedSize := upInt(256)
	expectedMin := gridCoordsUp{Offset0: -256, Offset1: -256}
	expectedMax := gridCoordsUp{Offset0: 256, Offset1: 256}
	expectedSize := upInt(512)

	if width(&subj) != expectedSize {
		t.Errorf("Expected width to equal %d, got %d", expectedSize, width(&subj))
	}
	if height(&subj) != expectedSize {
		t.Errorf("Expected height to equal %d, got %d", expectedSize, width(&subj))
	}
	if !subj.Min.equals(expectedMin) {
		t.Errorf("Expected min %s, got %s Size: %dx%d", expectedMin.String(), subj.Min.String(), width(&subj), height(&subj))
	}
	if !subj.Max.equals(expectedMax) {
		t.Errorf("Expected max %s, got %s Size: %dx%d", expectedMax.String(), subj.Max.String(), width(&subj), height(&subj))
	}
}

func up1(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0, Offset1: gc.Offset1 - 1}
}
func up(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0, Offset1: gc.Offset1 - 256}
}
func right1(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0 + 1, Offset1: gc.Offset1}
}
func right(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0 + 256, Offset1: gc.Offset1}
}
func down1(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0, Offset1: gc.Offset1 + 1}
}
func down(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0, Offset1: gc.Offset1 + 256}
}
func left1(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0 - 1, Offset1: gc.Offset1}
}
func left(gc gridCoordsUp) gridCoordsUp {
	return gridCoordsUp{Offset0: gc.Offset0 - 256, Offset1: gc.Offset1}
}

func TestResize(t *testing.T) {
	initialPoint := gridCoordsUp{Offset0: -128, Offset1: 127}

	expectedMin := gridCoordsUp{Offset0: -256, Offset1: 0}
	expectedMax := gridCoordsUp{Offset0: 0, Offset1: 256}
	expectedMinPoint := gridCoordsUp{Offset0: -256, Offset1: 0}
	expectedMaxPoint := gridCoordsUp{Offset0: -1, Offset1: 255}

	centerPoint := gridCoordsUp{
		Offset0: (expectedMin.Offset0 + expectedMax.Offset0) / 2,
		Offset1: (expectedMin.Offset1 + expectedMax.Offset1) / 2,
	}
	cornerNW := expectedMinPoint
	cornerNE := gridCoordsUp{Offset0: expectedMaxPoint.Offset0, Offset1: expectedMinPoint.Offset1}
	cornerSE := expectedMaxPoint
	cornerSW := gridCoordsUp{Offset0: expectedMinPoint.Offset0, Offset1: expectedMaxPoint.Offset1}

	var testsKeep = []struct {
		name  string
		point gridCoordsUp
	}{
		{name: "Should not resize if in the middle", point: centerPoint},
		{name: "Should not resize if initialPoint", point: initialPoint},
		{name: "Should not resize if on NW corner", point: cornerNW},
		{name: "Should not resize if on NE corner", point: cornerNE},
		{name: "Should not resize if on SE corner", point: cornerSE},
		{name: "Should not resize if on SW corner", point: cornerSW},
	}

	for _, tt := range testsKeep {
		t.Run(tt.name, func(t *testing.T) {
			subj := upArray{}
			subj.Initialize(initialPoint)
			subj.ResizeIfNeeded(tt.point)

			if !subj.Min.equals(expectedMin) {
				t.Errorf("Expected min %s, got %s Size: %dx%d", expectedMin.String(), subj.Min.String(), width(&subj), height(&subj))
			}

			if !subj.Max.equals(expectedMax) {
				t.Errorf("Expected max %s, got %s Size: %dx%d", expectedMax.String(), subj.Max.String(), width(&subj), height(&subj))
			}
		})
	}

	var testsResize = []struct {
		name   string
		point  gridCoordsUp
		newMin gridCoordsUp
		newMax gridCoordsUp
	}{
		{
			name:   "Should resize if outside left of NW corner",
			point:  left1(cornerNW),
			newMin: left(up(expectedMin)),
			newMax: expectedMax,
		},
		{
			name:   "Should resize if outside up of NW corner",
			point:  up1(cornerNW),
			newMin: up(left(expectedMin)),
			newMax: expectedMax,
		},
		{
			name:   "Should resize if outside up of NE corner",
			point:  up1(cornerNE),
			newMin: up(expectedMin),
			newMax: right(expectedMax),
		},
		{
			name:   "Should resize if outside right of NE corner",
			point:  right1(cornerNE),
			newMin: up(expectedMin),
			newMax: right(expectedMax),
		},
		{
			name:   "Should resize if outside right of SE corner",
			point:  right1(cornerSE),
			newMin: expectedMin,
			newMax: right(down(expectedMax)),
		},
		{
			name:   "Should resize if outside down of SE corner",
			point:  down1(cornerSE),
			newMin: expectedMin,
			newMax: down(right(expectedMax)),
		},
		{
			name:   "Should resize if outside down of SW corner",
			point:  down1(cornerSW),
			newMin: left(expectedMin),
			newMax: down(expectedMax),
		},
		{
			name:   "Should resize if outside left of SW corner",
			point:  left1(cornerSW),
			newMin: left(expectedMin),
			newMax: down(expectedMax),
		},
	}

	for _, tt := range testsResize {
		t.Run(tt.name, func(t *testing.T) {
			subj := upArray{}
			subj.Initialize(initialPoint)
			subj.ResizeIfNeeded(tt.point)

			if !subj.Min.equals(tt.newMin) {
				t.Errorf("Expected min %s, got %s Size: %dx%d", tt.newMin.String(), subj.Min.String(), width(&subj), height(&subj))
			}

			if !subj.Max.equals(tt.newMax) {
				t.Errorf("Expected max %s, got %s Size: %dx%d", tt.newMax.String(), subj.Max.String(), width(&subj), height(&subj))
			}
		})
	}
}

func TestCopy(t *testing.T) {
	largerBy := upInt(10)

	var testsCopy = []struct {
		name   string
		newMin gridCoordsUp
		newMax gridCoordsUp
	}{
		{
			name:   "Same rect",
			newMin: gridCoordsUp{0, 0},
			newMax: gridCoordsUp{5, 5},
		},
		{
			name:   "Lager at N side",
			newMin: gridCoordsUp{0, 0 - largerBy},
			newMax: gridCoordsUp{5, 5},
		},
		{
			name:   "Lager at NE side",
			newMin: gridCoordsUp{0, 0 - largerBy},
			newMax: gridCoordsUp{5 + largerBy, 5},
		},
		{
			name:   "Lager at E side",
			newMin: gridCoordsUp{0, 0},
			newMax: gridCoordsUp{5 + largerBy, 5},
		},
		{
			name:   "Lager at SE side",
			newMin: gridCoordsUp{0, 0},
			newMax: gridCoordsUp{5 + largerBy, 5 + largerBy},
		},
		{
			name:   "Lager at S side",
			newMin: gridCoordsUp{0, 0},
			newMax: gridCoordsUp{5, 5 + largerBy},
		},
		{
			name:   "Lager at SW side",
			newMin: gridCoordsUp{0 - largerBy, 0},
			newMax: gridCoordsUp{5, 5 + largerBy},
		},
		{
			name:   "Lager at W side",
			newMin: gridCoordsUp{0 - largerBy, 0},
			newMax: gridCoordsUp{5, 5},
		},
		{
			name:   "Lager at NW side",
			newMin: gridCoordsUp{0 - largerBy, 0 - largerBy},
			newMax: gridCoordsUp{5, 5},
		},
	}

	for _, tt := range testsCopy {
		t.Run(tt.name, func(t *testing.T) {
			subj := newUpArray(gridCoordsUp{0, 0}, gridCoordsUp{5, 5})

			p := GridCoords{1, 1}

			val, coordsDown := subj.Get(p)
			val[coordsDown] = 42
			val0, coordsDown0 := subj.Get(p)
			oldValue := val0[coordsDown0]

			if oldValue != 42 {
				t.Errorf("Expected same value before copy got %d", oldValue)
			}

			newMin := tt.newMin
			newMax := tt.newMax

			subj.Grow(newMin, newMax)

			val1, coordsDown1 := subj.Get(p)
			newValue := val1[coordsDown1]
			if newValue != 42 {
				t.Errorf("Expected same value after copy got %d", newValue)
			}
		})
	}
}
