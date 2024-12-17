package kamera_test

import (
	"testing"

	"github.com/setanarut/kamera/v2"
)

func TestLookAt(t *testing.T) {
	k := kamera.NewCamera(0, 0, 100, 100)
	k.LookAt(50, 50)

	x, y := k.TopLeft()
	if x != 0 && y != 0 {
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

func TestSetSize(t *testing.T) {
	cam := kamera.NewCamera(0, 0, 10, 10)
	cam.SetSize(20, 20)
	tx, ty := cam.TopLeft()
	if tx != -10 || ty != -10 || cam.Width() != 20 || cam.Height() != 20 {
		t.Error()
	}

}
