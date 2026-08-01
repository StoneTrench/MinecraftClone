package game

import (
	"image/color"

	"github.com/StoneTrench/go-mc-clone/game/events"
	"github.com/StoneTrench/go-mc-clone/game/log"
	"github.com/StoneTrench/go-mc-clone/game/modding"
	"github.com/StoneTrench/go-mc-clone/game/regentries"
	"github.com/StoneTrench/go-mc-clone/game/resources"
	"github.com/StoneTrench/go-mc-clone/metadata"
)

var EventBusMain = events.CreateEventBus()
var EventBusNetwork = events.CreateEventBus()

func Init() error {
	var err error = nil
	err = metadata.Init()
	if err != nil {
		return err
	}

	err = log.Init("logs", "latest")
	if err != nil {
		return err
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
	log.Info(metadata.GetFormattedApplicationLabel())

	RegistryBlocks.RegisterPanic("test", "air", regentries.BlockType{
		IsSolid: false,
		Color:   color.RGBA{255, 0, 255, 255},
	})
	RegistryBlocks.RegisterPanic("test", "stone", regentries.BlockType{
		IsSolid: true,
		Color:   color.RGBA{128, 128, 128, 255},
	})
	RegistryBlocks.RegisterPanic("test", "grass", regentries.BlockType{
		IsSolid: true,
		Color:   color.RGBA{128, 255, 128, 255},
	})

	err = modding.Init()
	if err != nil {
		log.Error(err)
	}
	RegistryBlocks.Freeze()
	RegistryBlocks.ApplyIdMap([]string{"test:air", "test:stone", "test:grass", "test:test"})

	return nil
}

func Update() error {
	EventBusMain.HandleDeferred(200)
	EventBusNetwork.HandleDeferred(200)

	return nil
}
