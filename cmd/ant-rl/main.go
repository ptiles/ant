package main

import (
	"flag"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ptiles/ant/pgrid"
	"github.com/ptiles/ant/utils"
	"image"
	"math"
	"os"
	"path/filepath"
)

var (
	programName = filepath.Base(os.Args[0])
	usageText   = `Run Langton's ant on Penrose tiling (pentagrid)

Usage of %s:
	%s [flags] [name RLRRLRR...]

Name should consist of letters R, L, r, l.

Flags:
`
	usageTextShort = "\nFor usage run: %s -h\n"
)

func toggleFullScreenWindow(windowWidth, windowHeight int) rl.Vector2 {
	if rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
		rl.SetWindowSize(windowWidth, windowHeight)
		return rl.Vector2{X: float32(windowWidth), Y: float32(windowHeight)}
	} else {
		monitor := rl.GetCurrentMonitor()
		monitorWidth := rl.GetMonitorWidth(monitor)
		monitorHeight := rl.GetMonitorHeight(monitor)
		rl.SetWindowSize(monitorWidth, monitorHeight)
		rl.ToggleFullscreen()
		return rl.Vector2{X: float32(monitorWidth), Y: float32(monitorHeight)}
	}
}

func main() {
	commonFlags := &utils.CommonFlags{}
	commonFlags.CommonFlagsSetup(pgrid.GridLinesTotal)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText, programName, programName)
		flag.PrintDefaults()
	}
	flag.Parse()
	commonFlags.ParseArgs()

	rules, err := utils.GetRules(commonFlags.AntName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid name.  Should consist of at least two letters R L r l.")
		fmt.Fprintf(os.Stderr, usageTextShort, programName)
		os.Exit(1)
	}

	utils.StartCPUProfile(commonFlags.Cpuprofile)
	defer utils.StopCPUProfile()

	field := pgrid.New(commonFlags.Radius, rules, commonFlags.InitialPoint)
	palette := utils.GetPaletteRainbow(len(rules))

	const (
		screenWidth  = 1680
		screenHeight = 1050
	)
	const zoomIncrement float32 = 0.025

	title := fmt.Sprintf("ant-rl %s %s", commonFlags.InitialPoint, commonFlags.AntName)
	rl.InitWindow(screenWidth, screenHeight, title)

	var camera rl.Camera2D
	camera.Zoom = 1.0

	topLeftScreen := rl.Vector2{X: 0, Y: 0}
	bottomRightScreen := rl.Vector2{X: screenWidth, Y: screenHeight}

	rl.SetTargetFPS(60)

	modifiedImagesCh := make(chan *image.RGBA, 256)
	fullViewRectCh := make(chan image.Rectangle)
	positionedImagesCh := make(chan *PositionedImage)
	commandCh := make(chan pgrid.CommandType, 1)

	go field.ControlledInfiniteStepper(modifiedImagesCh, commandCh, palette)

	go imageTilesServer(modifiedImagesCh, fullViewRectCh, positionedImagesCh)

	var textures = make(map[Pair]rl.Texture2D)

	for !rl.WindowShouldClose() {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			delta := rl.GetMouseDelta()
			delta = rl.Vector2Scale(delta, -1.0/camera.Zoom)

			camera.Target = rl.Vector2Add(camera.Target, delta)
		}

		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
			camera.Offset = rl.GetMousePosition()
			camera.Target = mouseWorldPos

			camera.Zoom *= float32(math.Pow(2, float64(wheel*zoomIncrement)))
			if camera.Zoom < 0.01 {
				camera.Zoom = 0.01
			}
		}

		if rl.IsKeyPressed(rl.KeyF) {
			bottomRightScreen = toggleFullScreenWindow(screenWidth, screenHeight)
		}
		if rl.IsKeyPressed(rl.KeyR) {
			commandCh <- pgrid.Reset

			for _, texture := range textures {
				rl.UnloadTexture(texture)
			}
			textures = make(map[Pair]rl.Texture2D)

			fullViewRectCh <- image.Rectangle{}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(camera)

		topLeft := rl.GetScreenToWorld2D(topLeftScreen, camera)
		bottomRight := rl.GetScreenToWorld2D(bottomRightScreen, camera)

		minPoint := image.Point{X: int(topLeft.X), Y: int(topLeft.Y)}
		maxPoint := image.Point{X: int(bottomRight.X), Y: int(bottomRight.Y)}

		fullViewRect := image.Rectangle{Min: minPoint, Max: maxPoint}
		fullViewRectCh <- fullViewRect

		for positionedImage := range positionedImagesCh {
			if positionedImage == nil {
				break
			}

			if texture, ok := textures[positionedImage.position]; ok {
				rl.UnloadTexture(texture)
			}

			texture := rl.LoadTextureFromImage(positionedImage.image)
			textures[positionedImage.position] = texture
		}

		for _, viewRect := range splitViewRect(fullViewRect) {
			texture := textures[viewRect]
			rl.DrawTexture(texture, int32(viewRect.X*256), int32(viewRect.Y*256), rl.White)
		}

		rl.EndMode2D()

		rl.DrawText("Mouse left button drag to move, mouse wheel to zoom", 10, 10, 20, rl.White)
		rl.DrawText(fmt.Sprintf("%.4f", camera.Zoom), 10, 40, 20, rl.White)
		rl.DrawFPS(int32(bottomRightScreen.X)-100, 10)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
