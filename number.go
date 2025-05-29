package glua

/*
#include "c/glua.h"
*/
import "C"
import (
	"strconv"

	"golang.org/x/exp/constraints"
)

const LuaNumberMaxSafeInteger int64 = (1 << 53) - 1

func pushNumber(L State, n LUA_NUMBER) {
	C.lua_pushnumber_wrap(L.c(), C.double(n))
}

func pushInt[V constraints.Integer](L State, n V) {
	// max(-n, n) is a quick way to get the absolute value of n, since golang doesn't have a built-in abs function for integers :D
	if max(-n, n) <= V(LuaNumberMaxSafeInteger) {
		pushNumber(L, LUA_NUMBER(n))
	} else {
		if n < 0 {
			L.PushString(strconv.FormatInt(int64(n), 10))
		} else {
			L.PushString(strconv.FormatUint(uint64(n), 10))
		}
	}
}

func (L State) PushNumber(n any) {
	switch v := n.(type) {
	case int:
		pushInt(L, v)
	case int8:
		pushNumber(L, LUA_NUMBER(v))
	case int16:
		pushNumber(L, LUA_NUMBER(v))
	case int32:
		pushNumber(L, LUA_NUMBER(v))
	case int64:
		pushInt(L, v)

	case uint:
		pushInt(L, v)
	case uint8:
		pushNumber(L, LUA_NUMBER(v))
	case uint16:
		pushNumber(L, LUA_NUMBER(v))
	case uint32:
		pushNumber(L, LUA_NUMBER(v))
	case uint64:
		pushInt(L, v)

	case float32:
		pushNumber(L, LUA_NUMBER(v))
	case float64:
		pushNumber(L, LUA_NUMBER(v))
	}
}
