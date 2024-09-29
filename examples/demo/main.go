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
| Key          | Action                     |
| ------------ | -------------------------- |
| WASD         | Move camera                |
| T            | Add 1.0 Trauma             |
| E/Q          | Zoom in/out                |
| Tab          | Look at random position    |
| ArrowUp/Down | Zoom 2x                    |
| Backspace    | Reset camera               |
| R            | Rotate                     |
| L            | Toggle Lerp                |
`
	w, h                                float64                  = 1024, 768
	camSpeed, zoomSpeedFactor, rotSpeed float64                  = 1.01, 1.02, 0.02
	targetX, targetY                    float64                  = w / 2, h / 2
	cam                                 *kamera.Camera           = kamera.NewCamera(targetX, targetY, w, h)
	dio                                 *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	gophersImage                        *ebiten.Image
)

func init() {
	cam.Lerp = true
}

type Game struct{}

func (g *Game) Update() error {
	// Use LookAt() only once in update
	cam.LookAt(targetX, targetY)

	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		cam.Lerp = !cam.Lerp
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		cam.AddTrauma(1.0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		cam.ZoomFactor *= 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		cam.ZoomFactor /= 2
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		targetX, targetY = rand.Float64()*w, rand.Float64()*h
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		targetX /= camSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		targetX *= camSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		targetY /= camSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		targetY *= camSpeed
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

	// Draw camera
	cam.Draw(gophersImage, dio, screen)

	// Draw camera crosshair
	cx, cy := float32(w/2), float32(h/2)
	vector.StrokeLine(screen, cx-10, cy, cx+10, cy, 1, color.White, true)
	vector.StrokeLine(screen, cx, cy-10, cx, cy+10, 1, color.White, true)
	// HUD
	ebitenutil.DebugPrintAt(screen, Controls, 10, 10)
	ebitenutil.DebugPrintAt(screen, cam.String(), 10, 200)
}

func (g *Game) Layout(width, height int) (int, int) {
	return int(w), int(h)
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(images.Spritesheet_png))
	if err != nil {
		log.Fatal(err)
	}
	gophersImage = ebiten.NewImageFromImage(img)
	ebiten.SetWindowSize(int(w), int(h))
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
