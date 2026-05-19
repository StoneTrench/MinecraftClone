package main

import (
	"github.com/StoneTrench/go-mc-clone/client"
	. "github.com/StoneTrench/go-mc-clone/shared"
)

func main() {
	var err error = nil
	err = InitLogging("logs", "latest")
	if err != nil {
		panic(err)
	}

	err = client.Init()
	if err != nil {
		LPanic(err)
	}
}
