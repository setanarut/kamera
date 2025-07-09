// kamera demo

package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var (
	Controls = `
Key           Action                     
------------  ------------------------
WASD          Move camera                
T             Add 1.0 Trauma             
Tab           Change camera smoothing type
Space         Look at random position    
E/Q           Zoom in/out                
ArrowUp/Down  Zoom 2x                    
Backspace     Reset camera               
R             Rotate                     
`
	w, h                                = 1024., 768.
	camSpeed, zoomSpeedFactor, rotSpeed = 2.0, 1.02, 0.02
	targetX, targetY                    = 0., 0.
	cam                                 = kamera.NewCamera(targetX, targetY, w, h)
	dio                                 = &ebiten.DrawImageOptions{}
	spriteSheet                         *ebiten.Image
)

type Game struct{}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		cam.AddTrauma(1.0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		cam.ZoomFactor *= 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		cam.ZoomFactor /= 2
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		targetX, targetY = rand.Float64()*200, rand.Float64()*200
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		switch cam.SmoothType {
		case kamera.None:
			cam.SetCenter(targetX, targetY)
			cam.SmoothType = kamera.Lerp
		case kamera.Lerp:
			cam.SetCenter(targetX, targetY)
			cam.SmoothType = kamera.SmoothDamp
		case kamera.SmoothDamp:
			cam.SetCenter(targetX, targetY)
			cam.SmoothType = kamera.None
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		cam.ZoomFactor /= zoomSpeedFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		cam.ZoomFactor *= zoomSpeedFactor
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		cam.Angle += rotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		cam.Angle -= rotSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		targetX, targetY = 0, 0
		cam.SetCenter(0, 0)
		cam.Reset()
	}

	aX, aY := Normalize(Axis())
	aX *= camSpeed
	aY *= camSpeed

	targetX += aX
	targetY += aY

	cam.LookAt(targetX, targetY)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	dio.GeoM.Reset() // GeoM must be reset
	cam.Draw(spriteSheet, dio, screen)

	// Draw camera crosshair
	cx, cy := float32(w/2), float32(h/2)
	vector.StrokeLine(screen, cx-100, cy, cx+100, cy, 1, color.White, true)
	vector.StrokeLine(screen, cx, cy-100, cx, cy+100, 1, color.White, true)
	// HUD
	ebitenutil.DebugPrintAt(screen, Controls, 10, 10)
	ebitenutil.DebugPrintAt(screen, cam.String(), 10, 250)
}

func (g *Game) Layout(width, height int) (int, int) {
	return int(w), int(h)
}

func main() {
	cam.SmoothType = kamera.SmoothDamp
	cam.ShakeEnabled = true
	img, _, err := image.Decode(bytes.NewReader(images.Gophers_jpg))
	if err != nil {
		log.Fatal(err)
	}
	spriteSheet = ebiten.NewImageFromImage(img)
	ebiten.SetWindowSize(int(w), int(h))
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func Normalize(x, y float64) (float64, float64) {
	magnitude := math.Sqrt(x*x + y*y)
	if magnitude == 0 {
		return 0, 0
	}
	return x / magnitude, y / magnitude
}

func Axis() (axisX, axisY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axisY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axisY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axisX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axisX += 1
	}
	return axisX, axisY
}
