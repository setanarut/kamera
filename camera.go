package kamera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/fastnoise"
)

var statText = `TargetX: %.2f
TargetY: %.2f
Cam Rotation: %.2f
Zoom factor: %.2f
LerpEnabled: %v
ShakeEnabled: %v`

// Camera object.
//
// Use the `Camera.LookAt()` to align the center of the camera to the target.
type Camera struct {

	// ZoomFactor is the camera zoom (scaling) factor. Default is 1.
	ZoomFactor float64

	// If LerpEnabled is true, the camera moves with linear interpolation.
	//
	// The default value is false.
	LerpEnabled bool

	// LerpSpeed ​​is the interpolation speed in the range [0-1].
	//
	// The default value is 0.1.
	LerpSpeed float64

	// If ShakeEnabled is false, AddTrauma() has no effect and shake is always 0.
	//
	// The default value is false
	ShakeEnabled bool

	// ShakeOptions holds the camera shake options.
	ShakeOptions *CameraShakeOptions

	// private
	drawOptions                                                        *ebiten.DrawImageOptions
	angle, actualAngle, tickSpeed, tick, trauma, w, h, zoomFactorShake float64
	tempTarget, centerOffset, topLeft, traumaOffset                    vec2
	bb                                                                 BB
}

// NewCamera returns new Camera
func NewCamera(lookAtX, lookAtY, w, h float64) *Camera {
	target := vec2{lookAtX, lookAtY}
	c := &Camera{
		ZoomFactor:   1.0,
		LerpEnabled:  false,
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
	c.bb = NewBBForExtents(lookAtX, lookAtY, w/2, h/2)
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
// Use this function only once in Update() and change only the (targetX, targetY)
func (cam *Camera) LookAt(targetX, targetY float64) {

	target := vec2{targetX, targetY}

	if cam.LerpEnabled {
		cam.tempTarget = cam.tempTarget.Lerp(target, cam.LerpSpeed)
		cam.topLeft = cam.tempTarget
	} else {
		cam.topLeft = target
	}

	if cam.ShakeEnabled {
		if cam.trauma > 0 {

			var shake = math.Pow(cam.trauma, 2)

			noiseValueX := cam.ShakeOptions.Noise.GetNoise3D(
				cam.tick*cam.ShakeOptions.TimeScale,
				0,
				0,
			)
			noiseValueY := cam.ShakeOptions.Noise.GetNoise3D(
				0,
				cam.tick*cam.ShakeOptions.TimeScale,
				0,
			)
			noiseValueAngle := cam.ShakeOptions.Noise.GetNoise3D(
				0,
				0,
				cam.tick*cam.ShakeOptions.TimeScale,
			)

			cam.traumaOffset.X = noiseValueX * cam.ShakeOptions.MaxX * shake
			cam.traumaOffset.Y = noiseValueY * cam.ShakeOptions.MaxY * shake
			cam.actualAngle = noiseValueAngle * cam.ShakeOptions.MaxAngle * shake

			noiseValueZoom := cam.ShakeOptions.Noise.GetNoise3D(
				cam.tick*cam.ShakeOptions.TimeScale+300,
				0,
				0,
			)
			cam.zoomFactorShake = noiseValueZoom * cam.ShakeOptions.MaxZoomFactor * shake
			cam.zoomFactorShake *= cam.ZoomFactor
			cam.zoomFactorShake += cam.ZoomFactor

			cam.trauma = clamp(
				cam.trauma-(cam.tickSpeed*cam.ShakeOptions.Decay),
				0,
				1,
			)

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

	} else {
		cam.zoomFactorShake = cam.ZoomFactor
		cam.actualAngle = cam.angle
		cam.topLeft = cam.topLeft.Add(cam.centerOffset)
		cam.trauma = 0
		cam.traumaOffset = vec2{}
	}
}

// AddTrauma adds trauma. Factor is in the range [0-1]
func (cam *Camera) AddTrauma(factor float64) {
	if cam.ShakeEnabled {
		cam.trauma = clamp(cam.trauma+factor, 0, 1)
	}
}

// TopLeft returns top left position of the camera in world-space
func (cam *Camera) TopLeft() (X float64, Y float64) {
	return cam.topLeft.X, cam.topLeft.Y
}

// BB returns camera's bounding box in world-space
func (cam *Camera) BB() BB {
	x, y := cam.Center()
	return NewBBForExtents(x, y, cam.w*0.5, cam.h*0.5)
}

// Center returns center point of the camera in world-space
func (cam *Camera) Center() (X float64, Y float64) {
	center := cam.topLeft.Sub(cam.centerOffset)
	return center.X, center.Y
}

// ActualAngle returns camera rotation angle (including the angle of trauma shaking.).
//
// The unit is radian.
func (cam *Camera) ActualAngle() (angle float64) {
	return cam.actualAngle
}

// Angle returns camera rotation angle (The angle of trauma shake is not included.).
//
// The unit is radian.
func (cam *Camera) Angle() (angle float64) {
	return cam.angle
}

// SetAngle sets rotation. The unit is radian.
func (cam *Camera) SetAngle(angle float64) {
	cam.angle = angle
}

// Width returns width of the camera
func (cam *Camera) Width() float64 {
	return cam.w
}

// Height returns height of the camera
func (cam *Camera) Height() float64 {
	return cam.h
}

// SetSize sets camera rectangle size
func (cam *Camera) SetSize(w, h float64) {
	cx, cy := cam.Center()
	cam.w, cam.h = w, h
	cam.centerOffset = vec2{-(w * 0.5), -(h * 0.5)}
	cam.LookAt(cx, cy)
}

// Reset resets rotation and zoom factor to zero
func (cam *Camera) Reset() {
	cam.angle, cam.ZoomFactor, cam.zoomFactorShake = 0.0, 1.0, 1.0
}

// String returns camera values as string
func (cam *Camera) String() string {
	return fmt.Sprintf(
		statText,
		cam.topLeft.X-cam.centerOffset.X,
		cam.topLeft.Y-cam.centerOffset.Y,
		cam.actualAngle,
		cam.zoomFactorShake,
		cam.LerpEnabled,
		cam.ShakeEnabled,
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
