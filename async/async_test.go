package async

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelMap(t *testing.T) {
	t.Run("correct output and order", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []int{2, 4, 6, 8, 10}
		
		result := ParallelMap(input, func(x int) int {
			time.Sleep(10 * time.Millisecond) // Simulate work
			return x * 2
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("ParallelMap() = %v, want %v", result, expected)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		result := ParallelMap(input, func(x int) int { return x * 2 })
		
		if len(result) != 0 {
			t.Errorf("ParallelMap() with empty slice should return empty slice")
		}
	})

	t.Run("type transformation", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := []string{"1", "2", "3"}
		
		result := ParallelMap(input, func(x int) string {
			return fmt.Sprintf("%d", x)
		})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("ParallelMap() = %v, want %v", result, expected)
		}
	})
}

func TestErrGroup(t *testing.T) {
	t.Run("all succeed", func(t *testing.T) {
		eg := &ErrGroup{}
		
		for i := 0; i < 5; i++ {
			eg.Go(func() error {
				time.Sleep(10 * time.Millisecond)
				return nil
			})
		}

		err := eg.Wait()
		if err != nil {
			t.Errorf("ErrGroup.Wait() = %v, want nil", err)
		}
	})

	t.Run("one fails", func(t *testing.T) {
		eg := &ErrGroup{}
		
		for i := 0; i < 5; i++ {
			i := i // capture loop variable
			eg.Go(func() error {
				time.Sleep(10 * time.Millisecond)
				if i == 2 {
					return fmt.Errorf("error at %d", i)
				}
				return nil
			})
		}

		err := eg.Wait()
		if err == nil {
			t.Error("ErrGroup.Wait() expected error, got nil")
		}
	})

	t.Run("no goroutines", func(t *testing.T) {
		eg := &ErrGroup{}
		err := eg.Wait()
		
		if err != nil {
			t.Errorf("ErrGroup.Wait() with no goroutines = %v, want nil", err)
		}
	})
}

func TestPool(t *testing.T) {
	t.Run("executes tasks concurrently", func(t *testing.T) {
		pool := NewPool(3)
		defer pool.Close()

		var counter int64

		// Submit 10 tasks
		tasks := make([]func(), 10)
		for i := 0; i < 10; i++ {
			tasks[i] = func() {
				atomic.AddInt64(&counter, 1)
				time.Sleep(10 * time.Millisecond)
			}
		}

		pool.Submit(tasks...)
		pool.Wait()

		if counter != 10 {
			t.Errorf("Pool executed %d tasks, want 10", counter)
		}
	})

	t.Run("handles empty task list", func(t *testing.T) {
		pool := NewPool(2)
		defer pool.Close()

		pool.Submit() // Submit nothing
		pool.Wait()   // Should not hang
	})

	t.Run("handles nil tasks", func(t *testing.T) {
		pool := NewPool(2)
		defer pool.Close()

		pool.Submit(nil, func() {}, nil)
		pool.Wait()
	})

	t.Run("invalid worker count", func(t *testing.T) {
		pool := NewPool(0)
		defer pool.Close()

		var executed bool
		pool.Submit(func() { executed = true })
		pool.Wait()

		if !executed {
			t.Error("Pool with 0 workers should default to 1 worker")
		}
	})
}

func TestDebounce(t *testing.T) {
	t.Run("debounces rapid calls", func(t *testing.T) {
		var counter int64
		debounced := Debounce(func() {
			atomic.AddInt64(&counter, 1)
		}, 50*time.Millisecond)

		// Rapid calls - only the last one should execute
		for i := 0; i < 5; i++ {
			debounced()
			time.Sleep(10 * time.Millisecond)
		}

		// Wait for debounce period
		time.Sleep(100 * time.Millisecond)

		finalCount := atomic.LoadInt64(&counter)
		if finalCount != 1 {
			t.Errorf("Debounce executed %d times, want 1", finalCount)
		}
	})

	t.Run("executes after delay", func(t *testing.T) {
		var executed int32
		debounced := Debounce(func() {
			atomic.StoreInt32(&executed, 1)
		}, 30*time.Millisecond)

		debounced()

		// Check immediately - should not be executed yet
		if atomic.LoadInt32(&executed) != 0 {
			t.Error("Debounced function executed too early")
		}

		// Wait for delay and check again
		time.Sleep(50 * time.Millisecond)
		if atomic.LoadInt32(&executed) != 1 {
			t.Error("Debounced function was not executed after delay")
		}
	})
}

func TestThrottle(t *testing.T) {
	t.Run("throttles rapid calls", func(t *testing.T) {
		var counter int64
		throttled := Throttle(func() {
			atomic.AddInt64(&counter, 1)
		}, 50*time.Millisecond)

		// Rapid calls
		for i := 0; i < 5; i++ {
			throttled()
			time.Sleep(10 * time.Millisecond)
		}

		// First call should execute immediately, others should be throttled
		firstCount := atomic.LoadInt64(&counter)
		if firstCount != 1 {
			t.Errorf("Throttle executed %d times, want 1", firstCount)
		}

		// Wait for throttle period and call again
		time.Sleep(60 * time.Millisecond)
		throttled()

		finalCount := atomic.LoadInt64(&counter)
		if finalCount != 2 {
			t.Errorf("Throttle executed %d times after delay, want 2", finalCount)
		}
	})

	t.Run("first call executes immediately", func(t *testing.T) {
		var executed int32
		throttled := Throttle(func() {
			atomic.StoreInt32(&executed, 1)
		}, 100*time.Millisecond)

		throttled()

		if atomic.LoadInt32(&executed) != 1 {
			t.Error("First throttled call should execute immediately")
		}
	})
}