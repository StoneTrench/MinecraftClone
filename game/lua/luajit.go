package lua

/*
#cgo windows CFLAGS: -I${SRCDIR}/include
#cgo windows LDFLAGS: -L${SRCDIR}/include -lluajit-5.1

#include <stdlib.h>
#include <stdio.h>

#ifdef lua_strlen
#undef lua_strlen
#endif

#include "include/luaconf.h"
#include "include/lua.h"
#include "include/lualib.h"
#include "include/lauxlib.h"
#include "include/luajit.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// func lua_pop(L *C.lua_State, n C.int) {
// 	C.lua_settop(L, -(n)-1)
// }

type State struct {
	L *C.lua_State
}

func (s State) Init() *State {
	L := C.luaL_newstate()
	if L == nil {
		return nil
	}
	C.luaL_openlibs(L)
	return &State{L: L}
}

func (s *State) DoString(code string) error {
	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	if C.luaL_loadstring(s.L, cCode) != 0 {
		return errors.New(s.getLastError())
	}
	if C.lua_pcall(s.L, 0, C.LUA_MULTRET, 0) != 0 {
		return fmt.Errorf("lua error: %s", s.getLastError())
	}
	return nil
}

func (s *State) getLastError() string {
	errStr := C.GoString(C.lua_tolstring(s.L, -1, nil))
	C.lua_settop(s.L, -2)
	return errStr
}

func (s *State) Deinit() {
	if s.L != nil {
		C.lua_close(s.L)
		s.L = nil
	}
}
