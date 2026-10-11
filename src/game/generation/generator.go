package generation

import (
	"github.com/StoneTrench/go-mat-lib/vec"
	"github.com/StoneTrench/go-mc-clone/src/engine/core/async"
	"github.com/StoneTrench/go-mc-clone/src/engine/resources"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
)

type Program struct {
	Phases        []Phase
	BlockRegistry *resources.Registry[world.BlockType, world.BlockId]
}
type IPhaseResult interface {
	SetBlock(local_pos vec.Vector3I, id world.BlockId) error
	GetBlock(local_pos vec.Vector3I) (world.BlockId, error)
}

type Phase func(i int, r PhaseRuntime) world.Chunk

type PhaseRuntime struct {
	ChunkPos    vec.Vector3[int32]
	ChunkOffset vec.Vector3[int32]

	Results []world.Chunk
	Program *Program
}

func (r *PhaseRuntime) LookupId(ident string) world.BlockId {
	return r.Program.BlockRegistry.GetIdPanic(ident)
}

func Execute(prog *Program, chunk_pos vec.Vector3[int32], last_phase int) *async.Future[world.Chunk] {
	return async.NewFn(func() (world.Chunk, error) {
		r := PhaseRuntime{}
		r.Program = prog
		r.ChunkPos = chunk_pos
		r.ChunkOffset = world.ChunkToWorld(chunk_pos)
		r.Results = make([]world.Chunk, len(prog.Phases))

		n := max(0, min(len(prog.Phases), last_phase+1))

		for i := range n {
			phase := prog.Phases[i]
			r.Results[i] = phase(i, r)
		}

		return r.Results[n-1], nil
	})
}
