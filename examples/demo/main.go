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
	MainCamera   *kamera.Camera
	RandomColors []*color.RGBA
	RandomPoints []vec
	Obj          *ebiten.Image
	CamSpeed     float64
	ZoomSpeed    float64
	FontSize     float64
	DIO          *ebiten.DrawImageOptions
	halfSize     float64
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
		g.MainCamera.SetRotation(g.MainCamera.Rotation() + 1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.MainCamera.SetRotation(g.MainCamera.Rotation() - 1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.MainCamera.Reset()
	}
	// tick for object rotation
	tick += 0.02
	if tick > math.Pi*2 {
		tick = 0.0
	}
	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	for i, randomPoint := range g.RandomPoints {

		g.DIO.GeoM.Reset()
		g.DIO.ColorScale.Reset()

		g.DIO.GeoM.Translate(-g.halfSize, -g.halfSize)
		g.DIO.GeoM.Rotate(tick)
		g.DIO.GeoM.Translate(randomPoint.X, randomPoint.Y)

		g.DIO.ColorScale.ScaleWithColor(g.RandomColors[i])

		// Draw camera
		g.MainCamera.Draw(g.Obj, g.DIO, screen)
	}

	ebitenutil.DebugPrintAt(screen, Controls, 10, 10)
	ebitenutil.DebugPrintAt(screen, g.MainCamera.String(), 10, 200)

	// draw circle at center of camera
	vector.DrawFilledCircle(screen, float32(g.ScreenSize.X/2), float32(g.ScreenSize.Y/2), 4, color.White, false)
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.ScreenSize.X, g.ScreenSize.Y
}

func main() {
	// minf, maxf := math.SmallestNonzeroFloat64, math.MaxFloat64
	bound := 2000.0
	objCount := 500
	objSize := 64
	w, h := 854, 480

	g := &Game{
		ScreenSize:   &image.Point{int(w), int(h)},
		ZoomSpeed:    3,
		MainCamera:   kamera.NewCamera(0, 0, float64(w), float64(h)),
		CamSpeed:     5,
		RandomPoints: RandomPoints(-bound, bound, -bound, bound, objCount),
		RandomColors: RandomColors(objCount),
		Obj:          ebiten.NewImage(objSize, objSize),
		DIO:          &ebiten.DrawImageOptions{},
		halfSize:     float64(objSize) / 2,
	}

	g.MainCamera.Lerp = true
	g.Obj.Fill(color.White)
	TargetX, TargetY = g.RandomPoints[0].X, g.RandomPoints[0].Y
	ebiten.SetWindowSize(g.ScreenSize.X, g.ScreenSize.Y)
	ebiten.RunGame(g)

}

func RandomColors(n int) []*color.RGBA {
	colors := make([]*color.RGBA, 0)
	for range n {
		colors = append(colors, &color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
	}
	return colors
}

func RandomPoints(minX, maxX, minY, maxY float64, n int) []vec {
	points := make([]vec, n)
	for i := range points {
		points[i] = vec{X: minX + rand.Float64()*(maxX-minX), Y: minY + rand.Float64()*(maxY-minY)}
	}
	return points
}
