package glua

/*
#include <stdint.h>
extern int gmod13_open(uintptr_t L);
extern int gmod13_close(uintptr_t L);
*/
import "C"
import "fmt"

var GMOD13_OPEN func(L State) int
var GMOD13_CLOSE func(L State) int

//export gmod13_open
func gmod13_open(L State) C.int {
	err := LoadLuaShared()
	if err != nil {
		fmt.Printf("Error loading lua shared: %v\n", *err)
		return 0
	}

	if GMOD13_OPEN != nil {
		return C.int(GMOD13_OPEN(L))
	}

	return 0
}

//export gmod13_close
func gmod13_close(L State) C.int {
	var res C.int = 0

	if GMOD13_CLOSE != nil {
		res = C.int(GMOD13_CLOSE(L))
	}

	UnloadLuaShared()

	return res
}
