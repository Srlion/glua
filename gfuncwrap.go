package glua

/*
#include "c/glua.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/Srlion/safereg"
)

type GoFunc = func(L State) int

var FuncRegistry *safereg.Registry

func InitGoFuncRegistry(L State) {
	FuncRegistry = safereg.New()
}

func registerGoFunc(fn GoFunc, oneTimeUse bool) uintptr {
	var handle uintptr
	if oneTimeUse {
		handle = FuncRegistry.Store(func(L State) int {
			FuncRegistry.Release(handle)
			return fn(L)
		})
	} else {
		handle = FuncRegistry.Store(fn)
	}
	return handle
}

func callGoFunc(L State, fn GoFunc) (res int, err error) {
	// allow panics to return to Lua with no issues, by making the C wrapper for check w/e u return to
	// detect if its an error or not to cause a lua error, we want that to happen from C side so we make sure go side is safe
	// im not sure if calling lua_error from go side will cause issues with the go runtime, so we just let the C side handle it
	// this won't cause issues if go functions panic when they are called from go functions by lua, eg. go > lua > go as long as it's called with PCall
	// we use panics so we don't have to keep checking for err with every single func call, could be costly but idgaf
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	res = fn(L)
	return
}

//export goLuaCallback
func goLuaCallback(L State, cRes *C.int, cErr **C.char) {
	var res int
	var err error

	// Get the function handle from the upvalue
	funcHandle := L.GetLightUserData(lua_upvalueindex(1))
	fn, exists := FuncRegistry.Get(funcHandle)
	if !exists {
		err = errors.New("attempt to call a nil value")
		goto handleRet
	}

	res, err = callGoFunc(L, fn.(GoFunc))

handleRet:
	if err != nil {
		*cErr = C.CString(err.Error())
	} else {
		*cRes = C.int(res)
	}
}

//	Pushes a Go function to the Lua stack.
//	If the function you pushing will be used one time, then use PushOneTimeGoFunc
//
// # Example
//
//	func printHello(L glua.State) int {
//		fmt.Println("Hello from Go!")
//		return 0
//	}
//	L := glua.NewState()
//	L.PushGoFunc(printHello)
//	L.SetGlobal("test")
func (L State) PushGoFunc(goFunc GoFunc) {
	handle := registerGoFunc(goFunc, false)
	L.PushLightUserData(handle)
	L.PushCClosure(C.lua_call_go, 1)
}

//	Pushes a Go function to the Lua stack that will be used only once.
//	After the function is called, it will be unregistered.
//
// # Example:
//
//	L := glua.NewState()
//	L.PushOneTimeGoFunc(func(L glua.State) int {
//		fmt.Println("Hello from Go!")
//		return 0
//	})
//	L.SetGlobal("test")
func (L State) PushOneTimeGoFunc(goFunc GoFunc) {
	handle := registerGoFunc(goFunc, true)
	L.PushLightUserData(handle)
	L.PushCClosure(C.lua_call_go, 1)
}

func (L State) PushCFunc(fn unsafe.Pointer) {
	C.lua_pushcclosure_wrap(L.c(), fn, 0)
}

func (L State) PushCClosure(fn unsafe.Pointer, n int) {
	C.lua_pushcclosure_wrap(L.c(), fn, C.int(n))
}
