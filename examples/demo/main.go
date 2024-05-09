// kamera demo

package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
)

var Controls = `
| Key       | Action                     |
| -----     | -------------------------- |
| WASD      | Move camera                |
| T         | Add 0.6 Trauma             |
| E         | Zoom in                    |
| Q         | Zoom out                   |
| Backspace | Reset Rotation/Zoom        |
| R         | Rotate                     |
| X         | Look at the random object. |
| L         | Toggle Lerp                |
`

type vec struct {
	X, Y float64
}

type Game struct {
	ScreenSize   *image.Point
	GameObjects  []*ebiten.Image
	MainCamera   *kamera.Camera
	CamSpeed     float64
	ZoomSpeed    float64
	FontSize     float64
	RandomPoints []vec
	DIO          *ebiten.DrawImageOptions
	halfSize     float64
}

// MakeObjects creates random images with random colors
func MakeObjects(n, size int) []*ebiten.Image {
	imgs := make([]*ebiten.Image, n)
	for i := range imgs {
		imgs[i] = ebiten.NewImage(size, size)
		imgs[i].Fill(color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
	}
	return imgs
}

func RandomPoints(minX, maxX, minY, maxY float64, n int) []vec {
	points := make([]vec, n)
	for i := range points {
		points[i] = vec{X: minX + rand.Float64()*(maxX-minX), Y: minY + rand.Float64()*(maxY-minY)}
	}
	return points
}

var delta vec
var tick = 0.0
var TargetX, TargetY float64

func (g *Game) Update() error {

	// Use LookAt() only once in update
	g.MainCamera.LookAt(TargetX, TargetY)

	// reset delta
	delta.X = 0
	delta.Y = 0

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		index := rand.Intn(len(g.RandomPoints))
		TargetX = g.RandomPoints[index].X
		TargetY = g.RandomPoints[index].Y
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.MainCamera.Lerp = !g.MainCamera.Lerp
	}

	// trauma
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.MainCamera.AddTrauma(0.6)
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		delta.X = -g.CamSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		delta.X = g.CamSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		delta.Y = -g.CamSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		delta.Y = g.CamSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		delta.Y = g.CamSpeed
	}
	// Check for diagonal movement
	if delta.X != 0 && delta.Y != 0 {
		factor := g.CamSpeed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}

	// Move Camera (WASD)
	TargetX += delta.X
	TargetY += delta.Y

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		if g.MainCamera.ZoomFactor > -4800 {
			g.MainCamera.ZoomFactor -= g.ZoomSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		if g.MainCamera.ZoomFactor < 4800 {
			g.MainCamera.ZoomFactor += g.ZoomSpeed
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.MainCamera.Rotation += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.MainCamera.Rotation -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.MainCamera.Reset()
	}
	// tick for rotation
	tick += 0.02
	if tick > math.Pi*2 {
		tick = 0.0
	}
	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, randomPoint := range g.RandomPoints {
		g.DIO.GeoM.Reset()
		// rotate objects about center and move to random point
		g.DIO.GeoM.Translate(-g.halfSize, -g.halfSize)
		g.DIO.GeoM.Rotate(tick)
		g.DIO.GeoM.Translate(randomPoint.X, randomPoint.Y)
		// render objects
		g.MainCamera.Draw(g.GameObjects[i], g.DIO, screen)
	}
	ebitenutil.DebugPrintAt(screen, Controls, 10, 10)
	ebitenutil.DebugPrintAt(screen, g.MainCamera.String(), 10, 200)
	// draw circle at center of camera

	vector.DrawFilledCircle(screen, float32(g.ScreenSize.X/2), float32(g.ScreenSize.Y/2), 4, color.White, false)
}

func (oyn *Game) Layout(w, h int) (int, int) {
	return oyn.ScreenSize.X, oyn.ScreenSize.Y
}

func main() {
	// minf, maxf := math.SmallestNonzeroFloat64, math.MaxFloat64
	minf, maxf := -5000.0, 5000.0
	enemyCount := 2000
	enemySize := 64
	w, h := 854, 480

	game := &Game{
		ScreenSize:   &image.Point{int(w), int(h)},
		ZoomSpeed:    3,
		GameObjects:  MakeObjects(enemyCount, enemySize),
		MainCamera:   kamera.NewCamera(0, 0, float64(w), float64(h)),
		CamSpeed:     5,
		RandomPoints: RandomPoints(minf, maxf, minf, maxf, enemyCount),
		DIO:          &ebiten.DrawImageOptions{},
		halfSize:     float64(enemySize) / 2,
	}

	game.MainCamera.Lerp = true

	ebiten.SetWindowSize(game.ScreenSize.X, game.ScreenSize.Y)
	ebiten.RunGame(game)

}
