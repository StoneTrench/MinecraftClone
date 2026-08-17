// Package client is the client.
package client

import (
	"fmt"

	"github.com/StoneTrench/go-mat-lib/other"
	"github.com/StoneTrench/go-mat-lib/vec"
	"github.com/StoneTrench/go-mc-clone/src/engine"
	"github.com/StoneTrench/go-mc-clone/src/engine/log"
	"github.com/StoneTrench/go-mc-clone/src/game"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func raylib_trace_log_callback(level int, msg string) {
	msg = fmt.Sprintf("(RayLib) %s", msg)

	switch rl.TraceLogLevel(level) {
	case rl.LogInfo:
		log.Info(msg)
	case rl.LogWarning:
		log.WarnStr(msg)
	case rl.LogError:
		log.ErrorStr(msg)
	case rl.LogFatal:
		log.ErrorStr(msg)
	default:
		log.Info(fmt.Sprintf("(Unknown TraceLogLevel) %s", msg))
	}
}

func Init() error {
	err := game.Init()
	if err != nil {
		panic(err)
	}

	level := world.CreateLevel()
	chunk_pos := vec.InitVector3[int32](0, 0, 0)
	chunk := level.GetChunk(chunk_pos, true)

	for local_pos := range world.IterateChunk() {
		fx := float32(local_pos.X)
		fz := float32(local_pos.Z)
		y := local_pos.Y
		sinX, _ := other.Sincos(fx / 10.0)
		_, cosZ := other.Sincos(fz / 10.0)
		h := int32((sinX*cosZ)*12 + 15)

		if y < h {
			chunk.SetBlock(local_pos, 1)
		} else if y < h+1 {
			chunk.SetBlock(local_pos, 2)
		} else {
			chunk.SetBlock(local_pos, 0)
		}
	}

	rl.SetTraceLogCallback(raylib_trace_log_callback)
	rl.InitWindow(800, 600, engine.GetFormattedApplicationLabel())
	defer rl.CloseWindow()

	cam := rl.Camera{
		Position:   rl.NewVector3(0, 0, 0),
		Target:     rl.NewVector3(0, 0, -1),
		Up:         rl.NewVector3(0, 1, 0),
		Projection: rl.CameraPerspective,
		Fovy:       45,
	}

	msh, err := GenerateChunkMesh(level, chunk_pos, game.RegistryBlocks)
	if err != nil {
		return fmt.Errorf("failed to generate chunk mesh, %w", err)
	}
	defer rl.UnloadMesh(&msh)
	mat := rl.LoadMaterialDefault()
	defer rl.UnloadMaterial(mat)

	rl.DisableCursor()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyDown(rl.KeyF2) {
			rl.TakeScreenshot("screenshot.png")
		}

		rl.UpdateCamera(&cam, rl.CameraFree)

		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.RayWhite)

			rl.BeginMode3D(cam)
			{
				rl.DrawMesh(msh, mat, rl.MatrixIdentity())
				rl.DrawGrid(100, 1)
			}
			rl.EndMode3D()
		}
		rl.EndDrawing()
	}

	return nil
}
