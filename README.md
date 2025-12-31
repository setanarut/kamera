[![GoDoc](https://godoc.org/github.com/setanarut/kamera?status.svg)](https://pkg.go.dev/github.com/setanarut/kamera/v2)

# Kamera

Camera package for Ebitengine.

```Go
import "github.com/setanarut/kamera/v2"
```

https://github.com/user-attachments/assets/58787fbe-5788-43d4-9560-482dcc965efb


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

1. Clone this repository
2. In the terminal, change directory to the examples folder `cd examples`
3. Run a demo with `go run ./folder_name` (`go run ./platformer`, `go run ./director`).