module github.com/StoneTrench/go-mc-clone

go 1.26.2

require (
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/StoneTrench/go-mat-lib v1.0.0
	github.com/tetratelabs/wazero v1.12.0
	golang.org/x/exp v0.0.0-20260908205506-85c1c2202aba
)

require (
	github.com/gen2brain/raylib-go/raylib v0.60.1
	golang.org/x/sys v0.48.0 // indirect
)

require (
	github.com/ebitengine/purego v0.10.0 // indirect
	github.com/jupiterrider/ffi v0.7.0 // indirect
	github.com/ojrac/opensimplex-go v1.0.2
)

replace (
	// github.com/BrownNPC/Raylib-Go-Wasm/wasm-runtime => ./Raylib-Go-Wasm/wasm-runtime
	github.com/StoneTrench/go-mat-lib => ../GoMatLib
// github.com/gen2brain/raylib-go/raygui => ./Raylib-Go-Wasm/raygui
// github.com/gen2brain/raylib-go/raylib => ./Raylib-Go-Wasm/raylib
)
