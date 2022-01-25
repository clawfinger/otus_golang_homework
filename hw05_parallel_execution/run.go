package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func runner(taskChan chan Task, errChan chan error, wg *sync.WaitGroup, stopped *int32) {
End:
	for {
		select {
		case task := <-taskChan:
			err := task()
			if atomic.LoadInt32(stopped) == 1 {
				break End
			}
			errChan <- err
		default:
			if atomic.LoadInt32(stopped) == 1 {
				break End
			}
		}
	}
	wg.Done()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.

	taskCount := len(tasks) + 1
	var wg sync.WaitGroup
	tasksCounter := 0

	errChan := make(chan error)
	taskChan := make(chan Task, taskCount)

	var stopped int32

	errCount := 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go runner(taskChan, errChan, &wg, &stopped)
	}

	for _, task := range tasks {
		taskChan <- task
	}
	var res error
	for {
		err := <-errChan
		tasksCounter++
		if err != nil {
			errCount++
		}
		if errCount == m {
			res = ErrErrorsLimitExceeded
			break
		}
		if tasksCounter == len(tasks) {
			break
		}
	}
	atomic.StoreInt32(&stopped, 1)
	wg.Wait()
	return res
}
