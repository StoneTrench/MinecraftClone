package main

import (
	"log/slog"

	"github.com/StoneTrench/go-mc-clone/src/client"
)

func main() {
	err := client.Init()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
