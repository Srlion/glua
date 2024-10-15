package glua

/*
#include "c/glua.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type GoFunc = func(L State) int

var (
	funcMap   = make(map[uint64]GoFunc)
	funcMapMu sync.Mutex
	nextID    uint64 = 0
)

func registerGoFunc(fn GoFunc, oneTimeUse bool) uint64 {
	funcMapMu.Lock()
	defer funcMapMu.Unlock()

	id := nextID
	if oneTimeUse {
		funcMap[id] = func(L State) int {
			ret := fn(L)
			unRegisterGoFunc(id)
			return ret
		}
	} else {
		funcMap[id] = fn
	}
	nextID++

	return id
}

func getGoFunc(id uint64) (GoFunc, bool) {
	funcMapMu.Lock()
	defer funcMapMu.Unlock()

	fn, ok := funcMap[id]
	return fn, ok
}

func unRegisterGoFunc(id uint64) {
	funcMapMu.Lock()
	defer funcMapMu.Unlock()

	delete(funcMap, id)
}

//export goLuaCallback
func goLuaCallback(L State, res *C.int, err **C.char) {
	// Get the upvalue, which should be the function handle (id)
	id := L.GetNumber(lua_upvalueindex(1))

	// Retrieve the Go function using the handle
	fn, ok := getGoFunc(uint64(id))
	if !ok {
		fmt.Printf("Failed to get Go function with handle %v\n", id)
		*err = C.CString("Failed to get Go function") // it will be freed by the C side
		return
	}

	// allow panics to return to Lua with no issues, by making the C wrapper for check w/e u return to
	// detect if its an error or not to cause a lua error, we want that to happen from C side so we make sure go side is safe
	// im not sure if calling lua_error from go side will cause issues with the go runtime, so we just let the C side handle it
	// this won't cause issues if go functions panic when they are called from go functions by lua, eg. go > lua > go as long as it's called with PCall
	// we use panics so we don't have to keep checking for err with every single func call, could be costly but idgaf
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Go function panic: %v\n", r)
			*err = C.CString(fmt.Sprintf("%v", r)) // it will be freed by the C side
		}
	}()

	*res = C.int(fn(L))
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
	PushInteger(L, handle)

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
	PushInteger(L, handle)

	L.PushCClosure(C.lua_call_go, 1)
}

func (L State) PushCFunc(fn unsafe.Pointer) {
	C.lua_pushcclosure_wrap(L.c(), fn, 0)
}

func (L State) PushCClosure(fn unsafe.Pointer, n int) {
	C.lua_pushcclosure_wrap(L.c(), fn, C.int(n))
}
