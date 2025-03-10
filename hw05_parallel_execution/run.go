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

	tasksChan := make(chan Task)
	errorsChan := make(chan struct{}, n)
	var errorCount int32
	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(tasksChan, errorsChan, &wg, &errorCount)
	}

	go func() {
		defer close(tasksChan)
		for _, task := range tasks {
			tasksChan <- task
		}
	}()

	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	for range errorsChan {
		if m > 0 && int(atomic.LoadInt32(&errorCount)) > m {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func worker(tasksChan <-chan Task, errorsChan chan<- struct{}, wg *sync.WaitGroup, errorCount *int32) {
	defer wg.Done()

	for task := range tasksChan {
		if err := task(); err != nil {
			atomic.AddInt32(errorCount, 1)
			select {
			case errorsChan <- struct{}{}:
			default:
			}
		}
	}
}
