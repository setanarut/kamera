// Package kamera provides a camera object for Ebitengine v2.
package kamera

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/setanarut/fastnoise"
)

// SmoothType is the camera movement smoothing type.
type SmoothType int

const (
	// None is instant movement to the target. No smoothing.
	None SmoothType = iota
	// Lerp is Lerp() function.
	Lerp
	// SmoothDamp is SmoothDamp() function.
	SmoothDamp
)

// Camera object.
//
// Use the `Camera.LookAt()` to align the center of the camera to the target.
type Camera struct {

	// ZoomFactor is the camera zoom (scaling) factor. Default is 1.
	ZoomFactor float64

	// SmoothType is the camera movement smoothing type.
	SmoothType SmoothType

	// SmoothOptions holds the camera movement smoothing settings
	SmoothOptions *SmoothOptions

	// If ShakeEnabled is false, AddTrauma() has no effect and shake is always 0.
	//
	// The default value is false
	ShakeEnabled bool

	// ShakeOptions holds the camera shake options.
	ShakeOptions *ShakeOptions

	// private
	drawOptions                                                           *ebiten.DrawImageOptions
	drawOptionsCM                                                         *colorm.DrawImageOptions
	angle, actualAngle, tickSpeed, tick, trauma, w, h, zoomFactorShake    float64
	tempTargetX, centerOffsetX, topLeftX, traumaOffsetX, currentVelocityX float64
	tempTargetY, centerOffsetY, topLeftY, traumaOffsetY, currentVelocityY float64
}

// NewCamera returns new Camera
func NewCamera(lookAtX, lookAtY, w, h float64) *Camera {
	c := &Camera{
		ZoomFactor:    1.0,
		SmoothType:    None,
		SmoothOptions: DefaultSmoothOptions(),
		ShakeOptions:  DefaultCameraShakeOptions(),
		// private
		w:               w,
		h:               h,
		angle:           0,
		zoomFactorShake: 1.0,
		trauma:          0,
		drawOptions:     &ebiten.DrawImageOptions{},
		centerOffsetX:   -(w * 0.5),
		centerOffsetY:   -(h * 0.5),
		tickSpeed:       1.0 / 60.0,
		tick:            0,
	}

	c.LookAt(lookAtX, lookAtY)
	c.tempTargetX = lookAtX
	c.tempTargetY = lookAtY
	return c
}

