module github.com/setanarut/kamera/v2

go 1.24.5

require github.com/hajimehoshi/ebiten/v2 v2.8.8

require (
	github.com/ebitengine/gomobile v0.0.0-20250329061421-6d0a8e981e4c // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/setanarut/fastnoise v1.1.1
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
)

require github.com/setanarut/tilecollider v1.5.1

require golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect

retract (
	v2.93.0 // GeoM.Reset() bug.
	v2.6.0 // nil pointer dereference (jpg) in demo code.
	[v2.0.0, v2.2.0]
)
