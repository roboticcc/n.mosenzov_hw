package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return fmt.Errorf("number of workers must be positive, got %d", n)
	}
	if len(tasks) == 0 {
		return nil
	}

	tasksChan := make(chan Task, len(tasks))
	stopChan := make(chan struct{})
	defer close(stopChan)

	var errorCount int32
	var wg sync.WaitGroup

	for _, task := range tasks {
		tasksChan <- task
	}
	close(tasksChan)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(tasksChan, stopChan, &wg, &errorCount, m)
	}

	wg.Wait()

	if m > 0 && int(atomic.LoadInt32(&errorCount)) > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(tasksChan <-chan Task, stopChan chan struct{}, wg *sync.WaitGroup, errorCount *int32, m int) {
	defer wg.Done()

	for {
		select {
		case <-stopChan:
			return
		case task, ok := <-tasksChan:
			if !ok {
				return
			}

			if err := task(); err != nil {
				newCount := atomic.AddInt32(errorCount, 1)
				if m > 0 && int(newCount) > m {
					signalStop(stopChan)
					return
				}
			}
		}
	}
}

func signalStop(stopChan chan<- struct{}) {
	select {
	case stopChan <- struct{}{}:
	default:
	}
}
