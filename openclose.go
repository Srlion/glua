package glua

/*
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"sync/atomic"
)

var GMOD13_OPEN func(L State) int
var GMOD13_CLOSE func(L State) int

var IS_STATE_OPEN = atomic.Bool{}

//export go_gmod13_open
func go_gmod13_open(L State) C.int {
	err := LoadLuaShared()
	if err != nil {
		fmt.Printf("Error loading lua shared: %v\n", *err)
		return 0
	}

	IS_STATE_OPEN.Store(true)

	InitGoTasks(L)
	InitGoPtrRegistry(L)
	InitGoFuncRegistry(L)
	InitThinkQueue(L)

	if GMOD13_OPEN != nil {
		return C.int(GMOD13_OPEN(L))
	}

	return 0
}

//export go_gmod13_close
func go_gmod13_close(L State) C.int {
	var res C.int = 0

	WaitGoTasks()

	if GMOD13_CLOSE != nil {
		res = C.int(GMOD13_CLOSE(L))
	}

	IS_STATE_OPEN.Store(false)

	UnloadLuaShared()

	return res
}
