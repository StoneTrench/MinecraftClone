package main

import (
	"github.com/StoneTrench/go-mc-clone/client"
	. "github.com/StoneTrench/go-mc-clone/game"
	"github.com/Tnze/go-mc/level"
	// "github.com/StoneTrench/go-mc-clone/game/lua"
)

func main() {

	

	client.GenerateChunkMesh()

	// var err error = nil
	// err = InitLogging("logs", "latest")
	// if err != nil {
	// 	panic(err)
	// }

	// LInfo(GetFormattedApplicationLabel())

	// l := lua.State{}.Init()
	// defer l.Deinit()

	// err = l.DoString("print('Hello from lua! :3')")
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.Init()
	// if err != nil {
	// 	panic(err)
	// }
}
