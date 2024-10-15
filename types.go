package glua

// #include "c/glua.h"
import "C"

type LUA_NUMBER float64

const LUA_MULTRET = -1

const LUA_REGISTRYINDEX = -10000
const LUA_ENVIRONINDEX = -10001
const LUA_GLOBALSINDEX int = -10002

const LUA_TNONE = -1
const LUA_TNIL = 0
const LUA_TBOOLEAN = 1
const LUA_TLIGHTUSERDATA = 2
const LUA_TNUMBER = 3
const LUA_TSTRING = 4
const LUA_TTABLE = 5
const LUA_TFUNCTION = 6
const LUA_TUSERDATA = 7
const LUA_TTHREAD = 8

const LUA_OK = 0
const LUA_YIELD = 1
const LUA_ERRRUN = 2
const LUA_ERRSYNTAX = 3
const LUA_ERRMEM = 4
const LUA_ERRERR = 5

const LUA_ERRFILE = LUA_ERRERR + 1

const LUA_REFNIL = -1
const LUA_NOREF = -2

const LUA_IDSIZE = 60

type State C.uintptr_t

func (L State) c() C.uintptr_t {
	return C.uintptr_t(L)
}

func lua_upvalueindex(i int) int {
	return LUA_GLOBALSINDEX - i
}
