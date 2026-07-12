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

// WorldToChunk converts a world coordinate.
func WorldToChunk(pos Vector3[int32]) (chunk Vector3[int32]) {
	chunk.X = pos.X >> CHUNK_SIZE_EXP
	chunk.Y = pos.Y >> CHUNK_SIZE_EXP
	chunk.Z = pos.Z >> CHUNK_SIZE_EXP
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
