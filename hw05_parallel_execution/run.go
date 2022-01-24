package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func runner(task Task, errChan chan<- error, doneChan chan<- struct{}) {
	err := task()
	doneChan <- struct{}{}
	if err != nil {
		errChan <- err
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.

	tasksLen := len(tasks)

	errChan := make(chan error)
	doneChan := make(chan struct{})
	errCount := 0

	pendingTasksIndex := -1
	doneCounter := 0
	if len(tasks) > n {
		pendingTasksIndex = n
	}
	for i := 0; i < n || i < tasksLen; i++ {
		go runner(tasks[i], errChan, doneChan)
	}
Mark:
	for {
		select {
		case <-doneChan:
			doneCounter++
			if doneCounter == tasksLen {
				break Mark
			}
			if pendingTasksIndex > 0 && pendingTasksIndex < tasksLen {
				go runner(tasks[pendingTasksIndex], errChan, doneChan)
				pendingTasksIndex++
			}
		case <-errChan:
			errCount++
			if errCount == m {
				return ErrErrorsLimitExceeded
			}
		}
	}
	return nil
}
