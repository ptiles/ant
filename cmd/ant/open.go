package main

import (
	"os/exec"
	"runtime"
)

func open(fileName string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", fileName).Run()
	case "windows":
		exec.Command("start", fileName).Run()
	default:
		exec.Command("xdg-open", fileName).Run()
	}
}
