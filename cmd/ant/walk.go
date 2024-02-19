package main

import "github.com/ptiles/ant/store"

func walk(coords store.PackedCoordinates, steps []bool, maxValue uint8) bool {
	value := (store.Get(coords) + 1) % maxValue
	store.Set(coords, value)
	return steps[value]
}
