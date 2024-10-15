package glua

/*
#include "c/glua.h"
*/
import "C"
import (
	"strconv"

	"golang.org/x/exp/constraints"
)

const LUA_NUMBER_MAX_SAFE_INTEGER uint64 = 9007199254740991
const LUA_NUMBER_MIN_SAFE_INTEGER int64 = -9007199254740991

func (L State) pushNumber(n LUA_NUMBER) {
	C.lua_pushnumber_wrap(L.c(), C.double(n))
}

// Methods can't use generics yet

/*
Pushes an integer onto the stack.
*/
func PushInteger[V constraints.Integer](L State, value V) {
	if value > 0 {
		if uint64(value) <= LUA_NUMBER_MAX_SAFE_INTEGER {
			L.pushNumber(LUA_NUMBER(value))
		} else {
			L.PushString(strconv.FormatUint(uint64(value), 10))
		}
	} else if value < 0 {
		if int64(value) >= LUA_NUMBER_MIN_SAFE_INTEGER {
			L.pushNumber(LUA_NUMBER(value))
		} else {
			L.PushString(strconv.FormatInt(int64(value), 10))
		}
	} else {
		L.pushNumber(0)
	}
}

/*
Pushes a float onto the stack.
*/
func PushFloat[V constraints.Float](L State, value V) {
	L.pushNumber(LUA_NUMBER(value))
}
