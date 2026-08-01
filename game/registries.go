package game

import (
	"github.com/StoneTrench/go-mc-clone/game/regentries"
	"github.com/StoneTrench/go-mc-clone/game/resources"
)

var RegistryBlocks = resources.CreateRegistry[regentries.BlockType, regentries.BlockId](0, "blocks")
