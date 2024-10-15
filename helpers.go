package glua

import "C"
import (
	"bytes"
	"unsafe"
)

func goStringN(p *C.char, n C.size_t) string {
	return string(unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
}

func goBytes(p unsafe.Pointer, n C.size_t) []byte {
	return bytes.Clone(unsafe.Slice((*byte)(p), n))
}
