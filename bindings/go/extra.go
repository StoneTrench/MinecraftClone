package api_bindings

import "unsafe"

var u_malloc_keepalive [][]byte

//go:wasmexport malloc
func uMalloc(size uint32) uint32 {
	if size == 0 {
		return 0
	}
	raw := make([]byte, size)
	u_malloc_keepalive = append(u_malloc_keepalive, raw)
	return uint32(uintptr(unsafe.Pointer(&raw[0])))
}

func uResetMalloc() {
	if len(u_malloc_keepalive) == 0 {
		return
	}
	u_malloc_keepalive = nil
}
