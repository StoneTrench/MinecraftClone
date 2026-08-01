package modding

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/StoneTrench/go-mc-clone/game/log"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var API_VERSION = semver.MustParse("0.0.0")

func wrap_api(builder wazero.HostModuleBuilder, moduleName, funcName string, fn func(ctx context.Context, m api.Module, args ...uint32) error) {
	builder.NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, args ...uint32) uint32 {
			err := fn(ctx, m, args...)
			if err != nil {
				log_api_error(moduleName, funcName, err)
				return 1
			}
			return 0
		}).Export(funcName)
}

func log_api_error(moduleName, name string, err error) {
	log.Error(fmt.Errorf("inside %s.%s and error occured, %w", moduleName, name, err))
}

func load_api_log(ctx context.Context, r wazero.Runtime) error {
	moduleName := "log"
	mod := r.NewHostModuleBuilder(moduleName)

	// wrap_api(
	// 	mod,
	// 	moduleName,
	// 	"info",
	// 	func(ctx context.Context, m api.Module, args ...uint32) error {
	// 		view, ok := m.Memory().Read(args[0], args[1])
	// 		if !ok {
	// 			return fmt.Errorf("failed to read memory")
	// 		}

	// 		log.Info(fmt.Sprintf("[Guest Module] %s", string(view)))

	// 		return nil
	// 	},
	// )
	// wrap_api(
	// 	mod,
	// 	moduleName,
	// 	"warning",
	// 	func(ctx context.Context, m api.Module, args ...uint32) error {
	// 		view, ok := m.Memory().Read(args[0], args[1])
	// 		if !ok {
	// 			return fmt.Errorf("failed to read memory")
	// 		}

	// 		log.Warn(fmt.Sprintf("[Guest Module] %s", string(view)))

	// 		return nil
	// 	},
	// )
	// wrap_api(
	// 	mod,
	// 	moduleName,
	// 	"error",
	// 	func(ctx context.Context, m api.Module, args ...uint32) error {
	// 		view, ok := m.Memory().Read(args[0], args[1])
	// 		if !ok {
	// 			return fmt.Errorf("failed to read memory")
	// 		}

	// 		log.Error(fmt.Sprintf("[Guest Module] %s", string(view)))

	// 		return nil
	// 	},
	// )

	_, err := mod.Instantiate(ctx) // Close the runtime, module closes too
	if err != nil {
		return fmt.Errorf("failed to initialize host module (%s), %w", moduleName, err)
	}
	return nil
}

func load_api_register(ctx context.Context, r wazero.Runtime) error {
	moduleName := "register"
	mod := r.NewHostModuleBuilder(moduleName)

	// wrap_api(
	// 	mod,
	// 	moduleName,
	// 	"block",
	// 	func(ctx context.Context, m api.Module, args ...uint32) error {
	// 		view, ok := m.Memory().Read(args[0], args[1])
	// 		if !ok {
	// 			return fmt.Errorf("failed to read memory")
	// 		}

	// 		ident, err := resources.RegistryBlocks.Register(m.Name(), string(view), types.BlockType{})
	// 		if err != nil {
	// 			return err
	// 		}

	// 		log.Info(fmt.Sprintf("Registered block %s", ident))

	// 		return nil
	// 	},
	// )

	_, err := mod.Instantiate(ctx) // Close the runtime, module closes too
	if err != nil {
		return fmt.Errorf("failed to initialize host module (%s), %w", moduleName, err)
	}
	return nil
}

// func load_api_world(ctx context.Context, r wazero.Runtime) error {
// 	moduleName := "world"
// 	mod := r.NewHostModuleBuilder(moduleName)

// 	wrap_api(
// 		mod,
// 		moduleName,
// 		"set_block",
// 		func(ctx context.Context, m api.Module, args ...uint32) error {
// 			x := args[0]
// 			y := args[1]
// 			z := args[2]
// 			view, ok := m.Memory().Read(args[3], args[4])

// 			if !ok {
// 				return fmt.Errorf("failed to read memory")
// 			}

// 			log.Error(fmt.Sprintf("[Guest Module] %s", string(view)))

// 			return nil
// 		},
// 	)

// 	_, err := mod.Instantiate(ctx) // Close the runtime, module closes too
// 	if err != nil {
// 		return fmt.Errorf("failed to initialize host module (%s), %w", moduleName, err)
// 	}
// 	return nil
// }
