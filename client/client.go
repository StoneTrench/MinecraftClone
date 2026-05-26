// Package client is the client.
package client

import (
	fw "github.com/StoneTrench/go-mc-clone/client/rendering"
	// util "github.com/StoneTrench/go-mc-clone/shared"
)

func Init() error {
	var err error = nil
	err = fw.Init()
	if err != nil {
		return err
	}

	for fw.IsRunning() {
		fw.PollEvents()

		err = fw.DrawFrames()
		if err != nil {
			return err
		}
	}

	fw.Deinit()
	return nil
}
