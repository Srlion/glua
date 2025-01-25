package glua

/*
#include "c/think_queue.h"
#include "c/glua.h"
*/
import "C"
import (
	"math/rand"
	"slices"
	"strconv"
	"sync"
	"time"
)

// The usage for C as a Think callback is because CGO is slow, so we need to only call it ONLY when we need to.

var (
	thinkQueue   chan GoFunc
	thinkQueueMu sync.Mutex

	thinkFuncs []GoFunc
)

func InitThinkQueue(L State) {
	thinkQueue = make(chan GoFunc, 100) // We need to use a buffered channel to make use of queueing
	thinkQueueMu = sync.Mutex{}

	thinkFuncs = make([]GoFunc, 0)

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
	thinkQueueMu.Lock()
	defer thinkQueueMu.Unlock()

	count := 0

	// process think functions
	thinkFuncs = slices.DeleteFunc(thinkFuncs, func(fn GoFunc) bool {
		L.SetTop(0)                   // completely empty the lua stack
		res, err := callGoFunc(L, fn) // we use callGoFunc to safely handle panics
		if err != nil {
			L.ErrorNoHalt(err.Error())
		}
		if res == 1 {
			count++
			return true // Remove the function from the think functions
		}
		return false
	})

loop:
	for {
		select {
		case task := <-thinkQueue:
			L.SetTop(0)                   // completely empty the lua stack
			_, err := callGoFunc(L, task) // we use callGoFunc to safely handle panics
			if err != nil {
				L.ErrorNoHalt(err.Error())
			}
			count++
		default:
			// No more tasks to process
			break loop
		}
	}

	C.decrement_tasks_count_by(C.int(count))
}

func WaitLuaThink(fn GoFunc) {
	if IS_STATE_OPEN.Load() == false {
		return
	}

	thinkQueue <- fn // Queue the task

	thinkQueueMu.Lock()
	{
		C.increment_tasks_count()
	}
	thinkQueueMu.Unlock()
}

// LuaThink is a function that will be called every frame
// This function is thread-safe
// Return 1 to stop calling the function
func LuaThink(fn GoFunc) {
	if IS_STATE_OPEN.Load() == false {
		return
	}

	thinkQueueMu.Lock()
	{
		thinkFuncs = append(thinkFuncs, fn) // Add the function to the think functions
		C.increment_tasks_count()
	}
	thinkQueueMu.Unlock()
}

func (L State) PollThinkQueue() {
	thinkQueueProcess(L)
}
