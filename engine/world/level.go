package world

import (
	"errors"

	. "github.com/StoneTrench/go-mc-clone/engine"

	. "github.com/StoneTrench/go-mat-lib"
)

type BlockId uint16
type Level struct {
	Chunks map[Vector3I]*Chunk
}
type Chunk struct {
	IsGenerated bool
	Nbt         map[Vector3[uint8]]any
	Blocks      [CHUNK_SIZE][CHUNK_SIZE][CHUNK_SIZE]BlockId
}

func GetChunk(level Level, chunk_pos Vector3I, create_new_chunk bool) *Chunk {
	if chunk, exists := level.Chunks[chunk_pos]; exists {
		return chunk
	}

	if !create_new_chunk {
		return nil
	}

	chunk := &Chunk{
		IsGenerated: false,
		Nbt:         make(map[Vector3[uint8]]any),
		Blocks:      [CHUNK_SIZE][CHUNK_SIZE][CHUNK_SIZE]BlockId{},
	}
	level.Chunks[chunk_pos] = chunk

	return chunk
}
func SetChunk(level Level, chunk_pos Vector3I, chunk *Chunk) {
	level.Chunks[chunk_pos] = chunk
}

func SetBlock(level Level, level_pos Vector3I, id BlockId, create_new_chunk bool) error {
	chunk_pos := WorldToChunk(level_pos)
	local_pos := WorldToLocal(level_pos)

	chunk := GetChunk(level, chunk_pos, create_new_chunk)
	if chunk == nil {
		return errors.New("no chunk, outside of world")
	}

	chunk.Blocks[local_pos.X][local_pos.Y][local_pos.Z] = id
	return nil
}
func GetBlock(level Level, level_pos Vector3I) (BlockId, error) {
	chunk_pos := WorldToChunk(level_pos)
	local_pos := WorldToLocal(level_pos)

	chunk := GetChunk(level, chunk_pos, false)
	if chunk == nil {
		return 0, errors.New("no chunk, outside of world")
	}

	return chunk.Blocks[local_pos.X][local_pos.Y][local_pos.Z], nil
}
