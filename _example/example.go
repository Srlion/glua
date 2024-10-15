package main

import (
	"github.com/Srlion/glua"
)

func init() {
	glua.GMOD13_OPEN = gmod13_open
	glua.GMOD13_CLOSE = gmod13_close
}

func test(L glua.State) int {
	var str string = L.CheckString(1)
	println("string:", str)

	L.PushString("Hello from Go!")
	return 1
}

func gmod13_open(L glua.State) int {
	println("hello from gmod13_open!")

	L.PushGoFunc(test)
	L.SetGlobal("test")

	return 0
}

func gmod13_close(L glua.State) int {
	return 0
}

// Required by Go when using `-buildmode=c-shared`
func main() {}
