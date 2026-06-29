// Package client is the client.
package client

import (
	"time"

	fw "github.com/StoneTrench/go-mc-clone/client/rendering"
	// util "github.com/StoneTrench/go-mc-clone/shared"
)

func Init() error {
	var err error = nil
	err = fw.Init()
	if err != nil {
		return err
	}

	var time_current float64 = 0
	lastTime := time.Now()

	for fw.IsRunning() {
		fw.PollEvents()

		now := time.Now()
		time_delta := now.Sub(lastTime).Seconds()
		lastTime = now

		err = fw.DrawFrames(time_current, time_delta)
		if err != nil {
			return err
		}

		time_current += time_delta
	}

	fw.Deinit()
	return nil
}
