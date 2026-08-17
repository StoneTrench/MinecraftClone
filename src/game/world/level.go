package world

import (
	"sync"

	"github.com/StoneTrench/go-mat-lib/vec"
)

// TODO: Split the level into regions, which is good for when you want to generate large structures.
// Because then chunks can query the region to see if they should place a structure in themselves.

type Level struct {
	mu     sync.RWMutex
	Chunks map[vec.Vector3I]*Chunk
}

func CreateLevel() *Level {
	return &Level{
		mu:     sync.RWMutex{},
		Chunks: make(map[vec.Vector3I]*Chunk),
	}
}

func (l *Level) GetChunk(chunk_pos vec.Vector3I, create_new_chunk bool) *Chunk {
	if chunk, exists := l.Chunks[chunk_pos]; exists {
		return chunk
	}

	if !create_new_chunk {
		return nil
	}

	chunk := CreateChunk()
	l.Chunks[chunk_pos] = chunk

	return chunk
}
func (l *Level) SetChunk(chunk_pos vec.Vector3I, chunk *Chunk) {
	chunk_copy := *chunk
	l.Chunks[chunk_pos] = &chunk_copy
}

func (l *Level) SetBlock(level_pos vec.Vector3I, id BlockId, create_new_chunk bool) error {
	return l.GetChunk(WorldToChunk(level_pos), create_new_chunk).SetBlock(WorldToLocal(level_pos), id)
}
func (l *Level) GetBlock(level_pos vec.Vector3I, create_new_chunk bool) (BlockId, error) {
	return l.GetChunk(WorldToChunk(level_pos), create_new_chunk).GetBlock(WorldToLocal(level_pos))
}
