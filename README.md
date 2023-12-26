# kamera

[kamera](https://pkg.go.dev/github.com/setanarut/kamera) is a camera package for Ebitengine.

## Example

Run example on your local machine

```console
go run github.com/setanarut/kamera/examples/topdown@latest
```

| Key   | Action                     |
| ----- | -------------------------- |
| WASD  | Move camera                |
| E     | Zoom in                    |
| Q     | Zoom out                   |
| Space | Reset camera               |
| Key 1 | look at the random object. |
| Key 2 | Reset Zoom                 |


## Usage

Apply all world-space `DrawImageOptions.GeoM` transformations before `Camera.Render()`. Pay attention to the order, the order will affect the camera drawing. 

A pseudo code:

```Go
func (g *Game) Update() error {
    g.MainCamera.LookAt(player.X, player.Y)
 }
func (g *Game) Draw(screen *ebiten.Image) {
    //  <- playerDrawImageOptions.GeoM world-space transformations will come here, or in Update()
    g.MainCamera.Render(player, playerDrawImageOptions, screen)
 }
```

See the [./examples/topdown/main.go](./examples/topdown/main.go) for a real working code example. 