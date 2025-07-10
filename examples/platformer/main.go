package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/tilecollider"
)

const crossHairLength float32 = 90.0

var helpText = `CAMERA CONTROLS
E/Q - Zoom in/out
R - Rotate  
T - Add 1.0 trauma
Arrow Keys - Decrease/Increase camera smoothing speed
Tab - Change camera smoothing type
X - Enable/Disable Shake
Backspace - Reset camera

PLAYER CONTROLS
WASD - Move player    
Space - Jump
Shift - Run
`

var tileMap = [][]uint8{
	{1, 0, 1, 0, 1, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 0, 0, 0, 1, 0, 1},
	{0, 0, 0, 1, 0, 1, 0, 1},
	{0, 0, 0, 0, 0, 1, 0, 1},
	{1, 0, 1, 1, 1, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1}}

var (
	screenWidth, screenHeight = 854, 480
	offset                    = [2]int{0, 0}
	gridSize                  = [2]int{8, 8}
	tileSize                  = [2]int{64, 64}
	playerBox                 = [4]float64{70, 70, 24, 32} // x, y, w, h
	vel                       = [2]float64{0, 4}
	cam                       = kamera.NewCamera(playerBox[0], playerBox[1], float64(screenWidth), float64(screenHeight))
	controller                = NewPlayerController()
	collider                  = tilecollider.NewCollider(tileMap, tileSize[0], tileSize[1])
	tileDIO                   = &ebiten.DrawImageOptions{}
	playerDIO                 = &ebiten.DrawImageOptions{}
	tileImage                 = ebiten.NewImage(tileSize[0], tileSize[1])
	playerImage               = ebiten.NewImage(24, 32)
)

