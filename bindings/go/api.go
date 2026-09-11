package api_bindings

import (
	"unsafe"
)

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}
type BlockType struct {
	Name    string
	Color   Color
	IsSolid bool
}

//go:wasmimport engine_bindings log-info
func engine_bindings_LogInfo(msg_str uint32, msg_len uint32)

func LogInfo(msg string) {
	var msg_str uint32
	var msg_len uint32
	msg_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(msg))))
	msg_len = (uint32)(len(msg))
	engine_bindings_LogInfo(msg_str, msg_len)
}

//go:wasmimport engine_bindings log-warn
func engine_bindings_LogWarn(msg_str uint32, msg_len uint32)

func LogWarn(msg string) {
	var msg_str uint32
	var msg_len uint32
	msg_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(msg))))
	msg_len = (uint32)(len(msg))
	engine_bindings_LogWarn(msg_str, msg_len)
}

//go:wasmimport engine_bindings log-error
func engine_bindings_LogError(msg_str uint32, msg_len uint32)

func LogError(msg string) {
	var msg_str uint32
	var msg_len uint32
	msg_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(msg))))
	msg_len = (uint32)(len(msg))
	engine_bindings_LogError(msg_str, msg_len)
}

//go:wasmimport engine_bindings register-block
func engine_bindings_RegisterBlock(block_Name_f_str uint32, block_Name_f_len uint32, block_IsSolid_f_b uint32, block_Color_f_R_f_w uint32, block_Color_f_G_f_w uint32, block_Color_f_B_f_w uint32, block_Color_f_A_f_w uint32, namespaceid_str uint32, namespaceid_len uint32)

func RegisterBlock(block BlockType) string {
	var block_Name_f string
	block_Name_f = (string)(block.Name)
	var block_Name_f_str uint32
	var block_Name_f_len uint32
	block_Name_f_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(block_Name_f))))
	block_Name_f_len = (uint32)(len(block_Name_f))
	var block_Color_f Color
	block_Color_f = (Color)(block.Color)
	var block_Color_f_R_f uint8
	block_Color_f_R_f = (uint8)(block_Color_f.R)
	var block_Color_f_R_f_w uint32
	block_Color_f_R_f_w = (uint32)(block_Color_f_R_f)
	var block_Color_f_G_f uint8
	block_Color_f_G_f = (uint8)(block_Color_f.G)
	var block_Color_f_G_f_w uint32
	block_Color_f_G_f_w = (uint32)(block_Color_f_G_f)
	var block_Color_f_B_f uint8
	block_Color_f_B_f = (uint8)(block_Color_f.B)
	var block_Color_f_B_f_w uint32
	block_Color_f_B_f_w = (uint32)(block_Color_f_B_f)
	var block_Color_f_A_f uint8
	block_Color_f_A_f = (uint8)(block_Color_f.A)
	var block_Color_f_A_f_w uint32
	block_Color_f_A_f_w = (uint32)(block_Color_f_A_f)
	var block_IsSolid_f bool
	block_IsSolid_f = (bool)(block.IsSolid)
	var block_IsSolid_f_b uint32
	if block_IsSolid_f {
		block_IsSolid_f_b = (uint32)(1)
	} else {
		block_IsSolid_f_b = (uint32)(0)
	}
	var namespaceid_str uint32
	var namespaceid_len uint32
	engine_bindings_RegisterBlock(block_Name_f_str, block_Name_f_len, block_IsSolid_f_b, block_Color_f_R_f_w, block_Color_f_G_f_w, block_Color_f_B_f_w, block_Color_f_A_f_w, (uint32)(uintptr(unsafe.Pointer(&namespaceid_str))), (uint32)(uintptr(unsafe.Pointer(&namespaceid_len))))
	var namespaceid string
	namespaceid = (string)(unsafe.String((*byte)(unsafe.Pointer(uintptr(namespaceid_str))), int(namespaceid_len)))
	uResetMalloc()
	return namespaceid
}

//go:wasmimport engine_bindings append-lang
func engine_bindings_AppendLang(name_str uint32, name_len uint32)

func AppendLang(name string) {
	var name_str uint32
	var name_len uint32
	name_str = (uint32)(uintptr(unsafe.Pointer(unsafe.StringData(name))))
	name_len = (uint32)(len(name))
	engine_bindings_AppendLang(name_str, name_len)
}
