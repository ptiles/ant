package main

import "github.com/ptiles/ant/store"

func walk(coords store.PackedCoordinates, steps []bool, maxValue uint8) bool {
	value := store.Get(coords)
	store.Set(coords, (value+1)%maxValue)
	return steps[value]
}
