package game

import (
	"fmt"
	"image/color"

	"github.com/StoneTrench/go-mc-clone/src/engine"
	"github.com/StoneTrench/go-mc-clone/src/engine/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/resources"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
	"github.com/StoneTrench/go-mc-clone/src/modding"
)

func Init() error {
	var err error = nil
	err = log.Init(false)
	if err != nil {
		return fmt.Errorf("failed to initialize logging, %w", err)
	}

	err = engine.Init()
	if err != nil {
		log.Errorf("failed to initialize metadata, %w", err)
	}

	err = resources.LangAppendFallback("assets/lang/en-us.json")
	if err != nil {
		log.Error(err)
	}
	err = resources.LangAppend("assets/lang/en-us.json")
	if err != nil {
		log.Error(err)
	}
	log.Info(resources.LangTranslate("test.message", "en-us"))
	log.Info(engine.GetFormattedApplicationLabel())

	RegistryBlocks.RegisterPanic("test", "air", world.BlockType{
		IsSolid: false,
		Color:   color.RGBA{255, 0, 255, 255},
	})
	RegistryBlocks.RegisterPanic("test", "stone", world.BlockType{
		IsSolid: true,
		Color:   color.RGBA{128, 128, 128, 255},
	})
	RegistryBlocks.RegisterPanic("test", "grass", world.BlockType{
		IsSolid: true,
		Color:   color.RGBA{128, 255, 128, 255},
	})

	err = modding.Init()
	if err != nil {
		log.Error(err)
	}
	RegistryBlocks.Freeze()
	RegistryBlocks.ApplyIdMap([]string{"test:air", "test:stone", "test:grass"})

	return nil
}

func Update() error {

	return nil
}
