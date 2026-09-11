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
}

func main() {}
