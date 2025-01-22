package glua

/*
#include "c/glua.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

type GoFunc = func(L State) int

var (
	funcMap   map[unsafe.Pointer]GoFunc
	funcMapMu sync.RWMutex
)

func InitGoFuncs(L State) {
	funcMap = make(map[unsafe.Pointer]GoFunc)
	funcMapMu = sync.RWMutex{}
}

func registerGoFunc(fn GoFunc, oneTimeUse bool) unsafe.Pointer {
	funcMapMu.Lock()
	defer funcMapMu.Unlock()

	funcPtr := unsafe.Pointer(&fn)

	if oneTimeUse {
		funcMap[funcPtr] = func(L State) int {
			ret := fn(L)
			unRegisterGoFunc(funcPtr)
			return ret
		}
	} else {
		funcMap[funcPtr] = fn
	}

	return funcPtr
}

func getGoFunc(funcPtr unsafe.Pointer) (GoFunc, bool) {
	funcMapMu.RLock()
	defer funcMapMu.RUnlock()

	fn, ok := funcMap[funcPtr]
	return fn, ok
}

func unRegisterGoFunc(funcPtr unsafe.Pointer) {
	funcMapMu.Lock()
	defer funcMapMu.Unlock()

	delete(funcMap, funcPtr)
}

func callGoFunc(L State, funcPtr unsafe.Pointer) (res int, err error) {
	// Retrieve the Go function using the handle
	fn, ok := getGoFunc(funcPtr)
	if !ok {
		err = errors.New("no go function found")
		return
	}

	// allow panics to return to Lua with no issues, by making the C wrapper for check w/e u return to
	// detect if its an error or not to cause a lua error, we want that to happen from C side so we make sure go side is safe
	// im not sure if calling lua_error from go side will cause issues with the go runtime, so we just let the C side handle it
	// this won't cause issues if go functions panic when they are called from go functions by lua, eg. go > lua > go as long as it's called with PCall
	// we use panics so we don't have to keep checking for err with every single func call, could be costly but idgaf
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("%v", r)
			println(errMsg)
			err = errors.New(errMsg)
		}
	}()

	res = fn(L)
	return
}

//export goLuaCallback
func goLuaCallback(L State, cRes *C.int, cErr **C.char) {
	ptr := L.GetLightUserData(lua_upvalueindex(1))
	res, err := callGoFunc(L, ptr)
	if err != nil {
		*cErr = C.CString(err.Error())
	} else {
		*cRes = C.int(res)
	}
}

//export goLuaCPCallback
func goLuaCPCallback(L State, cErr **C.char) {
	ptr := UnwrapGoPointer(L.GetLightUserData(1))
	if _, err := callGoFunc(L, ptr); err != nil {
		*cErr = C.CString(err.Error())
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
