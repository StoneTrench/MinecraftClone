package modding

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const API_VERSION = 0

func load_api(ctx context.Context, r wazero.Runtime, mod *Mod) error {
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr uint32, length uint32) {
			mem := m.Memory()
			buf, ok := mem.Read(ptr, length)
			if !ok {
				log_api_error("env", "log", fmt.Errorf("failed to read string from guest memory"))
				return
			}
			slog.Info(fmt.Sprintf("(Mod) (%s) %s", mod.Header.Name, string(buf)))
		}).
		Export("log").
		NewFunctionBuilder().
		WithFunc(func() uint64 {
			return uint64(API_VERSION)
		}).
		Export("api_version").
		Instantiate(ctx)

	if err != nil {
		log_api_error("host", "instantiate", err)
		return fmt.Errorf("failed to build host API: %w", err)
	}

	return nil
}

func log_api_error(moduleName, name string, err error) {
	slog.Error(fmt.Errorf("inside %s.%s and error occured, %w", moduleName, name, err).Error())
}
