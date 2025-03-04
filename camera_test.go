package kamera_test

import (
	"testing"

	"github.com/setanarut/kamera/v2"
)

func TestLookAt(t *testing.T) {
	k := kamera.NewCamera(0, 0, 100, 100)
	k.LookAt(50, 50)

	if k.X != 0 && k.Y != 0 {
		t.Error()
	}
}

func TestCenter(t *testing.T) {
	k := kamera.NewCamera(0, 0, 100, 100)
	k.LookAt(2.5, 4.2)

	x, y := k.Center()
	if x != 2.5 && y != 4.2 {
		t.Error()
	}
}
func TestSetCenter(t *testing.T) {
	k := kamera.NewCamera(0, 0, 100, 100)
	k.SmoothType = kamera.SmoothDamp
	k.SetCenter(3, 5)
	x, y := k.Center()
	if x != 3 && y != 5 {
		t.Error()
	}
	if k.TempTargetX != 3 || k.TempTargetY != 5 {
		t.Error()
	}
	k.SmoothType = kamera.Lerp
	k.SetCenter(5, 3)
	x, y = k.Center()
	if x != 5 && y != 3 {
		t.Error()
	}
	if k.TempTargetX != 5 || k.TempTargetY != 3 {
		t.Error()
	}
}
