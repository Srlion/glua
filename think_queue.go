package glua

/*
#include "c/think_queue.h"
#include "c/glua.h"
*/
import "C"
import (
	"math/rand"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// The usage for C as a Think callback is because CGO is slow, so we need to only call it ONLY when we need to.

var thinkQueue chan unsafe.Pointer
var thinkQueueCountMu sync.Mutex

func InitThinkQueue(L State) {
	thinkQueue = make(chan unsafe.Pointer, 100) // We need to use a buffered channel to make use of queueing
	thinkQueueCountMu = sync.Mutex{}
	C.reset_tasks_count()

	L.GetGlobal("timer")
	{
		L.GetField(-1, "Create")
		{
			randomUniqueName := strconv.FormatInt(int64(rand.Int()), 10) + strconv.FormatInt(time.Now().UnixNano(), 10)
			L.PushString("GoLuaThinkQueue" + randomUniqueName)
			L.PushNumber(0) // Delay (0 = next frame)
			L.PushNumber(0) // Repetitions (0 = infinite)
			L.PushCFunc(C.think_queue_think)
		}
		L.TryCall(4, 0)
	}
	L.Pop()
}

//export thinkQueueProcess
func thinkQueueProcess(L State) {
	count := 0
loop:
	for {
		select {
		case task := <-thinkQueue:
			L.TryCPCall(C.lua_cpcall_go, task)
			count++
		default:
			// No more tasks to process
			break loop
		}
	}

	// Decrement the task count by the number of tasks processed
	if count > 0 { // Only lock if there are tasks to decrement
		thinkQueueCountMu.Lock()
		{
			C.decrement_tasks_count_by(C.int(count))
		}
		thinkQueueCountMu.Unlock()
	}
}

func WaitLuaThink(fn GoFunc) {
	if !IS_STATE_OPEN {
		return
	}

	thinkQueue <- registerGoFunc(fn, true)

	thinkQueueCountMu.Lock()
	{
		C.increment_tasks_count()
	}
	thinkQueueCountMu.Unlock()
}

func (L State) PollThinkQueue() {
	thinkQueueProcess(L)
}
