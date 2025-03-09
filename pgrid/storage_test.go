package pgrid

import (
	"testing"
)

func height(ua *upArray) offsetInt {
	return ua.Max.Offset0 - ua.Min.Offset0
}

func width(ua *upArray) offsetInt {
	return ua.Max.Offset1 - ua.Min.Offset1
}

func TestInitialize(t *testing.T) {
	initial := GridCoords{Offset0: -128, Offset1: 127}
	expectedMin := GridCoords{Offset0: -256, Offset1: 0}
	expectedMax := GridCoords{Offset0: 0, Offset1: 256}

	subj := upArray{}
	subj.Initialize(initial)

	expectedSize := offsetInt(256)
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

func up1(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0, Offset1: gc.Offset1 - 1}
}
func up(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0, Offset1: gc.Offset1 - 256}
}
func right1(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0 + 1, Offset1: gc.Offset1}
}
func right(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0 + 256, Offset1: gc.Offset1}
}
func down1(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0, Offset1: gc.Offset1 + 1}
}
func down(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0, Offset1: gc.Offset1 + 256}
}
func left1(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0 - 1, Offset1: gc.Offset1}
}
func left(gc GridCoords) GridCoords {
	return GridCoords{Offset0: gc.Offset0 - 256, Offset1: gc.Offset1}
}

func TestResize(t *testing.T) {
	t.SkipNow()

	initialPoint := GridCoords{Offset0: -128, Offset1: 127}

	expectedMin := GridCoords{Offset0: -256, Offset1: 0}
	expectedMax := GridCoords{Offset0: 0, Offset1: 256}
	expectedMinPoint := GridCoords{Offset0: -256, Offset1: 0}
	expectedMaxPoint := GridCoords{Offset0: -1, Offset1: 255}

	centerPoint := GridCoords{
		Offset0: (expectedMin.Offset0 + expectedMax.Offset0) / 2,
		Offset1: (expectedMin.Offset1 + expectedMax.Offset1) / 2,
	}
	cornerNW := expectedMinPoint
	cornerNE := GridCoords{Offset0: expectedMaxPoint.Offset0, Offset1: expectedMinPoint.Offset1}
	cornerSE := expectedMaxPoint
	cornerSW := GridCoords{Offset0: expectedMinPoint.Offset0, Offset1: expectedMaxPoint.Offset1}

	var testsKeep = []struct {
		name  string
		point GridCoords
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
		point  GridCoords
		newMin GridCoords
		newMax GridCoords
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
	coordsDown := gridCoordsDown{}
	largerBy := offsetInt(10)

	var testsCopy = []struct {
		name   string
		newMin GridCoords
		newMax GridCoords
	}{
		{
			name:   "Same rect",
			newMin: GridCoords{0, 0},
			newMax: GridCoords{5, 5},
		},
		{
			name:   "Lager at N side",
			newMin: GridCoords{0, 0 - largerBy},
			newMax: GridCoords{5, 5},
		},
		{
			name:   "Lager at NE side",
			newMin: GridCoords{0, 0 - largerBy},
			newMax: GridCoords{5 + largerBy, 5},
		},
		{
			name:   "Lager at E side",
			newMin: GridCoords{0, 0},
			newMax: GridCoords{5 + largerBy, 5},
		},
		{
			name:   "Lager at SE side",
			newMin: GridCoords{0, 0},
			newMax: GridCoords{5 + largerBy, 5 + largerBy},
		},
		{
			name:   "Lager at S side",
			newMin: GridCoords{0, 0},
			newMax: GridCoords{5, 5 + largerBy},
		},
		{
			name:   "Lager at SW side",
			newMin: GridCoords{0 - largerBy, 0},
			newMax: GridCoords{5, 5 + largerBy},
		},
		{
			name:   "Lager at W side",
			newMin: GridCoords{0 - largerBy, 0},
			newMax: GridCoords{5, 5},
		},
		{
			name:   "Lager at NW side",
			newMin: GridCoords{0 - largerBy, 0 - largerBy},
			newMax: GridCoords{5, 5},
		},
	}

	for _, tt := range testsCopy {
		t.Run(tt.name, func(t *testing.T) {
			subj := newUpArray(GridCoords{0, 0}, GridCoords{5, 5})

			p := GridCoords{1, 1}
			subj.Get(p)[coordsDown] = 42
			oldValue := subj.Get(p)[coordsDown]

			if oldValue != 42 {
				t.Errorf("Expected same value before copy got %d", oldValue)
			}

			newMin := tt.newMin
			newMax := tt.newMax

			subj.Maps, subj.Stride = subj.Copy(newMin, newMax)
			subj.Min = newMin
			subj.Max = newMax

			newValue := subj.Get(p)[coordsDown]
			if newValue != 42 {
				t.Errorf("Expected same value after copy got %d", newValue)
			}
		})
	}
}
