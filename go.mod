module github.com/setanarut/kamera/v2

go 1.23.4

require github.com/hajimehoshi/ebiten/v2 v2.8.6

require (
	github.com/ebitengine/gomobile v0.0.0-20241016134836-cc2e38a7c0ee // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.2 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/setanarut/fastnoise v1.1.1
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

require github.com/setanarut/tilecollider v1.4.7

require golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect

retract (
	v2.93.0 // GeoM.Reset() bug.
	v2.6.0 // nil pointer dereference (jpg) in demo code.
	[v2.0.0, v2.2.0]
)
