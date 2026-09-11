package engine

import (
	"fmt"

	"github.com/StoneTrench/go-mc-clone/src/engine/log"
)

func Init(logDir string) error {
	err := log.Init(logDir, false)
	if err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}

	err = init_meta()
	if err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	return nil
}
