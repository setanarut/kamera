// topdown kamera demo

/*
| Key   | Action                     |
| ----- | -------------------------- |
| WASD  | Move camera                |
| E     | Zoom in                    |
| Q     | Zoom out                   |
| Space | Reset camera               |
| Key 1 | look at the random object. |
| Key 2 | Reset Zoom                 |
*/
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
	"github.com/setanarut/kamera"
)

type Vec2 struct {
	X, Y float64
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

func RandomPoints(minX, maxX, minY, maxY float64, n int) []Vec2 {
	points := make([]Vec2, n)
	for i := range points {
		points[i] = Vec2{X: minX + rand.Float64()*(maxX-minX), Y: minY + rand.Float64()*(maxY-minY)}
	}
	return points
}

var delta Vec2
var tick = 0.0

func (g *Game) Update() error {
	// reset delta
	delta.X = 0
	delta.Y = 0
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		index := rand.Intn(len(g.RandomPoints))
		g.MainCamera.LookAt(g.RandomPoints[index].X, g.RandomPoints[index].Y)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.MainCamera.ZoomFactor = 0
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
	// Check for diagonal movement
	if delta.X != 0 && delta.Y != 0 {
		factor := g.CamSpeed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}
	g.MainCamera.X += delta.X
	g.MainCamera.Y += delta.Y

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

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
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
		g.MainCamera.Render(g.GameObjects[i], g.DIO, screen)
	}
	ebitenutil.DebugPrint(screen, g.MainCamera.String())
	// draw circle at center of camera
	vector.DrawFilledCircle(screen, float32(g.MainCamera.W)/2,
		float32(g.MainCamera.H)/2, 4, color.White, false)
}

func (oyn *Game) Layout(w, h int) (int, int) {
	return oyn.InternalResolution.X, oyn.InternalResolution.Y
}

type Game struct {
	InternalResolution *image.Point
	GameObjects        []*ebiten.Image
	MainCamera         *kamera.Camera
	CamSpeed           float64
	ZoomSpeed          float64
	FontSize           float64
	RandomPoints       []Vec2
	DIO                *ebiten.DrawImageOptions
	halfSize           float64
}

func main() {
	// minf, maxf := math.SmallestNonzeroFloat64, math.MaxFloat64
	minf, maxf := -5000.0, 5000.0
	enemyCount := 2000
	enemySize := 64
	w, h := 854, 480
	game := &Game{
		InternalResolution: &image.Point{int(w), int(h)},
		ZoomSpeed:          3,
		GameObjects:        MakeObjects(enemyCount, enemySize),
		MainCamera:         kamera.NewCamera(float64(w), float64(h)),
		CamSpeed:           5,
		RandomPoints:       RandomPoints(minf, maxf, minf, maxf, enemyCount),
		DIO:                &ebiten.DrawImageOptions{},
		halfSize:           float64(enemySize) / 2,
	}

	ebiten.SetWindowSize(game.InternalResolution.X, game.InternalResolution.Y)
	ebiten.RunGame(game)

}
