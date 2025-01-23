package glua

import "sync"

var TasksWG sync.WaitGroup

func InitGoTasks(L State) {
	TasksWG = sync.WaitGroup{}
}

func Go(fn func()) {
	TasksWG.Add(1)
	go func() {
		defer TasksWG.Done()
		fn()
	}()
}

func WaitGoTasks() {
	TasksWG.Wait()
}
