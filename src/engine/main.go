package engine

import (
	"fmt"

	"github.com/StoneTrench/go-mc-clone/src/engine/core/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/platform"
)

func Init(logDir string) error {
	var err error = nil

	err = platform.InitMeta()
	if err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	err = log.Init(logDir, false)
	if err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}

	return nil
}
