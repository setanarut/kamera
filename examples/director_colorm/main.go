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
	"github.com/hajimehoshi/ebiten/v2/colorm"
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
	camSpeed, zoomSpeedFactor, rotSpeed = 7.0, 1.02, 0.02
	targetX, targetY                    = w / 2, h / 2
	cam                                 = kamera.NewCamera(targetX, targetY, w, h)
	colormDIO                           = &colorm.DrawImageOptions{}
	spriteSheet                         *ebiten.Image
)

type Game struct{}

func (g *Game) Update() error {

	aX, aY := Normalize(Axis())

	targetX += aX * camSpeed
	targetY += aY * camSpeed

	cam.LookAt(targetX, targetY)

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
		targetX, targetY = rand.Float64()*w, rand.Float64()*h
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		switch cam.SmoothType {
		case kamera.None:
			cam.SmoothType = kamera.Lerp
		case kamera.Lerp:
			cam.SmoothType = kamera.SmoothDamp
		case kamera.SmoothDamp:
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
		cam.SetAngle(cam.Angle() + rotSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		cam.SetAngle(cam.Angle() - rotSpeed)
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		targetX, targetY = w/2, h/2
		cam.Reset()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	var cm colorm.ColorM

	cm.ChangeHSV(2, 1, 0.5)

	// Draw backgorund tiles
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			colormDIO.GeoM.Reset()
			colormDIO.GeoM.Translate(float64(x*300), float64(y*300))
			cam.DrawWithColorM(spriteSheet, cm, colormDIO, screen)
		}
	}

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
