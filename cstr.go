package glua

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type CString struct {
	c *C.char
}

func (cstr *CString) free() {
	if cstr.c != nil {
		C.free(unsafe.Pointer(cstr.c))
		cstr.c = nil
	}
}

func CStr(gs string) *CString {
	cstr := &CString{C.CString(gs)}
	return cstr
}

func NoHeapCStr(str string) *C.char {
	cstr := unsafe.StringData(str)
	return (*C.char)(unsafe.Pointer(cstr))
}
