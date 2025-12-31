// kamera demo

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/kamera/v2/examples"
	"github.com/setanarut/v"
)

var Controls = `
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
C             Toggle DrawWithColorM() %v                
`

var (
	w, h                                = 1024., 768.
	camSpeed, zoomSpeedFactor, rotSpeed = 6.0, 1.02, 0.02
	target                              = v.Vec{}
	mainCamera                          = kamera.NewCamera(target.X, target.Y, w, h)
	dio                                 = &ebiten.DrawImageOptions{}
	cdio                                = &colorm.DrawImageOptions{}
	clrm                                = colorm.ColorM{}
	spriteSheet                         *ebiten.Image

	colormEnabled bool
)

type Game struct{}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		mainCamera.AddTrauma(1.0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		colormEnabled = !colormEnabled
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		mainCamera.ZoomFactor *= 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		mainCamera.ZoomFactor /= 2
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		target.X, target.Y = rand.Float64()*200, rand.Float64()*200
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		switch mainCamera.SmoothType {
		case kamera.None:
			mainCamera.SetCenter(target.X, target.Y)
			mainCamera.SmoothType = kamera.Lerp
		case kamera.Lerp:
			mainCamera.SetCenter(target.X, target.Y)
			mainCamera.SmoothType = kamera.SmoothDamp
		case kamera.SmoothDamp:
			mainCamera.SetCenter(target.X, target.Y)
			mainCamera.SmoothType = kamera.None
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		mainCamera.ZoomFactor /= zoomSpeedFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		mainCamera.ZoomFactor *= zoomSpeedFactor
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		mainCamera.Angle += rotSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		mainCamera.Angle -= rotSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		target = v.Vec{}
		mainCamera.SetCenter(0, 0)
		mainCamera.Reset()
	}

	a := examples.Axis().Unit().Scale(camSpeed)
	target = target.Add(a)
	mainCamera.LookAt(target.X, target.Y)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if colormEnabled {
		cdio.GeoM.Reset() // GeoM must be reset
		clrm.Reset()      // ColorM must be reset
		clrm.ChangeHSV(2, 1, 0.5)
		mainCamera.DrawWithColorM(spriteSheet, clrm, cdio, screen)
	} else {
		dio.GeoM.Reset() // GeoM must be reset
		mainCamera.Draw(spriteSheet, dio, screen)
	}

	// Draw camera crosshair
	cx, cy := float32(w/2), float32(h/2)
	vector.StrokeLine(screen, cx-100, cy, cx+100, cy, 1, color.White, true)
	vector.StrokeLine(screen, cx, cy-100, cx, cy+100, 1, color.White, true)
	// HUD
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf(Controls, colormEnabled), 10, 10)
	ebitenutil.DebugPrintAt(screen, mainCamera.String(), 10, 250)
}

func (g *Game) Layout(width, height int) (int, int) {
	return int(w), int(h)
}

func main() {
	mainCamera.SmoothType = kamera.SmoothDamp
	mainCamera.ShakeEnabled = true
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
