// Package client is the client.
package client

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	game "github.com/StoneTrench/go-mc-clone/game"
)

func Init() error {
	rl.InitWindow(800, 600, game.GetFormattedApplicationLabel())

	for !rl.WindowShouldClose() {

	}

	rl.CloseWindow()
	return nil
}
