module github.com/StoneTrench/go-mc-clone

go 1.26.2

require github.com/StoneTrench/go-mat-lib v1.0.0

require (
	github.com/ebitengine/purego v0.10.0 // indirect
	github.com/jupiterrider/ffi v0.7.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
)

require (
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/gen2brain/raylib-go/raylib v0.60.0
	github.com/tetratelabs/wazero v1.12.0
	golang.org/x/sys v0.44.0 // indirect
)

replace github.com/StoneTrench/go-mat-lib => ../GoMatLib
