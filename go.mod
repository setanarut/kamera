module github.com/setanarut/kamera/v2

go 1.25

require github.com/hajimehoshi/ebiten/v2 v2.9.5

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.1 // indirect
	github.com/jezek/xgb v1.2.0 // indirect
	github.com/setanarut/fastnoise v1.2.0
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)

retract (
	v2.93.0 // GeoM.Reset() bug.
	v2.6.0 // nil pointer dereference (jpg) in demo code.
	[v2.0.0, v2.2.0]
)
