package modding

import (
	"fmt"

	"github.com/StoneTrench/go-mc-clone/bindings/host"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

var host_module_parts []host.HostBinding

// Registers a new Api function to be loaded into the mods.
func RegisterApiFunction(b host.HostBinding) {
	host_module_parts = append(host_module_parts, b)
}

func load_api(r wazero.Runtime, mod *Mod) error {
	host := r.NewHostModuleBuilder("engine_bindings")

	_, err := wasi_snapshot_preview1.Instantiate(mod.ctx, r)
	if err != nil {
		return fmt.Errorf("failed to instantiate wasip1 for (%s): %w", mod.String(), err)
	}

	for _, p := range host_module_parts {
		host.NewFunctionBuilder().WithFunc(p.Fn).Export(p.Name)
	}

	comp_host, err := host.Compile(mod.ctx)
	if err != nil {
		return fmt.Errorf("failed to compile host module for (%s): %w", mod.String(), err)
	}
	_, err = r.InstantiateModule(mod.ctx, comp_host, wazero.NewModuleConfig())
	if err != nil {
		return fmt.Errorf("failed to instantiate host module for (%s): %w", mod.String(), err)
	}

	return nil
}
