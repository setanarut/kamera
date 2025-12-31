package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/kamera/v2/examples"
	"github.com/setanarut/v"
	"golang.org/x/image/colornames"
)

const CameraDeadZoneX = 40.0
const CameraDeadZoneOffset = 100.0

const (
	ScreenWidth   = 620
	ScreenHeight  = 360
	MoveSpeedX    = 6.125
	JumpPower     = -14.46
	Gravity       = 0.66
	PlatformSpeed = 0.05
)

//go:embed bg.png
var bgImageBytes []byte

var cam *kamera.Camera

var (
	player        = coll.NewAABB(1000, ScreenHeight, 12, 24)
	playerDelta   = v.Vec{0, 0}
	playerHitInfo = &coll.Hit{}
)

var (
	bg    *ebiten.Image
	im    = ebiten.NewImage(1, 1)
	dio   = &colorm.DrawImageOptions{}
	bgdio = &ebiten.DrawImageOptions{}
	clrm  = colorm.ColorM{}
)
var (
	platform               = coll.NewAABB(1000, 130, 64, 16)
	platformDelta          = v.Vec{}
	platformRotationCenter = platform.Pos
	platformRadius         = 60.0
	platformAngle          = 0.0
	floorY                 float64
)

func init() {

	im.Fill(color.White)
	var err error
	bg, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(bgImageBytes))
	if err != nil {
		fmt.Println(err)
	}
	floorY = float64(bg.Bounds().Dy()) - 70
	player.SetBottom(floorY)

	// init camera
	cam = kamera.NewCamera(player.Pos.X, ScreenHeight/2, ScreenWidth, ScreenHeight)
	cam.ShakeEnabled = true
	cam.SmoothType = kamera.SmoothDamp
	cam.SmoothOptions.SmoothDampTimeX = 0.15

}
func calculateDeadZoneX(playerX float64) float64 {
	camCenterX := cam.X - cam.CenterOffsetX
	deltaX := playerX - camCenterX
	targetX := playerX
	if CameraDeadZoneX > 0 {
		if math.Abs(deltaX) <= CameraDeadZoneX {
			targetX = camCenterX
		} else {
			if deltaX > 0 {
				targetX = playerX - CameraDeadZoneX
			} else {
				targetX = playerX + CameraDeadZoneX
			}
		}
	}
	return targetX
}

type Game struct{}

func (g *Game) Update() error {

	cam.LookAt(calculateDeadZoneX(player.Pos.X+CameraDeadZoneOffset), ScreenHeight/2)

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		cam.AddTrauma(1.0)
	}

	movePlatform()
	playerDelta.Y += Gravity
	speed := examples.Axis().Unit().Scale(MoveSpeedX)

	playerHitInfo.Reset()
	hit := coll.BoxBoxSweep2(platform, player, platformDelta, playerDelta, playerHitInfo)
	onGround := false
	if hit && playerHitInfo.Normal.Y == -1 {
		onGround = isOnPlatform(playerHitInfo)
	}
	if onGround {
		player.Pos = player.Pos.Add(playerDelta.Scale(playerHitInfo.Data))
		player.Pos.Y = platform.Pos.Y + platformDelta.Y - player.Half.Y - platform.Half.Y
		player.Pos.X += platformDelta.X + speed.X
		playerDelta.X = platformDelta.X + speed.X
		playerDelta.Y = platformDelta.Y
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			playerDelta.Y = JumpPower
		}
	} else {
		playerDelta.X = speed.X
		player.Pos = player.Pos.Add(playerDelta)
	}

	platform.Pos = platform.Pos.Add(platformDelta)

	if player.Bottom() >= floorY {
		player.SetBottom(floorY)
		playerDelta.Y = 0
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			playerDelta.Y = JumpPower
		}
	}

	// inherit only the platform's fractional X offset to avoid sub-pixel jitter.
	if onGround {
		player.Pos.X = math.Floor(player.Pos.X) + Fract(platform.Pos.X)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{10})
	for x := range 3 {
		bgdio.GeoM.Reset()
		bgdio.GeoM.Translate(float64(x)*float64(bg.Bounds().Dx()), 0)
		cam.Draw(bg, bgdio, screen)
	}
	fillAABB(platform, screen, color.Gray{100})
	fillAABB(player, screen, color.Gray{180})

	//draw deadZone guides
	centerX := (float32(ScreenWidth) / 2)
	r := centerX + CameraDeadZoneX - CameraDeadZoneOffset
	l := centerX - CameraDeadZoneX - CameraDeadZoneOffset
	vector.StrokeLine(screen, r, 0, r, 1000, 1, colornames.Cyan, true)
	vector.StrokeLine(screen, l, 0, l, 1000, 1, colornames.Cyan, true)

	ebitenutil.DebugPrintAt(screen, "Space - Jump\nA/D - Move\nT - Trauma", 10, 10)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Platformer Example")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(fmt.Errorf("error running game: %w", err))
	}

}
func fillAABB(box *coll.AABB, s *ebiten.Image, c color.Color) {
	clrm.Reset()
	clrm.ScaleWithColor(c)
	dio.GeoM.Reset()
	dio.GeoM.Scale(box.Half.X*2, box.Half.Y*2)
	dio.GeoM.Translate(-box.Half.X, -box.Half.Y)
	dio.GeoM.Translate(box.Pos.X, box.Pos.Y)
	cam.DrawWithColorM(im, clrm, dio, s)
}

// Fract returns the fractional part of x.
func Fract(x float64) float64 {
	if x >= 0 {
		return x - math.Floor(x)
	}
	return x - math.Ceil(x)
}

func isOnPlatform(hit *coll.Hit) bool {
	playerPosAtHit := player.Pos.Add(playerDelta.Scale(hit.Data))
	platformPosAtHit := platform.Pos.Add(platformDelta.Scale(hit.Data))
	playerBottomAtHit := playerPosAtHit.Y + player.Half.Y
	platformTopAtHit := platformPosAtHit.Y - platform.Half.Y
	return playerBottomAtHit <= platformTopAtHit
}

func movePlatform() {
	platformAngle += PlatformSpeed
	newPlatCenterX := platformRotationCenter.X + math.Cos(platformAngle)*platformRadius
	newPlatCenterY := platformRotationCenter.Y + math.Sin(platformAngle)*platformRadius
	newPlatPos := v.Vec{X: newPlatCenterX, Y: newPlatCenterY}
	platformDelta = newPlatPos.Sub(platform.Pos)
}
