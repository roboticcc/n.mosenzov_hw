package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors - concurrency check", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		n := 3
		taskCount := 1200

		var concurrentTasks int32
		var maxConcurrentTasks int32
		var mu sync.Mutex
		done := make(chan struct{})

		tasks := make([]Task, taskCount)
		for i := 0; i < taskCount; i++ {
			tasks[i] = func() error {
				current := atomic.AddInt32(&concurrentTasks, 1)

				mu.Lock()
				if current > maxConcurrentTasks {
					maxConcurrentTasks = current
				}
				mu.Unlock()

				for j := 0; j < 100000; j++ {
					_ = j * j
				}

				atomic.AddInt32(&concurrentTasks, -1)
				return nil
			}
		}

		errChan := make(chan error, 1)
		go func() {
			err := Run(tasks, n, 1)
			errChan <- err
			close(done)
		}()

		<-done
		err := <-errChan
		require.NoError(t, err)

		mu.Lock()
		require.Equal(t, int32(n), maxConcurrentTasks, "Expected %d tasks to run concurrently, but got max %d", n, maxConcurrentTasks)
		mu.Unlock()

		require.Equal(t, int32(0), atomic.LoadInt32(&concurrentTasks), "All tasks should be completed")
	})
	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}
