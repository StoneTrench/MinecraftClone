package client

import (
	"fmt"

	lib "github.com/StoneTrench/go-mat-lib/vec"
	world "github.com/StoneTrench/go-mc-clone/game/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func GenerateChunkMesh(chunk world.Chunk) *rl.Mesh {
	mesh := rl.Mesh{}

	solid := [32][32]uint64{}

	for x := int32(0); x < world.CHUNK_SIZE; x++ {
		for y := int32(0); y < world.CHUNK_SIZE; y++ {
			for z := int32(0); z < world.CHUNK_SIZE; z++ {
				local_pos := lib.InitVector3(x, y, z)

				block := chunk.Blocks[world.LocalToIndex(local_pos)]

				if block != 0 {
					solid[x][y] |= 1 << z
				}
			}
		}
	}

	fmt.Printf("solid: %v\n", solid)

	return &mesh
}
