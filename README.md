# kamera

[kamera](https://pkg.go.dev/github.com/setanarut/kamera) is a camera package for Ebitengine.

## Features

- Camera trauma shake effect with Opensimplex noise
- Lerp motion
- Rotate/Zoom

## Example

Run example on your local machine

```console
go run github.com/setanarut/kamera/examples/demo@latest
```

## Usage

A pseudo code:

```Go
func (g *Game) Update() error {
    g.MainCamera.LookAt(player.X, player.Y)
 }
func (g *Game) Draw(screen *ebiten.Image) {
    // Apply all world-space `DrawImageOptions.GeoM` transformations before `Camera.Draw()`
    g.MainCamera.Draw(player, playerDrawImageOptions, screen)
 }
```

See the [./examples/demo/main.go](./examples/demo/main.go) for a real working code example. 