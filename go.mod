module github.com/setanarut/kamera/v2

go 1.23.2

require github.com/hajimehoshi/ebiten/v2 v2.8.0

require github.com/setanarut/vec v1.1.1 // indirect

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/setanarut/cm v1.13.0
	github.com/setanarut/fastnoise v1.1.1
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

retract (
	v2.6.0 // nil pointer dereference (jpg) in demo code.
	[v2.0.0, v2.2.0]
)
