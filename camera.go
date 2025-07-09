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
	// Top-left X position of camera
	X float64
	// Top-left Y position of camera
	Y float64
	// Width is camera's width
	Width float64
	// Height is camera's height
	Height float64
	// Amgle is camera angle (without the angle of trauma shaking).
	//
	// The unit is radian.
	Angle float64
	// ActualAngle is camera angle (including the angle of trauma shaking).
	//
	// The unit is radian.
	ActualAngle float64
	// ZoomFactor is the camera zoom (scaling) factor. Default is 1.
	ZoomFactor float64
	// SmoothType is the camera movement smoothing type.
	SmoothType SmoothType
	// Trauma factor. Factor is in the range [0-1]. Use AddTrauma() function
	Trauma float64
	// SmoothOptions holds the camera movement smoothing settings
	SmoothOptions *SmoothOptions
	// ShakeOptions holds the camera shake options.
	ShakeOptions *ShakeOptions
	// If ShakeEnabled is false, AddTrauma() has no effect and shake is always 0.
	//
	// The default value is false
	ShakeEnabled bool
	// XAxisSmoothingDisabled disables the smoothing of the X axis if it's true.
	XAxisSmoothingDisabled bool
	// YAxisSmoothingDisabled disables the smoothing of the Y axis if it's true.
	YAxisSmoothingDisabled bool
	// Internal camera values. Do not change directly.
	TickSpeed, Tick, ZoomFactorShake float64
	// Internal camera values. Do not change directly.
	TempTargetX, CenterOffsetX, TraumaOffsetX, CurrentVelocityX float64
	// Internal camera values. Do not change directly.
	TempTargetY, CenterOffsetY, TraumaOffsetY, CurrentVelocityY float64
}

