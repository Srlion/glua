package glua

import "C"
import "unsafe"

var quick_go_ptr = NewGoPointer(nil)

// This mallocs a new pointer to a Go object and returns the struct
func NewGoPointer(ptr unsafe.Pointer) unsafe.Pointer {
	newPtr := (*C.uintptr_t)(C.malloc(C.sizeof_uintptr_t))
	*newPtr = C.uintptr_t(uintptr(ptr))
	return unsafe.Pointer(newPtr)
}

// This unwraps a *C.uintptr_t back to what it was before
func UnwrapGoPointer(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(*(*uintptr)(ptr))
}

// This can be used when calling a C function to bypass the rule of not passing Go pointers to C
//
// YOU SHOULD ONLY USE THIS IF YOU WON"T BE STORING THE POINTER IN C!
//
// THIS SHOULDN'T BE USED CONCURRENTLY!
func QuickGoPtr(ptr unsafe.Pointer) unsafe.Pointer {
	*(*uintptr)(quick_go_ptr) = uintptr(ptr)
	return quick_go_ptr
}
