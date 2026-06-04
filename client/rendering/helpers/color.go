package helpers

import (
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

func NewClearColorFloat(r, g, b, a float32) vk.ClearColorValue {
	var c vk.ClearColorValue
	ptr := (*[4]float32)(unsafe.Pointer(&c))
	ptr[0], ptr[1], ptr[2], ptr[3] = r, g, b, a
	return c
}