// NewCamera returns new Camera
func NewCamera(lookAtX, lookAtY, w, h float64) *Camera {
	c := &Camera{
		ZoomFactor:      1.0,
		SmoothType:      None,
		SmoothOptions:   DefaultSmoothOptions(),
		ShakeOptions:    DefaultCameraShakeOptions(),
		Width:           w,
		Height:          h,
		Angle:           0,
		ZoomFactorShake: 1.0,
		Trauma:          0,
		CenterOffsetX:   -(w * 0.5),
		CenterOffsetY:   -(h * 0.5),
		TickSpeed:       1.0 / 60.0,
		Tick:            0,
	}

	c.LookAt(lookAtX, lookAtY)
	c.TempTargetX = lookAtX
	c.TempTargetY = lookAtY
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

// smoothDampX gradually changes a value towards a desired goal over time for X axis.
func (cam *Camera) smoothDampX(targetX float64) float64 {
	// Ensure smooth time is not too small to avoid division by zero
	smoothTimeX := math.Max(0.0001, cam.SmoothOptions.SmoothDampTimeX)

	// Calculate exponential decay factor for X
	omegaX := 2.0 / smoothTimeX
	xX := omegaX * 0.016666666666666666
	expX := 1.0 / (1.0 + xX + 0.48*xX*xX + 0.235*xX*xX*xX)

	// Calculate change with max speed
	changeX := cam.TempTargetX - targetX
	originalToX := targetX
	maxChangeX := cam.SmoothOptions.SmoothDampMaxSpeedX * smoothTimeX
	maxChangeXSq := maxChangeX * maxChangeX

	// Limit change
	if changeX*changeX > maxChangeXSq {
		changeX = math.Copysign(maxChangeX, changeX)
	}

	targetX = cam.TempTargetX - changeX

	// Calculate velocity and output with exponential decay
	tempVelocityX := (cam.CurrentVelocityX + changeX*omegaX) * 0.016666666666666666
	cam.CurrentVelocityX = (cam.CurrentVelocityX - tempVelocityX*omegaX) * expX
	outputX := targetX + (changeX+tempVelocityX)*expX

	// Check if we've overshot the target
	origMinusCurrentX := originalToX - cam.TempTargetX
	outMinusOrigX := outputX - originalToX

	if origMinusCurrentX*outMinusOrigX > 0 {
		outputX = originalToX
		cam.CurrentVelocityX = (outputX - originalToX) / 0.016666666666666666
	}

	return outputX
}

// smoothDampY gradually changes a value towards a desired goal over time for Y axis.
func (cam *Camera) smoothDampY(targetY float64) float64 {
	// Ensure smooth time is not too small to avoid division by zero
	smoothTimeY := math.Max(0.0001, cam.SmoothOptions.SmoothDampTimeY)

	// Calculate exponential decay factor for Y
	omegaY := 2.0 / smoothTimeY
	xY := omegaY * 0.016666666666666666
	expY := 1.0 / (1.0 + xY + 0.48*xY*xY + 0.235*xY*xY*xY)

	// Calculate change with max speed
	changeY := cam.TempTargetY - targetY
	originalToY := targetY
	maxChangeY := cam.SmoothOptions.SmoothDampMaxSpeedY * smoothTimeY
	maxChangeYSq := maxChangeY * maxChangeY

	// Limit change
	if changeY*changeY > maxChangeYSq {
		changeY = math.Copysign(maxChangeY, changeY)
	}

	targetY = cam.TempTargetY - changeY

	// Calculate velocity and output with exponential decay
	tempVelocityY := (cam.CurrentVelocityY + changeY*omegaY) * 0.016666666666666666
	cam.CurrentVelocityY = (cam.CurrentVelocityY - tempVelocityY*omegaY) * expY
	outputY := targetY + (changeY+tempVelocityY)*expY

	// Check if we've overshot the target
	origMinusCurrentY := originalToY - cam.TempTargetY
	outMinusOrigY := outputY - originalToY

	if origMinusCurrentY*outMinusOrigY > 0 {
		outputY = originalToY
		cam.CurrentVelocityY = (outputY - originalToY) / 0.016666666666666666
	}

	return outputY
}

// LookAt aligns the midpoint of the camera viewport to the target.
//
// Camera motion smoothing is only applied with this method.
// Use this function only once in Update() and change only the (targetX, targetY)
func (cam *Camera) LookAt(targetX, targetY float64) {
	switch cam.SmoothType {
	case SmoothDamp:
		if !cam.XAxisSmoothingDisabled && !cam.YAxisSmoothingDisabled {
			cam.TempTargetX = cam.smoothDampX(targetX)
			cam.TempTargetY = cam.smoothDampY(targetY)
			cam.X = cam.TempTargetX
			cam.Y = cam.TempTargetY
		} else if !cam.XAxisSmoothingDisabled && cam.YAxisSmoothingDisabled {
			cam.TempTargetX = cam.smoothDampX(targetX)
			cam.X = cam.TempTargetX
			cam.Y = targetY
		} else if cam.XAxisSmoothingDisabled && !cam.YAxisSmoothingDisabled {
			cam.TempTargetY = cam.smoothDampY(targetY)
			cam.Y = cam.TempTargetY
			cam.X = targetX
		} else {
			cam.X = targetX
			cam.Y = targetY
		}
	case Lerp:
		if !cam.XAxisSmoothingDisabled && !cam.YAxisSmoothingDisabled {
			cam.TempTargetX = lerp(cam.TempTargetX, targetX, cam.SmoothOptions.LerpSpeedX)
			cam.TempTargetY = lerp(cam.TempTargetY, targetY, cam.SmoothOptions.LerpSpeedY)
			cam.X = cam.TempTargetX
			cam.Y = cam.TempTargetY
		} else if !cam.XAxisSmoothingDisabled && cam.YAxisSmoothingDisabled {
			cam.TempTargetX = lerp(cam.TempTargetX, targetX, cam.SmoothOptions.LerpSpeedX)
			cam.X = cam.TempTargetX
			cam.Y = targetY
		} else if cam.XAxisSmoothingDisabled && !cam.YAxisSmoothingDisabled {
			cam.TempTargetY = lerp(cam.TempTargetY, targetY, cam.SmoothOptions.LerpSpeedY)
			cam.Y = cam.TempTargetY
			cam.X = targetX
		} else {
			cam.X = targetX
			cam.Y = targetY
		}
	case None:
		cam.X = targetX
		cam.Y = targetY
	default:
		cam.X = targetX
		cam.Y = targetY
	}
	if cam.ShakeEnabled {
		if cam.Trauma > 0 {
			var shake = math.Pow(cam.Trauma, 2)
			noiseValueX := cam.ShakeOptions.Noise.GetNoise3D(cam.Tick*cam.ShakeOptions.TimeScale, 0, 0)
			noiseValueY := cam.ShakeOptions.Noise.GetNoise3D(0, cam.Tick*cam.ShakeOptions.TimeScale, 0)
			noiseValueAngle := cam.ShakeOptions.Noise.GetNoise3D(0, 0, cam.Tick*cam.ShakeOptions.TimeScale)

			cam.TraumaOffsetX = noiseValueX * cam.ShakeOptions.MaxX * shake
			cam.TraumaOffsetY = noiseValueY * cam.ShakeOptions.MaxY * shake
			cam.ActualAngle = noiseValueAngle * cam.ShakeOptions.MaxAngle * shake

			noiseValueZoom := cam.ShakeOptions.Noise.GetNoise3D(cam.Tick*cam.ShakeOptions.TimeScale+300, 0, 0)
			cam.ZoomFactorShake = noiseValueZoom * cam.ShakeOptions.MaxZoomFactor * shake
			cam.ZoomFactorShake *= cam.ZoomFactor
			cam.ZoomFactorShake += cam.ZoomFactor

			// clamp
			cam.Trauma = min(max(cam.Trauma-(cam.TickSpeed*cam.ShakeOptions.Decay), 0), 1)

		} else {
			cam.ActualAngle = 0.0
			cam.ZoomFactorShake = cam.ZoomFactor
		}

		// offset
		cam.ActualAngle += cam.Angle
		cam.X += cam.TraumaOffsetX
		cam.Y += cam.TraumaOffsetY
		cam.X += cam.CenterOffsetX
		cam.Y += cam.CenterOffsetY

		// tick
		cam.Tick += cam.TickSpeed
		if cam.Tick > 1000000 {
			cam.Tick = 0
		}

	} else {
		cam.ZoomFactorShake = cam.ZoomFactor
		cam.ActualAngle = cam.Angle

		cam.X += cam.CenterOffsetX
		cam.Y += cam.CenterOffsetY

		cam.Trauma = 0
		cam.TraumaOffsetX, cam.TraumaOffsetY = 0, 0
	}
}

// AddTrauma adds trauma. Factor is in the range [0-1]
func (cam *Camera) AddTrauma(factor float64) {
	if cam.ShakeEnabled {
		cam.Trauma = min(max(cam.Trauma+factor, 0), 1) // clamp
	}
}

// Right returns the right edge position of the camera in world-space.
func (cam *Camera) Right() float64 {
	return cam.X + cam.Width
}

// Bottom returns the bottom edge position of the camera in world-space.
func (cam *Camera) Bottom() float64 {
	return cam.Y + cam.Height
}

// SetTopLeft sets top-left position of the camera in world-space.
//
// Unlike the LookAt() method, the position is set directly (teleport).
func (cam *Camera) SetTopLeft(x, y float64) {
	cam.X, cam.Y = x, y
	cam.TempTargetX, cam.TempTargetY = cam.Center()

}

// SetCenter sets center position of the camera in world-space.
//
// Unlike the LookAt() method, the position is set directly (teleport).
//
// Can be used to cancel follow camera and teleport to target.
func (cam *Camera) SetCenter(x, y float64) {
	cam.TempTargetX, cam.TempTargetY = x, y
	cam.LookAt(x, y)
}

// Center returns center point of the camera in world-space
func (cam *Camera) Center() (X float64, Y float64) {
	return cam.X - cam.CenterOffsetX, cam.Y - cam.CenterOffsetY
}

// SetSize sets camera rectangle size
func (cam *Camera) SetSize(w, h float64) {
	cam.Width, cam.Height = w, h
	cam.CenterOffsetX = -(w * 0.5)
	cam.CenterOffsetY = -(h * 0.5)
}

// Reset resets rotation and zoom factor to zero
func (cam *Camera) Reset() {
	cam.Angle, cam.ZoomFactor, cam.ZoomFactorShake = 0.0, 1.0, 1.0
}

const cameraStats = `TargetX: %.2f
TargetY: %.2f
Top-left X: %.2f
Top-left Y: %.2f
Size: %.2f %.2f
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
		cam.X-cam.CenterOffsetX,
		cam.Y-cam.CenterOffsetY,
		cam.X,
		cam.Y,
		cam.Width, cam.Height,
		cam.ActualAngle,
		cam.ZoomFactorShake,
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
func (cam *Camera) ApplyCameraTransform(g *ebiten.GeoM) {
	g.Translate(-cam.X, -cam.Y)                                           // camera movement
	g.Translate(cam.CenterOffsetX, cam.CenterOffsetY)                     // rotate and scale from center.
	g.Rotate(cam.ActualAngle)                                             // rotate
	g.Scale(cam.ZoomFactorShake, cam.ZoomFactorShake)                     // apply zoom factor
	g.Translate(math.Abs(cam.CenterOffsetX), math.Abs(cam.CenterOffsetY)) // restore center translation
}

// Draw applies the Camera's geometric transformation then draws the object on the screen with drawing options.
func (cam *Camera) Draw(worldObject *ebiten.Image, worldObjectOps *ebiten.DrawImageOptions, screen *ebiten.Image) {
	cam.ApplyCameraTransform(&worldObjectOps.GeoM)
	screen.DrawImage(worldObject, worldObjectOps)
}

// DrawWithColorM applies the Camera's geometric transformation then draws the object on the screen with colorm package drawing options.
func (cam *Camera) DrawWithColorM(worldObject *ebiten.Image, cm colorm.ColorM, worldObjectOps *colorm.DrawImageOptions, screen *ebiten.Image) {
	cam.ApplyCameraTransform(&worldObjectOps.GeoM)
	colorm.DrawImage(screen, worldObject, cm, worldObjectOps)
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
