[![GoDoc](https://godoc.org/github.com/setanarut/kamera?status.svg)](https://pkg.go.dev/github.com/setanarut/kamera/v2)

# Kamera

Camera package for Ebitengine.

```Go
import "github.com/setanarut/kamera/v2"
```

## Features

- Camera shake effect with [fastnoise](https://github.com/setanarut/fastnoise) library noise types.
- Smooth camera movement with three interpolation modes:
  - `None`: Direct camera movement without smoothing
  - `Lerp`: Linear interpolation for smooth transitions
  - `SmoothDamp`: Spring-like motion with acceleration and deceleration and maximum speed.
- Rotate/Zoom

## Usage

A pseudo code:

```Go
func (g *Game) Update() error {
  MainCamera.LookAt(playerCenter.X, playerCenter.Y)
  return nil
 }

func (g *Game) Draw(screen *ebiten.Image) {
  playerDrawImageOptions.GeoM.Reset() // GeoM must be reset
  playerDrawImageOptions.GeoM.Translate(playerTopLeft.X, playerTopLeft.Y) // Move player
  MainCamera.Draw(playerImage, playerDrawImageOptions, screen)
 }
```

## Examples

### Platformer

Run platformer example on your local machine

```console
go run github.com/setanarut/kamera/v2/examples/platformer@latest
```

### Director

Run director example on your local machine

```console
go run github.com/setanarut/kamera/v2/examples/director@latest
```

### Director (colorm package)

Run director_colorm example on your local machine

```console
go run github.com/setanarut/kamera/v2/examples/director_colorm@latest
```