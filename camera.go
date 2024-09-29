package kamera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/fastnoise"
)

// Camera object
// Use the `Camera.LookAt()` function to align the center of the camera to the target.
type Camera struct {

	// ZoomFactor is the camera zoom (scaling) factor. Default is 1.
	ZoomFactor float64

	// Interpolate camera movement
	Lerp      bool
	LerpSpeed float64

	// Camera shake options
	ShakeOptions *CameraShakeOptions

	// private
	drawOptions                                                        *ebiten.DrawImageOptions
	angle, actualAngle, tickSpeed, tick, trauma, w, h, zoomFactorShake float64
	tempTarget, centerOffset, topLeft, traumaOffset                    vec2
}

// NewCamera returns new Camera
func NewCamera(lookAtX, lookAtY, w, h float64) *Camera {
	target := vec2{lookAtX, lookAtY}
	c := &Camera{
		ZoomFactor:   1.0,
		Lerp:         false,
		LerpSpeed:    0.1,
		ShakeOptions: DefaultCameraShakeOptions(),
		// private
		w:               w,
		h:               h,
		angle:           0,
		zoomFactorShake: 1.0,
		trauma:          0,
		drawOptions:     &ebiten.DrawImageOptions{},
		centerOffset:    vec2{-(w * 0.5), -(h * 0.5)},
		traumaOffset:    vec2{},
		topLeft:         vec2{},
		tempTarget:      vec2{},
		tickSpeed:       1.0 / 60.0,
		tick:            0,
	}

	c.LookAt(lookAtX, lookAtY)
	c.tempTarget = target
	return c
}

func DefaultCameraShakeOptions() *CameraShakeOptions {
	opt := &CameraShakeOptions{
		Noise:         fastnoise.New[float64](),
		MaxX:          10.0,
		MaxY:          10.0,
		MaxAngle:      0.05,
		MaxZoomFactor: 0.1,
		Decay:         0.666,
		TimeScale:     10,
	}
	opt.Noise.Frequency = 0.5
	return opt
}

// LookAt aligns the midpoint of the camera viewport to the target.
// Use this function only once in Update() and change only the TargetX TargetY variables
func (cam *Camera) LookAt(targetX, targetY float64) {

	target := vec2{targetX, targetY}

	if cam.Lerp {
		cam.tempTarget = cam.tempTarget.Lerp(target, cam.LerpSpeed)
		cam.topLeft = cam.tempTarget
	} else {
		cam.topLeft = target
	}

	if cam.trauma > 0 {

		var shake = math.Pow(cam.trauma, 2)

		noiseValueX := cam.ShakeOptions.Noise.GetNoise3D(cam.tick*cam.ShakeOptions.TimeScale, 0, 0)
		noiseValueY := cam.ShakeOptions.Noise.GetNoise3D(0, cam.tick*cam.ShakeOptions.TimeScale, 0)
		noiseValueAngle := cam.ShakeOptions.Noise.GetNoise3D(0, 0, cam.tick*cam.ShakeOptions.TimeScale)

		cam.traumaOffset.X = noiseValueX * cam.ShakeOptions.MaxX * shake
		cam.traumaOffset.Y = noiseValueY * cam.ShakeOptions.MaxY * shake
		cam.actualAngle = noiseValueAngle * cam.ShakeOptions.MaxAngle * shake

		noiseValueZoom := cam.ShakeOptions.Noise.GetNoise3D(cam.tick*cam.ShakeOptions.TimeScale+300, 0, 0)
		cam.zoomFactorShake = (noiseValueZoom * cam.ShakeOptions.MaxZoomFactor * shake)
		cam.zoomFactorShake *= cam.ZoomFactor
		cam.zoomFactorShake += cam.ZoomFactor

		cam.trauma = clamp(cam.trauma-(cam.tickSpeed*cam.ShakeOptions.Decay), 0, 1)

	} else {
		cam.actualAngle = 0.0
		cam.zoomFactorShake = cam.ZoomFactor
	}

	// offset
	cam.actualAngle += cam.angle
	cam.topLeft = cam.topLeft.Add(cam.traumaOffset)
	cam.topLeft = cam.topLeft.Add(cam.centerOffset)

	// tick
	cam.tick += cam.tickSpeed
	if cam.tick > 1000000 {
		cam.tick = 0
	}
}

