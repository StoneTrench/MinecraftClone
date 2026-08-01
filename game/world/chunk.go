package world

import (
	"fmt"

	"github.com/StoneTrench/go-mat-lib/vec"
	"github.com/StoneTrench/go-mc-clone/game/ecs"
)

type Chunk struct {
	// The intention with this flag, is to allow setting blocks to chunks that haven't been generated yet,
	// and when the chunk gets generated, these values would get overlaid on it, preserving them, as if the
	// chunk was always there.
	IsGenerated bool
	Blocks      [CHUNK_VOLUME]BlockId
	Entities    map[ecs.EntityId]struct{}
}

func CreateChunk() *Chunk {
	return &Chunk{
		IsGenerated: false,
		Blocks:      [CHUNK_VOLUME]BlockId{},
		Entities:    make(map[ecs.EntityId]struct{}),
	}
}

func (c *Chunk) SetBlock(local_pos vec.Vector3I, id BlockId) error {
	if c == nil {
		return fmt.Errorf("nil chunk")
	}

	c.Blocks[LocalToIndex(local_pos)] = id
	return nil
}
func (c *Chunk) GetBlock(local_pos vec.Vector3I) (BlockId, error) {
	if c == nil {
		return 0, fmt.Errorf("nil chunk")
	}

	return c.Blocks[LocalToIndex(local_pos)], nil
}

func (c *Chunk) AddEntity(id ecs.EntityId) error {
	if c == nil {
		return fmt.Errorf("nil chunk")
	}
	if _, ok := c.Entities[id]; ok {
		return fmt.Errorf("entity %v already in chunk", id)
	}

	c.Entities[id] = struct{}{}
	return nil
}
func (c *Chunk) RemoveEntity(id ecs.EntityId) error {
	if c == nil {
		return fmt.Errorf("nil chunk")
	}

	if _, ok := c.Entities[id]; !ok {
		return fmt.Errorf("entity %v not in chunk", id)
	}
	delete(c.Entities, id)
	return nil
}
