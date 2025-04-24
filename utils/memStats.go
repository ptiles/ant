package utils

import (
	"runtime"
)

func MemStatsMB() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return (m.Sys - m.HeapReleased) / 1024 / 1024
}
