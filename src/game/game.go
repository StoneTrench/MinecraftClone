package game

import (
	"fmt"
	"image/color"
	"log/slog"
	"path"

	"github.com/StoneTrench/go-mc-clone/bindings/host"
	"github.com/StoneTrench/go-mc-clone/src/engine"
	"github.com/StoneTrench/go-mc-clone/src/engine/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/modding"
	"github.com/StoneTrench/go-mc-clone/src/engine/resources"
	"github.com/StoneTrench/go-mc-clone/src/game/world"
	"github.com/tetratelabs/wazero/api"
)

func Init() error {
	var err error = nil
	err = engine.Init("./logs/")
	if err != nil {
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	log.Info(engine.GetFormattedApplicationLabel())
	load_api()

	err = modding.LoadMods("./mods/")
	if err != nil {
		log.Error(err)
	}
	RegistryBlocks.Freeze()
	RegistryBlocks.ApplyIdMap([]string{"core:air", "core:stone", "core:dirt", "core:grass"})

	return nil
}

func load_api() {
	modding.RegisterApiFunction(host.LogInfo(func(m api.Module, msg string) {
		slog.Info(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))
	modding.RegisterApiFunction(host.LogWarn(func(m api.Module, msg string) {
		slog.Warn(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))
	modding.RegisterApiFunction(host.LogError(func(m api.Module, msg string) {
		slog.Error(fmt.Sprintf("(Mod) (%s) %s", m.Name(), msg))
	}))

	modding.RegisterApiFunction(host.RegisterBlock(func(m api.Module, block host.BlockType) (namespaceid string) {
		return RegistryBlocks.RegisterPanic(m.Name(), block.Name, world.BlockType{
			Color:   color.RGBA(block.Color),
			IsSolid: block.IsSolid,
		})
	}))
	modding.RegisterApiFunction(host.AppendLang(func(m api.Module, name string) {
		mod := modding.GetModWasm(m)
		if mod == nil {
			panic("cannot find mod")
		}
		if len(name) != 5 {
			panic("name must be 5 characters long")
		}
		err := resources.LangAppend(path.Join(mod.Folder, "assets", "lang", name+".json"))
		if err != nil {
			panic(err)
		}
	}))

	// modding.RegisterApiFunction(host.ConfigSetString(func(m api.Module, name, value string) {

	// }))
}
