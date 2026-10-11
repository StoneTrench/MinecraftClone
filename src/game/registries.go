package game

import (
	"github.com/StoneTrench/go-mc-clone/src/engine/resources"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
)

var RegistryBlocks = new(resources.Registry[world.BlockType, world.BlockId]).Init("blocks")
