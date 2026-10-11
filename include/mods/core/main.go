//go:build wasip1

package main

import (
	"fmt"

	eb "github.com/StoneTrench/go-mc-clone/bindings/go"
)

func RegBlock(name string, r, g, b, a uint8, sld bool) {
	namespace := eb.RegisterBlock(eb.BlockType{
		Name:    name,
		Color:   eb.Color{R: r, G: g, B: b, A: a},
		IsSolid: sld,
	})
	eb.LogInfo(fmt.Sprintf("Registered block: %s", namespace))
}

func init() {
	eb.LogInfo("Module initialized!")

	eb.AppendLang("en-us")
	RegBlock("air", 255, 0, 255, 255, false)
	RegBlock("stone", 128, 128, 128, 255, true)
	RegBlock("dirt", 255, 128, 128, 255, true)
	RegBlock("grass", 128, 255, 128, 255, true)

	RegBlock("water", 51, 102, 204, 180, false)
	RegBlock("sand", 219, 205, 135, 255, true)
	RegBlock("snow", 240, 245, 255, 255, true)
	RegBlock("iron_ore", 168, 122, 90, 255, true)
	RegBlock("log", 102, 71, 43, 255, true)
	RegBlock("leaves", 58, 122, 47, 220, true)
}

func main() {}
