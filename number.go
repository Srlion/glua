package glua

/*
#include "c/glua.h"
*/
import "C"
import (
	"strconv"

	"golang.org/x/exp/constraints"
)

const LuaNumberMaxSafeInteger = (1 << 53) - 1

func pushNumber(L State, n LUA_NUMBER) {
	C.lua_pushnumber_wrap(L.c(), C.double(n))
}

func pushIntSigned[T constraints.Signed](L State, n T) {
	const max = int64(LuaNumberMaxSafeInteger)
	v := int64(n)
	if v >= -max && v <= max {
		pushNumber(L, LUA_NUMBER(v))
		return
	}
	L.PushString(strconv.FormatInt(v, 10))
}

func pushIntUnsigned[T constraints.Unsigned](L State, n T) {
	const max = uint64(LuaNumberMaxSafeInteger)
	v := uint64(n)
	if v <= max {
		pushNumber(L, LUA_NUMBER(v))
		return
	}
	L.PushString(strconv.FormatUint(v, 10))
}

func (L State) PushNumber(n any) {
	switch v := n.(type) {
	case int:
		pushIntSigned(L, v)
	case int8:
		pushNumber(L, LUA_NUMBER(v))
	case int16:
		pushNumber(L, LUA_NUMBER(v))
	case int32:
		pushNumber(L, LUA_NUMBER(v))
	case int64:
		pushIntSigned(L, v)

	case uint:
		pushIntUnsigned(L, v)
	case uint8:
		pushNumber(L, LUA_NUMBER(v))
	case uint16:
		pushNumber(L, LUA_NUMBER(v))
	case uint32:
		pushNumber(L, LUA_NUMBER(v))
	case uint64:
		pushIntUnsigned(L, v)

	case float32:
		pushNumber(L, LUA_NUMBER(v))
	case float64:
		pushNumber(L, LUA_NUMBER(v))
	}
}
