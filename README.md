[![GoDoc](https://godoc.org/github.com/setanarut/kamera?status.svg)](https://pkg.go.dev/github.com/setanarut/kamera/v2)

# Kamera

Camera package for Ebitengine.

```Go
import "github.com/setanarut/kamera/v2"
```

## Features

- Camera shake effect with [fastnoise](https://github.com/setanarut/fastnoise) library noise types.
- Lerp motion
- Rotate/Zoom

## Example

Run example on your local machine

```console
go run github.com/setanarut/kamera/v2/examples/demo@latest
```

## Usage

A pseudo code:

```Go
func (g *Game) Update() error {
  g.MainCamera.LookAt(player.X, player.Y)
  // Apply all world-space `playerDrawImageOptions.GeoM{}` transform here
 }
func (g *Game) Draw(screen *ebiten.Image) {
  g.MainCamera.Draw(playerImage, playerDrawImageOptions, screen)
 }
```

See the [./examples/demo/main.go](./examples/demo/main.go) for a real working code example. 
