package world

import (
	"iter"

	. "github.com/StoneTrench/go-mat-lib/vec"
)

const CHUNK_SIZE_X_EXP = 5
const CHUNK_SIZE_Y_EXP = 6
const CHUNK_SIZE_Z_EXP = 5

const CHUNK_SIZE_X = 1 << CHUNK_SIZE_X_EXP
const CHUNK_SIZE_Y = 1 << CHUNK_SIZE_Y_EXP
const CHUNK_SIZE_Z = 1 << CHUNK_SIZE_Z_EXP

const CHUNK_SIZE_X_MASK = CHUNK_SIZE_X - 1
const CHUNK_SIZE_Y_MASK = CHUNK_SIZE_Y - 1
const CHUNK_SIZE_Z_MASK = CHUNK_SIZE_Z - 1

const CHUNK_VOLUME = CHUNK_SIZE_X * CHUNK_SIZE_Y * CHUNK_SIZE_Z

const OBJECT_VOXEL_SIZE = 0.5
const OBJECT_CHUNK_SIZE_X = CHUNK_SIZE_X * OBJECT_VOXEL_SIZE
const OBJECT_CHUNK_SIZE_Y = CHUNK_SIZE_Y * OBJECT_VOXEL_SIZE
const OBJECT_CHUNK_SIZE_Z = CHUNK_SIZE_Z * OBJECT_VOXEL_SIZE

func IterateChunk() iter.Seq[Vector3[int32]] {
	return func(yield func(Vector3[int32]) bool) {
		for x := range int32(CHUNK_SIZE_X) {
			for y := range int32(CHUNK_SIZE_Y) {
				for z := range int32(CHUNK_SIZE_Z) {
					if !yield(InitVector3(x, y, z)) {
						return
					}
				}
			}
		}
	}
}

// WorldToChunk finds which chunk the position is in.
func WorldToChunk(pos Vector3[int32]) (chunk Vector3[int32]) {
	chunk.X = pos.X >> CHUNK_SIZE_X_EXP
	chunk.Y = pos.Y >> CHUNK_SIZE_Y_EXP
	chunk.Z = pos.Z >> CHUNK_SIZE_Z_EXP
	return chunk
}

// WorldToChunk finds the origin of the chunk in the world.
func ChunkToWorld(pos Vector3[int32]) (chunk Vector3[int32]) {
	chunk.X = pos.X << CHUNK_SIZE_X_EXP
	chunk.Y = pos.Y << CHUNK_SIZE_Y_EXP
	chunk.Z = pos.Z << CHUNK_SIZE_Z_EXP
	return chunk
}

// WorldToLocal finds the position inside the chunk.
func WorldToLocal(pos Vector3[int32]) (local Vector3[int32]) {
	local.X = pos.X & CHUNK_SIZE_X_MASK
	local.Y = pos.Y & CHUNK_SIZE_Y_MASK
	local.Z = pos.Z & CHUNK_SIZE_Z_MASK
	return local
}

// LocalToIndex converts local coordinates into a single integer to index into the block array in a chunk.
func LocalToIndex(local Vector3[int32]) (index int32) {
	return local.Z + local.X*CHUNK_SIZE_Z + local.Y*CHUNK_SIZE_X*CHUNK_SIZE_Z
}

// IndexToLocal converts index into local coordinates in a chunk.
func IndexToLocal(index int32) (local Vector3[int32]) {
	local.Z = index % CHUNK_SIZE_Z
	local.Y = index / (CHUNK_SIZE_X * CHUNK_SIZE_Z)
	local.X = (index / CHUNK_SIZE_Z) % CHUNK_SIZE_Z
	return local
}

var NEIGHBOURS = [6]Vector3[int32]{
	InitVector3[int32](0, 0, 1),  // back
	InitVector3[int32](0, 0, -1), // front
	InitVector3[int32](0, 1, 0),  // top
	InitVector3[int32](0, -1, 0), // bottom
	InitVector3[int32](-1, 0, 0), // left
	InitVector3[int32](1, 0, 0),  // right
}
