package game

import (
	"github.com/StoneTrench/go-mc-clone/src/engine/resources"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
)

var RegistryBlocks = resources.CreateRegistry[world.BlockType, world.BlockId]("blocks")
