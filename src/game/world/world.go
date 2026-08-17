package world

import "github.com/StoneTrench/go-mc-clone/src/engine/ecs"

// World is a datastructure containing a game world.
// Static level is the world/environment.
// And then there's an ecs, where ships are basically levels with a transform component and collision.
type World struct {
	NamespaceIDMap map[string]uint16
	StaticLevel    Level
	Entities       ecs.EntityWorld
}
