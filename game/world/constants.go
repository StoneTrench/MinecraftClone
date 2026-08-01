package world

import (
	. "github.com/StoneTrench/go-mat-lib/vec"
)

const CHUNK_SIZE_EXP = 5

const CHUNK_SIZE = 1 << CHUNK_SIZE_EXP
const CHUNK_SIZE_MASK = CHUNK_SIZE - 1

const CHUNK_LAYER_AREA = CHUNK_SIZE * CHUNK_SIZE
const CHUNK_VOLUME = CHUNK_SIZE * CHUNK_SIZE * CHUNK_SIZE

const OBJECT_VOXEL_SIZE = 0.5
const OBJECT_CHUNK_SIZE = CHUNK_SIZE * OBJECT_VOXEL_SIZE

// WorldToChunk finds which chunk the position is in.
func WorldToChunk(pos Vector3[int32]) (chunk Vector3[int32]) {
	chunk.X = pos.X >> CHUNK_SIZE_EXP
	chunk.Y = pos.Y >> CHUNK_SIZE_EXP
	chunk.Z = pos.Z >> CHUNK_SIZE_EXP
	return chunk
}

// WorldToChunk finds the origin of the chunk in the world.
func ChunkToWorld(pos Vector3[int32]) (chunk Vector3[int32]) {
	chunk.X = pos.X << CHUNK_SIZE_EXP
	chunk.Y = pos.Y << CHUNK_SIZE_EXP
	chunk.Z = pos.Z << CHUNK_SIZE_EXP
	return chunk
}

// WorldToLocal finds the position inside the chunk.
func WorldToLocal(pos Vector3[int32]) (local Vector3[int32]) {
	local.X = pos.X & CHUNK_SIZE_MASK
	local.Y = pos.Y & CHUNK_SIZE_MASK
	local.Z = pos.Z & CHUNK_SIZE_MASK
	return local
}

// LocalToIndex converts local coordinates into a single integer to index into the block array in a chunk.
func LocalToIndex(local Vector3[int32]) (index int32) {
	return local.X + local.Y*CHUNK_SIZE + local.Z*CHUNK_LAYER_AREA
}

// IndexToLocal converts index into local coordinates in a chunk.
func IndexToLocal(index int32) (local Vector3[int32]) {
	local.X = index % CHUNK_SIZE
	local.Y = (index / CHUNK_SIZE) % CHUNK_SIZE
	local.Z = index / CHUNK_LAYER_AREA
	return local
}

var NEIGHBOURS = [6][3]int32{
	{0, 0, 1},  // back
	{0, 0, -1}, // front
	{0, 1, 0},  // top
	{0, -1, 0}, // bottom
	{-1, 0, 0}, // left
	{1, 0, 0},  // right
}

// func IterateNeighbours() iter.Seq[Vector3[int32]] {
// 	return func(yield func(Vector3[int32]) bool) {
// 		for i := range 6 {
// 			if !yield(InitVector3(NEIGHBOURS[i][0], NEIGHBOURS[i][1], NEIGHBOURS[i][2])) {
// 				return
// 			}
// 		}
// 	}
// }