// AddTrauma adds trauma. factor is in the range [0-1]
func (cam *Camera) AddTrauma(factor float64) {
	cam.trauma = clamp(cam.trauma+factor, 0, 1)
}

// TopLeft() returns top left position of the camera rectangle
func (cam *Camera) TopLeft() (X float64, Y float64) {
	return cam.topLeft.X, cam.topLeft.Y
}

// Target returns center point of the camera in world-space
func (cam *Camera) Target() (X float64, Y float64) {
	center := cam.topLeft.Sub(cam.centerOffset)
	return center.X, center.Y
}

// ActualAngle returns camera rotation angle (including the angle of trauma shaking.). The unit is radian.
func (cam *Camera) ActualAngle() (angle float64) {
	return cam.actualAngle
}

// Angle returns camera rotation angle (The angle of trauma shake is not included.). The unit is radian.
func (cam *Camera) Angle() (angle float64) {
	return cam.angle
}

// SetAngle sets rotation. The unit is radian.
func (cam *Camera) SetAngle(angle float64) {
	cam.angle = angle
}

// Width returns width  of the camera
func (cam *Camera) Width() float64 {
	return cam.w
}

// Width returns height  of the camera
func (cam *Camera) Height() float64 {
	return cam.h
}

// SetSize ses camera rectangle size
func (cam *Camera) SetSize(w, h float64) {
	cam.w, cam.h = w, h
	cam.centerOffset = vec2{-(w * 0.5), -(h * 0.5)}
}

// Reset resets rotation and zoom factor to zero
func (cam *Camera) Reset() {
	cam.angle, cam.ZoomFactor, cam.zoomFactorShake = 0.0, 1.0, 1.0
}

// String returns camera values as string
func (cam *Camera) String() string {
	x, y := cam.Target()
	return fmt.Sprintf(
		"TargetX: %.1f\nTargetY: %.1f\nCam Rotation: %.1f\nZoom factor: %.2f\nLerp: %v",
		x, y, cam.ActualAngle(), cam.zoomFactorShake, cam.Lerp,
	)
}

// ScreenToWorld converts screen-space coordinates to world-space
func (cam *Camera) ScreenToWorld(screenX, screenY int) (worldX float64, worldY float64) {
	g := ebiten.GeoM{}
	cam.ApplyCameraTransform(&g)
	if g.IsInvertible() {
		g.Invert()
		worldX, worldY := g.Apply(float64(screenX), float64(screenY))
		return worldX, worldY
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

// ApplyCameraTransform applies geometric transformation to given geoM
func (cam *Camera) ApplyCameraTransform(geoM *ebiten.GeoM) {
	geoM.Translate(-cam.topLeft.X, -cam.topLeft.Y)                             // camera movement
	geoM.Translate(cam.centerOffset.X, cam.centerOffset.Y)                     // rotate and scale from center.
	geoM.Rotate(cam.actualAngle)                                               // rotate
	geoM.Scale(cam.zoomFactorShake, cam.zoomFactorShake)                       // apply zoom factor
	geoM.Translate(math.Abs(cam.centerOffset.X), math.Abs(cam.centerOffset.Y)) // restore center translation
}

// Draw applies the Camera's geometric transformation then draws the object on the screen with drawing options.
func (cam *Camera) Draw(worldObject *ebiten.Image, worldObjectOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	cam.drawOptions = worldObjectOps
	cam.ApplyCameraTransform(&cam.drawOptions.GeoM)
	screen.DrawImage(worldObject, cam.drawOptions)
	cam.drawOptions.GeoM.Reset()
}

type CameraShakeOptions struct {
	// Noise generator for noise types and settings.
	Noise         *fastnoise.State[float64]
	MaxX          float64 // Maximum X-axis shake. 0 means disabled
	MaxY          float64 // Maximum Y-axis shake. 0 means disabled
	MaxAngle      float64 // Max shake angle (radians). 0 means disabled
	MaxZoomFactor float64 // Zoom factor strength [1-0]. 0 means disabled
	TimeScale     float64 // Noise time domain speed
	Decay         float64 // Decay for trauma
}
