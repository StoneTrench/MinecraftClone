package host

import (
	"fmt"
	"unsafe"

	"github.com/tetratelabs/wazero/api"
)

type HostBinding struct {
	Name string
	Fn   interface{}
}
type FnAllocator func(size uint32) (offset uint32)

func uPanic(m api.Module, err error) {
	panic(fmt.Errorf("error occured in (%s): %w", m.Name(), err))
}
func uPanicf(m api.Module, format string, args ...any) {
	uPanic(m, fmt.Errorf(format, args...))
}
func uPanics(m api.Module, msg string) {
	uPanicf(m, "%s", msg)
}

func uReadString(m api.Module, offset, length uint32) string {
	if length == 0 {
		return ""
	}
	raw, ok := m.Memory().Read(offset, length)
	if !ok {
		uPanics(m, "out of bounds")
	}
	return string(raw)
}
func uReadPlex[T any](m api.Module, offset uint32) T {
	var v T
	size := uint32(unsafe.Sizeof(v))
	if size == 0 {
		return v
	}
	raw, ok := m.Memory().Read(offset, size)
	if !ok {
		uPanics(m, "out of bounds")
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&v)), size), raw)
	return v
}
func uReadSlice[T any](m api.Module, offset, length uint32) []T {
	if length == 0 {
		return nil
	}
	var zero T
	elemSize := uint32(unsafe.Sizeof(zero))
	size := length * elemSize
	raw, ok := m.Memory().Read(offset, size)
	if !ok {
		uPanics(m, "out of bounds")
	}
	out := make([]T, length)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(&out[0])), size), raw)
	return out
}
func uReadPointer[T any](m api.Module, offset uint32) *T {
	if offset == 0 {
		return nil
	}
	v := new(T)
	size := uint32(unsafe.Sizeof(*v))
	raw, ok := m.Memory().Read(offset, size)
	if !ok {
		uPanics(m, "out of bounds")
	}
	copy(unsafe.Slice((*byte)(unsafe.Pointer(v)), size), raw)
	return v
}

func uWriteString(m api.Module, alloc FnAllocator, str string, ptr_ptr uint32, ptr_len uint32) {
	length := uint32(len(str))
	if length == 0 {
		uWriteAnyFixed[uint32](m, 0, ptr_ptr)
		uWriteAnyFixed[uint32](m, 0, ptr_len)
		return
	}
	offset := alloc(length)
	ok := m.Memory().Write(offset, ([]byte)(str))
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[uint32](m, offset, ptr_ptr)
	uWriteAnyFixed[uint32](m, length, ptr_len)
}
func uWritePlex[T any](m api.Module, alloc FnAllocator, v T, ptr_ptr uint32) {
	size := uint32(unsafe.Sizeof(v))
	if size == 0 {
		uWriteAnyFixed[uint32](m, 0, ptr_ptr)
		return
	}
	offset := alloc(size)
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&v)), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[uint32](m, offset, ptr_ptr)
}
func uWriteSlice[T any](m api.Module, alloc FnAllocator, v []T, ptr_ptr uint32, ptr_len uint32) {
	if len(v) == 0 {
		uWriteAnyFixed[uint32](m, 0, ptr_ptr)
		uWriteAnyFixed[uint32](m, 0, ptr_len)
		return
	}
	var zero T
	elemSize := uint32(unsafe.Sizeof(zero))
	length := uint32(len(v))
	size := length * elemSize
	offset := alloc(size)
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&v[0])), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[uint32](m, offset, ptr_ptr)
	uWriteAnyFixed[uint32](m, length, ptr_len)
}
func uWritePointer[T any](m api.Module, alloc FnAllocator, v *T, ptr_ptr uint32) {
	if v == nil {
		uWriteAnyFixed[uint32](m, 0, ptr_ptr)
		return
	}
	size := uint32(unsafe.Sizeof(*v))
	offset := alloc(size)
	raw := unsafe.Slice((*byte)(unsafe.Pointer(v)), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
	uWriteAnyFixed[uint32](m, offset, ptr_ptr)
}
func uWriteAnyFixed[T any](m api.Module, v T, offset uint32) {
	size := uint32(unsafe.Sizeof(v))
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&v)), size)
	ok := m.Memory().Write(offset, raw)
	if !ok {
		uPanics(m, "out of bounds")
	}
}
