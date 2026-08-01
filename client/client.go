// Package client is the client.
package client

import (
	"fmt"

	"github.com/StoneTrench/go-mat-lib/other"
	"github.com/StoneTrench/go-mat-lib/vec"
	"github.com/StoneTrench/go-mc-clone/game"
	"github.com/StoneTrench/go-mc-clone/game/log"
	"github.com/StoneTrench/go-mc-clone/game/world"
	"github.com/StoneTrench/go-mc-clone/metadata"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func raylib_trace_log_callback(level int, msg string) {
	msg = fmt.Sprintf("(RayLib) %s", msg)

	switch rl.TraceLogLevel(level) {
	case rl.LogInfo:
		log.Info(msg)
	case rl.LogWarning:
		log.Warn(msg)
	case rl.LogError:
		log.Error(msg)
	case rl.LogFatal:
		log.Panic(msg)
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

	for x := range int32(world.CHUNK_SIZE) {
		for y := range int32(world.CHUNK_SIZE) {
			for z := range int32(world.CHUNK_SIZE) {
				local_pos := vec.InitVector3(x, y, z)

				fx := float32(x)
				fz := float32(z)
				sinX, _ := other.Sincos(fx / 10.0)
				_, cosZ := other.Sincos(fz / 10.0)
				h := int32((sinX*cosZ)*12 + 10)

				if y < h {
					chunk.SetBlock(local_pos, 1)
				} else if y < h+1 {
					chunk.SetBlock(local_pos, 2)
				} else {
					chunk.SetBlock(local_pos, 0)
				}
			}
		}
	}

	rl.SetTraceLogCallback(raylib_trace_log_callback)
	rl.InitWindow(800, 600, metadata.GetFormattedApplicationLabel())
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
