package glua

import "C"
import (
	"unsafe"

	"github.com/Srlion/safereg"
)

var PtrRegistry *safereg.Registry
var quick_go_ptr any
var quick_go_ptr_handle uintptr

func InitGoPtrRegistry(L State) {
	PtrRegistry = safereg.New()
	quick_go_ptr = 0
	quick_go_ptr_handle = PtrRegistry.Store(0) // we just use it as first handle for quick go ptr
}

// This mallocs a new pointer to a Go object and returns the struct
func NewGoPointer(val any) unsafe.Pointer {
	newPtr := PtrRegistry.Store(val)
	return unsafe.Pointer(newPtr)
}

// This unwraps a *C.uintptr_t back to what it was before
func UnwrapGoPointer[T any](ptr unsafe.Pointer) T {
	handle := uintptr(ptr)
	// if it's a quick go ptr then return the value
	if handle == quick_go_ptr_handle {
		return quick_go_ptr.(T)
	}
	val, _ := PtrRegistry.Get(handle)
	return val.(T)
}

// This can be used when calling a C function to bypass the rule of not passing Go pointers to C
//
// YOU SHOULD ONLY USE THIS IF YOU WON"T BE STORING THE POINTER IN C!
//
// THIS SHOULDN'T BE USED CONCURRENTLY!
func QuickGoPtr(val any) unsafe.Pointer {
	quick_go_ptr = val
	return unsafe.Pointer(quick_go_ptr_handle)
}