func init() {
	cam.SmoothType = kamera.SmoothDamp
	cam.ShakeEnabled = true

	tileImage.Fill(color.RGBA{0, 0, 255, 0})
	playerImage.Fill(color.Gray{100})
	controller.SetPhyicsScale(2.2)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func Translate(bx *[4]float64, x, y float64) {
	bx[0] += x
	bx[1] += y
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		cam.Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		cam.Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		cam.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		cam.ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		cam.ZoomFactor *= 1.02
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		cam.ShakeEnabled = !cam.ShakeEnabled
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		cam.AddTrauma(1.0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		switch cam.SmoothType {
		case kamera.Lerp:
			cam.SmoothOptions.LerpSpeedY -= 0.01
			cam.SmoothOptions.LerpSpeedY = max(0, min(cam.SmoothOptions.LerpSpeedY, 1))

		case kamera.SmoothDamp:
			cam.SmoothOptions.SmoothDampTimeY += 0.01
			cam.SmoothOptions.SmoothDampTimeY = max(0, min(cam.SmoothOptions.SmoothDampTimeY, 10))

		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		switch cam.SmoothType {
		case kamera.Lerp:
			cam.SmoothOptions.LerpSpeedY += 0.01
			cam.SmoothOptions.LerpSpeedY = max(0, min(cam.SmoothOptions.LerpSpeedY, 1))
		case kamera.SmoothDamp:
			cam.SmoothOptions.SmoothDampTimeY -= 0.01
			cam.SmoothOptions.SmoothDampTimeY = max(0, min(cam.SmoothOptions.SmoothDampTimeY, 10))
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		switch cam.SmoothType {
		case kamera.Lerp:
			cam.SmoothOptions.LerpSpeedX -= 0.01
			cam.SmoothOptions.LerpSpeedX = max(0, min(cam.SmoothOptions.LerpSpeedX, 1))

		case kamera.SmoothDamp:
			cam.SmoothOptions.SmoothDampTimeX += 0.01
			cam.SmoothOptions.SmoothDampTimeX = max(0, min(cam.SmoothOptions.SmoothDampTimeX, 10))

		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		switch cam.SmoothType {
		case kamera.Lerp:
			cam.SmoothOptions.LerpSpeedX += 0.01
			cam.SmoothOptions.LerpSpeedX = max(0, min(cam.SmoothOptions.LerpSpeedX, 1))
		case kamera.SmoothDamp:
			cam.SmoothOptions.SmoothDampTimeX -= 0.01
			cam.SmoothOptions.SmoothDampTimeX = max(0, min(cam.SmoothOptions.SmoothDampTimeX, 10))
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		switch cam.SmoothType {
		case kamera.None:
			cam.SetCenter(playerBox[0]+playerBox[2]/2, playerBox[1]+playerBox[3]/2)
			cam.SmoothType = kamera.Lerp
		case kamera.Lerp:
			cam.SetCenter(playerBox[0]+playerBox[2]/2, playerBox[1]+playerBox[3]/2)
			cam.SmoothType = kamera.SmoothDamp
		case kamera.SmoothDamp:
			cam.SetCenter(playerBox[0]+playerBox[2]/2, playerBox[1]+playerBox[3]/2)
			cam.SmoothType = kamera.None
		}
	}

	if vel[1] < 0 {
		controller.IsOnFloor = false
	}
	vel = controller.ProcessVelocity(vel)
	dx, dy := collider.Collide(
		playerBox[0],
		playerBox[1],
		playerBox[2],
		playerBox[3],
		vel[0],
		vel[1],
		func(infos []tilecollider.CollisionInfo[uint8], _, _ float64) {
			for _, info := range infos {
				if info.Normal[1] == -1 {
					controller.IsOnFloor = true
				}
				if info.Normal[1] == 1 {
					controller.IsJumping = false
					vel[1] = 0
				}
				if info.Normal[0] == 1 || info.Normal[0] == -1 {
					vel[0] = 0
				}
			}
		},
	)

	Translate(&playerBox, dx, dy)

	// Update camera
	cam.LookAt(playerBox[0]+playerBox[2]/2, playerBox[1]+playerBox[3]/2)

	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

type Game struct{}

func (g *Game) Draw(s *ebiten.Image) {

	// Draw tiles
	for y, row := range tileMap {
		for x, value := range row {
			if value != 0 {
				tileDIO.GeoM.Reset()
				tileDIO.GeoM.Translate(float64(x*tileSize[0]), float64(y*tileSize[1]))
				cam.Draw(tileImage, tileDIO, s)
			}
		}
	}

	// Draw player
	playerDIO.GeoM.Reset()
	playerDIO.GeoM.Translate(playerBox[0], playerBox[1])
	cam.Draw(playerImage, playerDIO, s)

	// Draw camera crosshair
	cx, cy := float32(screenWidth/2), float32(screenHeight/2)
	vector.StrokeLine(s, cx-crossHairLength, cy, cx+crossHairLength, cy, 1, color.RGBA{255, 255, 0, 255}, false)
	vector.StrokeLine(s, cx, cy-crossHairLength, cx, cy+crossHairLength, 1, color.RGBA{255, 255, 0, 255}, false)

	// Draw help text
	ebitenutil.DebugPrintAt(s, helpText, 14, 0)
	ebitenutil.DebugPrintAt(s, "CAMERA STATS", 14, 220)
	ebitenutil.DebugPrintAt(s, cam.String(), 14, 235)
}

func Axis() (axisX, axisY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axisY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axisY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axisX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axisX += 1
	}
	return axisX, axisY
}

type PlayerController struct {
	MinSpeed         float64
	MaxSpeed         float64
	MaxWalkSpeed     float64
	MaxFallSpeed     float64
	MaxFallSpeedCap  float64
	MinSlowDownSpeed float64
	WalkAcceleration float64
	RunAcceleration  float64
	WalkFriction     float64
	SkidFriction     float64
	StompSpeed       float64
	StompSpeedCap    float64
	JumpSpeed        [3]float64
	LongJumpGravity  [3]float64
	Gravity          float64
	SpeedThresholds  [2]float64
	// states
	IsFacingLeft bool
	IsRunning    bool
	IsJumping    bool
	IsFalling    bool
	IsSkidding   bool
	IsCrouching  bool
	IsOnFloor    bool
	// private
	minSpeedValue       float64
	maxSpeedValue       float64
	accel               float64
	speedThresholdIndex int
}

func NewPlayerController() *PlayerController {
	const (
		minSpeed              = 0.07421875
		maxSpeed              = 2.5625
		maxWalkSpeed          = 1.5625
		maxFallSpeed          = 4.5
		maxFallSpeedCap       = 4
		minSlowDownSpeed      = 0.5625
		walkAcceleration      = 0.037109375
		runAcceleration       = 0.0556640625
		walkFriction          = 0.05078125
		skidFriction          = 0.1015625
		stompSpeed            = 4
		stompSpeedCap         = 4
		jumpSpeedNormal       = -4
		jumpSpeedRun          = -4
		jumpSpeedLong         = -5
		longJumpGravityNormal = 0.12
		longJumpGravityRun    = 0.11
		longJumpGravityLong   = 0.15
		gravity               = 0.43
		speedThreshold1       = 1
		speedThreshold2       = 2.3125
	)

	pc := &PlayerController{
		MinSpeed:         minSpeed,
		MaxSpeed:         maxSpeed,
		MaxWalkSpeed:     maxWalkSpeed,
		MaxFallSpeed:     maxFallSpeed,
		MaxFallSpeedCap:  maxFallSpeedCap,
		MinSlowDownSpeed: minSlowDownSpeed,
		WalkAcceleration: walkAcceleration,
		RunAcceleration:  runAcceleration,
		WalkFriction:     walkFriction,
		SkidFriction:     skidFriction,
		StompSpeed:       stompSpeed,
		StompSpeedCap:    stompSpeedCap,

		JumpSpeed:       [3]float64{jumpSpeedNormal, jumpSpeedRun, jumpSpeedLong},
		LongJumpGravity: [3]float64{longJumpGravityNormal, longJumpGravityRun, longJumpGravityLong},
		Gravity:         gravity,
		SpeedThresholds: [2]float64{speedThreshold1, speedThreshold2},
		IsFacingLeft:    false,
		IsRunning:       false,

		IsJumping:   false,
		IsFalling:   false,
		IsSkidding:  false,
		IsCrouching: false,
		IsOnFloor:   false,

		speedThresholdIndex: 0,
	}

	pc.minSpeedValue = pc.MinSpeed
	pc.maxSpeedValue = pc.MaxSpeed
	pc.accel = pc.WalkAcceleration

	return pc
}

func (pc *PlayerController) SetPhyicsScale(s float64) {
	pc.MinSpeed *= s
	pc.MaxSpeed *= s
	pc.MaxWalkSpeed *= s
	pc.MaxFallSpeed *= s
	pc.MaxFallSpeedCap *= s
	pc.MinSlowDownSpeed *= s
	pc.WalkAcceleration *= s
	pc.RunAcceleration *= s
	pc.WalkFriction *= s
	pc.SkidFriction *= s
	pc.StompSpeed *= s
	pc.StompSpeedCap *= s
	pc.JumpSpeed[0] *= s
	pc.JumpSpeed[1] *= s
	pc.JumpSpeed[2] *= s
	pc.LongJumpGravity[0] *= s
	pc.LongJumpGravity[1] *= s
	pc.LongJumpGravity[2] *= s
	pc.Gravity *= s
	pc.SpeedThresholds[0] *= s
	pc.SpeedThresholds[1] *= s
}

func (pc *PlayerController) ProcessVelocity(vel [2]float64) [2]float64 {
	inputAxisX, inputAxisY := Axis()

	if pc.IsOnFloor {
		pc.IsRunning = ebiten.IsKeyPressed(ebiten.KeyShift)
		pc.IsCrouching = ebiten.IsKeyPressed(ebiten.KeyDown)
		if pc.IsCrouching && inputAxisX != 0 {
			pc.IsCrouching = false
			inputAxisX = 0.0
		}
	}

	if pc.IsOnFloor {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			pc.IsJumping = true
			speed := math.Abs(vel[0])
			pc.speedThresholdIndex = 0
			if speed >= pc.SpeedThresholds[1] {
				pc.speedThresholdIndex = 2
			} else if speed >= pc.SpeedThresholds[0] {
				pc.speedThresholdIndex = 1
			}

			vel[1] = pc.JumpSpeed[pc.speedThresholdIndex]

		}
	} else {
		gravityValue := pc.Gravity
		if ebiten.IsKeyPressed(ebiten.KeySpace) && pc.IsJumping && vel[1] < 0 {
			gravityValue = pc.LongJumpGravity[pc.speedThresholdIndex]
		}
		vel[1] += gravityValue
		if vel[1] > pc.MaxFallSpeedCap {
			vel[1] = pc.MaxFallSpeedCap
		}
	}

	// Update states
	if vel[1] > 0 {
		pc.IsJumping = false
		pc.IsFalling = true
	} else if pc.IsOnFloor {
		pc.IsFalling = false
	}

	if inputAxisX != 0 {
		if pc.IsOnFloor {
			if vel[0] != 0 {
				pc.IsFacingLeft = inputAxisX < 0.0
				pc.IsSkidding = vel[0] < 0.0 != pc.IsFacingLeft
			}
			if pc.IsSkidding {
				pc.minSpeedValue = pc.MinSlowDownSpeed
				pc.maxSpeedValue = pc.MaxWalkSpeed
				pc.accel = pc.SkidFriction
			} else if pc.IsRunning {
				pc.minSpeedValue = pc.MinSpeed
				pc.maxSpeedValue = pc.MaxSpeed
				pc.accel = pc.RunAcceleration
			} else {
				pc.minSpeedValue = pc.MinSpeed
				pc.maxSpeedValue = pc.MaxWalkSpeed
				pc.accel = pc.WalkAcceleration
			}
		} else if pc.IsRunning && math.Abs(vel[0]) > pc.MaxWalkSpeed {
			pc.maxSpeedValue = pc.MaxSpeed
		} else {
			pc.maxSpeedValue = pc.MaxWalkSpeed
		}
		targetSpeed := inputAxisX * pc.maxSpeedValue

		// Manually implementing moveToward()
		if vel[0] < targetSpeed {
			vel[0] += pc.accel
			if vel[0] > targetSpeed {
				vel[0] = targetSpeed
			}
		} else if vel[0] > targetSpeed {
			vel[0] -= pc.accel
			if vel[0] < targetSpeed {
				vel[0] = targetSpeed
			}
		}

	} else if pc.IsOnFloor && vel[0] != 0 {
		if !pc.IsSkidding {
			pc.accel = pc.WalkFriction
		}
		if inputAxisY != 0 {
			pc.minSpeedValue = pc.MinSlowDownSpeed
		} else {
			pc.minSpeedValue = pc.MinSpeed
		}
		if math.Abs(vel[0]) < pc.minSpeedValue {
			vel[0] = 0.0
		} else {
			if vel[0] > 0 {
				vel[0] -= pc.accel
				if vel[0] < 0 {
					vel[0] = 0
				}
			} else {
				vel[0] += pc.accel
				if vel[0] > 0 {
					vel[0] = 0
				}
			}
		}
	}
	if math.Abs(vel[0]) < pc.MinSlowDownSpeed {
		pc.IsSkidding = false
	}

	return vel
}
