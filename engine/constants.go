package engine

import (
	. "github.com/StoneTrench/go-mat-lib"
)

const CHUNK_SIZE_EXP = 5

const CHUNK_SIZE = 1 << CHUNK_SIZE_EXP
const CHUNK_SIZE_MASK = CHUNK_SIZE - 1

// WorldToChunk converts a world coordinate
func WorldToChunk[T Integer](pos Vector3[T]) Vector3[T] {
	x := pos.X >> CHUNK_SIZE_EXP
	y := pos.Y >> CHUNK_SIZE_EXP
	z := pos.Z >> CHUNK_SIZE_EXP
	return Vector3[T]{X: x, Y: y, Z: z}
}

// WorldToLocal finds the position inside the chunk
func WorldToLocal[T Integer](pos Vector3[T]) Vector3[T] {
	x := pos.X & CHUNK_SIZE_MASK
	y := pos.Y & CHUNK_SIZE_MASK
	z := pos.Z & CHUNK_SIZE_MASK
	return Vector3[T]{X: x, Y: y, Z: z}
}
