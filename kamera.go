// kamera is a camera package for ebitengine.
package kamera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera object
type Camera struct {
	// (X, Y) is the top left corner of the camera viewport.
	// Use the `Camera.LookAt()` function to align the center of the camera to the target.
	// (W, H) is the width and height of the camera viewport.
	X, Y, W, H, Rotation, ZoomFactor float64

	DrawOptions *ebiten.DrawImageOptions
}

// NewCamera returns new Camera
func NewCamera(w, h float64) *Camera {
	return &Camera{
		W: w, H: h, X: 0, Y: 0, Rotation: 0, ZoomFactor: 0, DrawOptions: &ebiten.DrawImageOptions{},
	}
}

// LookAt aligns the midpoint of the camera viewport to the target.
func (cam *Camera) LookAt(targetX, targetY float64) {
	cam.X, cam.Y = targetX-cam.W*0.5, targetY-cam.H*0.5

}

// Reset resets all camera values to zero
func (cam *Camera) Reset() {
	cam.X, cam.Y = 0.0, 0.0
	cam.Rotation, cam.ZoomFactor = 0.0, 0
}

// String returns camera values as string
func (cam *Camera) String() string {
	return fmt.Sprintf(
		"CamX: %.1f\nCamY: %.1f\nCam Rotation: %.1f\nZoom factor: %.2f",
		cam.X, cam.Y, cam.Rotation, cam.ZoomFactor,
	)
}

// ScreenToWorld converts screen-space coordinates to world-space
func (cam *Camera) ScreenToWorld(screenX, screenY int) (worldX, worldY float64) {
	g := ebiten.GeoM{}
	cam.applyCameraTransform(&g)
	if g.IsInvertible() {
		g.Invert()
		worldX, worldY := g.Apply(float64(screenX), float64(screenY))
		return worldX, worldY
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

// applyCameraTransform applies geometric transformation to given geoM
func (cam *Camera) applyCameraTransform(geoM *ebiten.GeoM) {
	geoM.Translate(-cam.X, -cam.Y)                                                               // camera movement
	geoM.Translate(-cam.W*0.5, -cam.H*0.5)                                                       // rotate and scale from center.
	geoM.Rotate(float64(cam.Rotation) * 2 * math.Pi / 360)                                       // rotate
	geoM.Scale(math.Pow(1.01, float64(cam.ZoomFactor)), math.Pow(1.01, float64(cam.ZoomFactor))) // apply zoom factor
	geoM.Translate(cam.W*0.5, cam.H*0.5)                                                         // restore center translation
}

// Render applies the Camera's geometric transformation then draws the object on the screen with drawing options.
func (cam *Camera) Render(worldObject *ebiten.Image, worldObjectOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	cam.DrawOptions = worldObjectOps
	cam.applyCameraTransform(&cam.DrawOptions.GeoM)
	screen.DrawImage(worldObject, cam.DrawOptions)
	cam.DrawOptions.GeoM.Reset()
}
