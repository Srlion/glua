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
)

// The usage for C as a Think callback is because CGO is slow, so we need to only call it ONLY when we need to.

var (
	thinkQueue   chan GoFunc
	thinkFuncs   []GoFunc
	thinkFuncsMu sync.Mutex
)

func InitThinkQueue(L State) {
	thinkQueue = make(chan GoFunc, 100) // Make a buffered channel for the think queue
	thinkFuncs = make([]GoFunc, 0)      // Make a slice for the think functions
	thinkFuncsMu = sync.Mutex{}

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

	thinkFuncsMu.Lock()
	snapshot := thinkFuncs[:len(thinkFuncs):len(thinkFuncs)]
	thinkFuncsMu.Unlock()

	toRemove := map[int]struct{}{}
	for i, fn := range snapshot {
		L.SetTop(0)                   // completely empty the lua stack
		res, err := callGoFunc(L, fn) // we use callGoFunc to safely handle panics
		if err != nil {
			L.ErrorNoHalt(err.Error())
		}
		if res == 1 {
			count++
			toRemove[i] = struct{}{}
		}
	}

	if len(toRemove) > 0 {
		thinkFuncsMu.Lock()
		newThinkFuncs := make([]GoFunc, 0, len(thinkFuncs)-len(toRemove))
		for i, fn := range thinkFuncs {
			if _, ok := toRemove[i]; !ok {
				newThinkFuncs = append(newThinkFuncs, fn)
			}
		}
		thinkFuncs = newThinkFuncs
		thinkFuncsMu.Unlock()
	}

loop:
	for i := 0; i < 3; i++ { // Process 3 tasks per frame to prevent lag spikes
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

	thinkQueue <- fn

	C.increment_tasks_count() // concurrent increment
}

// LuaThink is a function that will be called every frame
// This function is thread-safe
// Return 1 to stop calling the function
func LuaThink(fn GoFunc) {
	if IS_STATE_OPEN.Load() == false {
		return
	}

	thinkFuncsMu.Lock()
	{
		thinkFuncs = append(thinkFuncs, fn) // Add the function to the think functions
	}
	thinkFuncsMu.Unlock()

	C.increment_tasks_count() // concurrent increment
}

func (L State) PollThinkQueue() {
	thinkQueueProcess(L)
}
