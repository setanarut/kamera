package examples

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/setanarut/coll"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/v"
)

var im = ebiten.NewImage(1, 1)
var dio = &ebiten.DrawImageOptions{}

func init() {
	im.Fill(color.White)
}

const WindowWidth, WindowHeight int = 500, 500

func StrokeCircleAt(dst *ebiten.Image, pos v.Vec, r float64, clr color.Color) {
	vector.StrokeCircle(dst, float32(pos.X), float32(pos.Y), float32(r), 2, clr, true)
}
func StrokeCircle(dst *ebiten.Image, c *coll.Circle, clr color.Color) {
	vector.StrokeCircle(dst, float32(c.Pos.X), float32(c.Pos.Y), float32(c.Radius), 2, clr, true)
}

func FillCircle(dst *ebiten.Image, c *coll.Circle, clr color.Color) {
	vector.FillCircle(dst, float32(c.Pos.X), float32(c.Pos.Y), float32(c.Radius), clr, true)
}
func FillCircleAt(dst *ebiten.Image, origin v.Vec, radius float64, clr color.Color) {
	vector.FillCircle(dst, float32(origin.X), float32(origin.Y), float32(radius), clr, true)
}

func StrokeBox(dst *ebiten.Image, box *coll.AABB, clr color.Color) {
	vector.StrokeRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), 1, clr, false)
}

func FillOBB(cam *kamera.Camera, screen *ebiten.Image, obox *coll.OBB, c color.Color) {
	dio.GeoM.Reset()
	dio.GeoM.Scale(obox.Half.X*2, obox.Half.Y*2)
	dio.GeoM.Translate(-obox.Half.X, -obox.Half.Y)
	dio.GeoM.Rotate(obox.Angle)
	dio.GeoM.Translate(obox.Pos.X, obox.Pos.Y)
}

func StrokeBoxAt(dst *ebiten.Image, pos, half v.Vec, clr color.Color) {
	vector.StrokeRect(
		dst,
		float32(pos.X-half.X),
		float32(pos.Y-half.Y),
		float32(half.X*2),
		float32(half.Y*2),
		1,
		clr,
		false,
	)
}
func FillBox(dst *ebiten.Image, box *coll.AABB, clr color.Color) {
	vector.FillRect(dst, float32(box.Left()), float32(box.Top()), float32(box.Half.X*2), float32(box.Half.Y*2), clr, false)
}

func CursorPos() v.Vec {
	curX, curY := ebiten.CursorPosition()
	return v.Vec{float64(curX), float64(curY)}
}

func PrintHitInfoAt(dst *ebiten.Image, hit *coll.Hit, x, y int, isPenetration bool) {
	if isPenetration {
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Normal: %v\nPenetration depth: %v", hit.Normal, hit.Data), x, y)
	} else {
		ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Normal: %v\nTime: %v", hit.Normal, hit.Data), x, y)
	}
}

func Axis() (axis v.Vec) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		axis.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		axis.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		axis.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		axis.X += 1
	}
	return
}

// func DrawRay(s *ebiten.Image, pos, dir v.Vec, length float64, clr color.Color, arrow bool) {
// 	end := pos.Add(dir.Unit().Scale(length))
// 	DrawLine(s, pos, end, clr)

// 	if arrow {
// 		arrowLen := 12.0
// 		arrowAngle := math.Pi / 7

// 		unitDir := dir.Unit()

// 		left := unitDir.Rotate(math.Pi - arrowAngle).Scale(arrowLen)
// 		right := unitDir.Rotate(-(math.Pi - arrowAngle)).Scale(arrowLen)

// 		DrawLine(s, end, end.Add(left), clr)
// 		DrawLine(s, end, end.Add(right), clr)
// 	}
// }

func DrawFloor(cam *kamera.Camera, y float64, clr color.Color, s *ebiten.Image) {
	cx, cy := cam.ApplyCameraTransformToPoint(-4000, y)
	vector.StrokeLine(s, float32(cx), float32(cy), float32(4000), float32(cy), 1, clr, true)
}
func DrawSegment(s *ebiten.Image, seg *coll.Segment, clr color.Color) {
	vector.StrokeLine(s, float32(seg.A.X), float32(seg.A.Y), float32(seg.B.X), float32(seg.B.Y), 1.5, clr, true)
}

func StrokeAABB(cam *kamera.Camera, box *coll.AABB, clr color.Color, dst *ebiten.Image) {
	x, y := cam.ApplyCameraTransformToPoint(box.Left(), box.Top())
	vector.StrokeRect(
		dst,
		float32(x),
		float32(y),
		float32(box.Width()),
		float32(box.Height()),
		1,
		clr,
		true,
	)
}
