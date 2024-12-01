// kamera demo

package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
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
	camSpeed, zoomSpeedFactor, rotSpeed = 7.0, 1.02, 0.02
	targetX, targetY                    = w / 2, h / 2
	cam                                 = kamera.NewCamera(targetX, targetY, w, h)
	dio                                 = &ebiten.DrawImageOptions{}
	spriteSheet                         *ebiten.Image
)

type Game struct{}

func (g *Game) Update() error {
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
		switch cam.Smoothing {
		case kamera.None:
			cam.Smoothing = kamera.Lerp
		case kamera.Lerp:
			cam.Smoothing = kamera.SmoothDamp
		case kamera.SmoothDamp:
			cam.Smoothing = kamera.None
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		targetX -= camSpeed / cam.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		targetX += camSpeed / cam.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		targetY -= camSpeed / cam.ZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		targetY += camSpeed / cam.ZoomFactor
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

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			dio.GeoM.Reset()
			dio.GeoM.Translate(float64(x*300), float64(y*300))
			cam.Draw(spriteSheet, dio, screen)
		}
	}

	// Draw camera
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
	cam.Smoothing = kamera.SmoothDamp
	img, _, err := image.Decode(bytes.NewReader(images.Smoke_png))
	if err != nil {
		log.Fatal(err)
	}
	spriteSheet = ebiten.NewImageFromImage(img)
	ebiten.SetWindowSize(int(w), int(h))
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}