func DefaultCameraShakeOptions() *ShakeOptions {
	opt := &ShakeOptions{
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

// smoothDamp gradually changes a value towards a desired goal over time.
func (cam *Camera) smoothDamp(targetX, targetY float64) (float64, float64) {
	// Ensure smooth times are not too small to avoid division by zero
	smoothTimeX := math.Max(0.0001, cam.SmoothOptions.SmoothDampTimeX)
	smoothTimeY := math.Max(0.0001, cam.SmoothOptions.SmoothDampTimeY)

	// Calculate exponential decay factors for X and Y
	omegaX := 2.0 / smoothTimeX
	omegaY := 2.0 / smoothTimeY

	xX := omegaX * 0.016666666666666666
	xY := omegaY * 0.016666666666666666

	expX := 1.0 / (1.0 + xX + 0.48*xX*xX + 0.235*xX*xX*xX)
	expY := 1.0 / (1.0 + xY + 0.48*xY*xY + 0.235*xY*xY*xY)

	// Calculate change with independent max speeds
	changeX := cam.tempTargetX - targetX
	changeY := cam.tempTargetY - targetY

	originalToX := targetX
	originalToY := targetY

	maxChangeX := cam.SmoothOptions.SmoothDampMaxSpeedX * smoothTimeX
	maxChangeY := cam.SmoothOptions.SmoothDampMaxSpeedY * smoothTimeY

	maxChangeXSq := maxChangeX * maxChangeX
	maxChangeYSq := maxChangeY * maxChangeY

	// Limit change independently for X and Y
	if changeX*changeX > maxChangeXSq {
		changeX = math.Copysign(maxChangeX, changeX)
	}

	if changeY*changeY > maxChangeYSq {
		changeY = math.Copysign(maxChangeY, changeY)
	}

	targetX = cam.tempTargetX - changeX
	targetY = cam.tempTargetY - changeY

	// Calculate velocity and output with independent exponential decay
	tempVelocityX := (cam.currentVelocityX + changeX*omegaX) * 0.016666666666666666
	tempVelocityY := (cam.currentVelocityY + changeY*omegaY) * 0.016666666666666666

	cam.currentVelocityX = (cam.currentVelocityX - tempVelocityX*omegaX) * expX
	cam.currentVelocityY = (cam.currentVelocityY - tempVelocityY*omegaY) * expY

	outputX := targetX + (changeX+tempVelocityX)*expX
	outputY := targetY + (changeY+tempVelocityY)*expY

	origMinusCurrentX := originalToX - cam.tempTargetX
	origMinusCurrentY := originalToY - cam.tempTargetY

	outMinusOrigX := outputX - originalToX
	outMinusOrigY := outputY - originalToY

	dot := origMinusCurrentX*outMinusOrigX + origMinusCurrentY*outMinusOrigY

	if dot > 0 {
		outputX = originalToX
		outputY = originalToY
		cam.currentVelocityX = (outputX - originalToX) / 0.016666666666666666
		cam.currentVelocityY = (outputY - originalToY) / 0.016666666666666666
	}

	return outputX, outputY
}

// LookAt aligns the midpoint of the camera viewport to the target.
// Use this function only once in Update() and change only the (targetX, targetY)
func (cam *Camera) LookAt(targetX, targetY float64) {
	switch cam.SmoothType {
	case SmoothDamp:
		cam.tempTargetX, cam.tempTargetY = cam.smoothDamp(targetX, targetY)
		cam.topLeftX = cam.tempTargetX
		cam.topLeftY = cam.tempTargetY

	case Lerp:
		cam.tempTargetX = lerp(cam.tempTargetX, targetX, cam.SmoothOptions.LerpSpeedX)
		cam.tempTargetY = lerp(cam.tempTargetY, targetY, cam.SmoothOptions.LerpSpeedY)
		cam.topLeftX, cam.topLeftY = cam.tempTargetX, cam.tempTargetY
	default:
		cam.topLeftX, cam.topLeftY = targetX, targetY
	}
	if cam.ShakeEnabled {
		if cam.trauma > 0 {
			var shake = math.Pow(cam.trauma, 2)
			noiseValueX := cam.ShakeOptions.Noise.GetNoise3D(cam.tick*cam.ShakeOptions.TimeScale, 0, 0)
			noiseValueY := cam.ShakeOptions.Noise.GetNoise3D(0, cam.tick*cam.ShakeOptions.TimeScale, 0)
			noiseValueAngle := cam.ShakeOptions.Noise.GetNoise3D(0, 0, cam.tick*cam.ShakeOptions.TimeScale)

			cam.traumaOffsetX = noiseValueX * cam.ShakeOptions.MaxX * shake
			cam.traumaOffsetY = noiseValueY * cam.ShakeOptions.MaxY * shake
			cam.actualAngle = noiseValueAngle * cam.ShakeOptions.MaxAngle * shake

			noiseValueZoom := cam.ShakeOptions.Noise.GetNoise3D(cam.tick*cam.ShakeOptions.TimeScale+300, 0, 0)
			cam.zoomFactorShake = noiseValueZoom * cam.ShakeOptions.MaxZoomFactor * shake
			cam.zoomFactorShake *= cam.ZoomFactor
			cam.zoomFactorShake += cam.ZoomFactor

			// clamp
			cam.trauma = min(max(cam.trauma-(cam.tickSpeed*cam.ShakeOptions.Decay), 0), 1)

		} else {
			cam.actualAngle = 0.0
			cam.zoomFactorShake = cam.ZoomFactor
		}

		// offset
		cam.actualAngle += cam.angle
		cam.topLeftX += cam.traumaOffsetX
		cam.topLeftY += cam.traumaOffsetY

		// tick
		cam.tick += cam.tickSpeed
		if cam.tick > 1000000 {
			cam.tick = 0
		}

	} else {
		cam.zoomFactorShake = cam.ZoomFactor
		cam.actualAngle = cam.angle

		cam.topLeftX += cam.centerOffsetX
		cam.topLeftY += cam.centerOffsetY

		cam.trauma = 0
		cam.traumaOffsetX, cam.traumaOffsetY = 0, 0
	}
}

// AddTrauma adds trauma. Factor is in the range [0-1]
func (cam *Camera) AddTrauma(factor float64) {
	if cam.ShakeEnabled {
		cam.trauma = min(max(cam.trauma+factor, 0), 1) // clamp
	}
}

// TopLeft returns top left position of the camera in world-space
func (cam *Camera) TopLeft() (X float64, Y float64) {
	return cam.topLeftX, cam.topLeftY
}

// Center returns center point of the camera in world-space
func (cam *Camera) Center() (X float64, Y float64) {
	return cam.topLeftX - cam.centerOffsetX, cam.topLeftY - cam.centerOffsetY
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
	cam.centerOffsetX = -(w * 0.5)
	cam.centerOffsetY = -(h * 0.5)
	cam.LookAt(cx, cy)
}

// Reset resets rotation and zoom factor to zero
func (cam *Camera) Reset() {
	cam.angle, cam.ZoomFactor, cam.zoomFactorShake = 0.0, 1.0, 1.0
}

const cameraStats = `TargetX: %.2f
TargetY: %.2f
Cam Rotation: %.2f
Zoom factor: %.2f
ShakeEnabled: %v
Smoothing Function: %s
LerpSpeedX: %.4f
LerpSpeedY: %.4f
SmoothDampTimeX: %.4f
SmoothDampTimeY: %.4f
SmoothDampMaxSpeedX: %.2f
SmoothDampMaxSpeedY: %.2f`

// String returns camera values as string
func (cam *Camera) String() string {
	smoothTypeStr := ""
	switch cam.SmoothType {
	case None:
		smoothTypeStr = "None"
	case Lerp:
		smoothTypeStr = "Lerp"
	case SmoothDamp:
		smoothTypeStr = "SmoothDamp"
	}

	return fmt.Sprintf(
		cameraStats,
		cam.topLeftX-cam.centerOffsetX,
		cam.topLeftY-cam.centerOffsetY,
		cam.actualAngle,
		cam.zoomFactorShake,
		cam.ShakeEnabled,
		smoothTypeStr,
		cam.SmoothOptions.LerpSpeedX,
		cam.SmoothOptions.LerpSpeedY,
		cam.SmoothOptions.SmoothDampTimeX,
		cam.SmoothOptions.SmoothDampTimeY,
		cam.SmoothOptions.SmoothDampMaxSpeedX,
		cam.SmoothOptions.SmoothDampMaxSpeedY,
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

// ApplyCameraTransformToPoint applies camera transformation to given point
func (cam *Camera) ApplyCameraTransformToPoint(x, y float64) (float64, float64) {
	geoM := ebiten.GeoM{}
	cam.ApplyCameraTransform(&geoM)
	return geoM.Apply(x, y)
}

// ApplyCameraTransform applies geometric transformation to given geoM
func (cam *Camera) ApplyCameraTransform(geoM *ebiten.GeoM) {
	geoM.Translate(-cam.topLeftX, -cam.topLeftY)                             // camera movement
	geoM.Translate(cam.centerOffsetX, cam.centerOffsetY)                     // rotate and scale from center.
	geoM.Rotate(cam.actualAngle)                                             // rotate
	geoM.Scale(cam.zoomFactorShake, cam.zoomFactorShake)                     // apply zoom factor
	geoM.Translate(math.Abs(cam.centerOffsetX), math.Abs(cam.centerOffsetY)) // restore center translation
}

// Draw applies the Camera's geometric transformation then draws the object on the screen with drawing options.
func (cam *Camera) Draw(worldObject *ebiten.Image, worldObjectOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	cam.drawOptions = worldObjectOps
	cam.ApplyCameraTransform(&cam.drawOptions.GeoM)
	screen.DrawImage(worldObject, cam.drawOptions)
	cam.drawOptions.GeoM.Reset()
}

// DrawWithColorM applies the Camera's geometric transformation then draws the object on the screen with colorm package drawing options.
func (cam *Camera) DrawWithColorM(worldObject *ebiten.Image, cm colorm.ColorM, worldObjectOps *colorm.DrawImageOptions, screen *ebiten.Image) {
	cam.drawOptionsCM = worldObjectOps
	cam.ApplyCameraTransform(&cam.drawOptionsCM.GeoM)
	colorm.DrawImage(screen, worldObject, cm, worldObjectOps)
	cam.drawOptions.GeoM.Reset()
}

type ShakeOptions struct {
	// Noise generator for noise types and settings.
	Noise         *fastnoise.State[float64]
	MaxX          float64 // Maximum X-axis shake. 0 means disabled
	MaxY          float64 // Maximum Y-axis shake. 0 means disabled
	MaxAngle      float64 // Max shake angle (radians). 0 means disabled
	MaxZoomFactor float64 // Zoom factor strength [1-0]. 0 means disabled
	TimeScale     float64 // Noise time domain speed
	Decay         float64 // Decay for trauma
}

// SmoothOptions is the camera movement smoothing options.
type SmoothOptions struct {
	// LerpSpeedX is the  X-axis linear interpolation speed every frame.
	// Value is in the range [0-1]. Default value is 0.09
	//
	// A smaller value will reach the target slower.
	LerpSpeedX float64
	// LerpSpeedY is the Y-axis linear interpolation speed every frame. Value is in the range [0-1].
	//
	// A smaller value will reach the target slower.
	LerpSpeedY float64

	// SmoothDampTimeX is the X-Axis approximate time it will take to reach the target.
	//
	// A smaller value will reach the target faster. Default value is 0.2
	SmoothDampTimeX float64
	// SmoothDampTimeY is the Y-Axis approximate time it will take to reach the target.
	//
	// A smaller value will reach the target faster. Default value is 0.2
	SmoothDampTimeY float64

	// SmoothDampMaxSpeedX is the maximum speed the camera can move while smooth damping in X-Axis
	//
	// Default value is 1000
	SmoothDampMaxSpeedX float64
	// SmoothDampMaxSpeedY is the maximum speed the camera can move while smooth damping in Y-Axis
	//
	// Default value is 1000
	SmoothDampMaxSpeedY float64
}

func DefaultSmoothOptions() *SmoothOptions {
	return &SmoothOptions{
		LerpSpeedX:          0.09,
		LerpSpeedY:          0.09,
		SmoothDampTimeX:     0.2,
		SmoothDampTimeY:     0.2,
		SmoothDampMaxSpeedX: 1000.0,
		SmoothDampMaxSpeedY: 1000.0,
	}
}

func lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}
