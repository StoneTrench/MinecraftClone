package main

import (
	"github.com/StoneTrench/go-mc-clone/client"
	"github.com/StoneTrench/go-mc-clone/game/log"
)

func main() {
	err := client.Init()
	if err != nil {
		log.Panic(err)
	}
}
