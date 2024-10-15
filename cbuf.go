package glua

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// CBytes struct to hold C-allocated byte buffer
type CBytes struct {
	c    *C.char
	size C.size_t
}

// Free releases the C-allocated memory
func (cb *CBytes) Free() {
	if cb.c != nil {
		C.free(unsafe.Pointer(cb.c))
		cb.c = nil
		cb.size = 0
	}
}

func CByt(buf []byte) *CBytes {
	if len(buf) == 0 {
		return &CBytes{c: (*C.char)(unsafe.Pointer(nil)), size: 0}
	}

	cbuf := C.CBytes(buf)
	return &CBytes{
		c:    (*C.char)(cbuf),
		size: C.size_t(len(buf)),
	}
}
