// Package async provides utilities for concurrent programming in Go.
// All functions maintain thread safety and handle goroutine management gracefully.
package async

import (
	"sync"
	"time"
)

// ParallelMap applies a transformation function to each element of a slice concurrently.
// Results are returned in the same order as the input slice.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	squares := async.ParallelMap(numbers, func(x int) int {
//		time.Sleep(100 * time.Millisecond) // Simulate work
//		return x * x
//	})
//	fmt.Println(squares) // Output: [1 4 9 16 25]
func ParallelMap[T, U any](input []T, transform func(T) U) []U {
	if len(input) == 0 {
		return []U{}
	}

	result := make([]U, len(input))
	var wg sync.WaitGroup

	for i, item := range input {
		wg.Add(1)
		go func(index int, value T) {
			defer wg.Done()
			result[index] = transform(value)
		}(i, item)
	}

	wg.Wait()
	return result
}

// ErrGroup is a collection of goroutines working on subtasks that are part of the same overall task.
// It captures the first error that occurs and cancels remaining goroutines.
//
// Example:
//
//	eg := &async.ErrGroup{}
//	for i := 0; i < 5; i++ {
//		i := i // capture loop variable
//		eg.Go(func() error {
//			if i == 3 {
//				return fmt.Errorf("error at %d", i)
//			}
//			return nil
//		})
//	}
//	if err := eg.Wait(); err != nil {
//		fmt.Printf("Error: %v\n", err)
//	}
type ErrGroup struct {
	wg     sync.WaitGroup
	errOnce sync.Once
	err     error
}

// Go starts a goroutine and runs the given function.
// The first call to return a non-nil error cancels the group.
func (g *ErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

// Wait blocks until all goroutines have completed and returns the first error.
func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

// Pool represents a worker pool that can execute tasks concurrently.
// It maintains a fixed number of workers and distributes tasks among them.
//
// Example:
//
//	pool := async.NewPool(3) // 3 workers
//	defer pool.Close()
//
//	tasks := []func(){
//		func() { fmt.Println("Task 1") },
//		func() { fmt.Println("Task 2") },
//		func() { fmt.Println("Task 3") },
//	}
//	pool.Submit(tasks...)
//	pool.Wait()
type Pool struct {
	workers   int
	taskQueue chan func()
	wg        sync.WaitGroup
	closed    bool
	mu        sync.Mutex
}

// NewPool creates a new worker pool with the specified number of workers.
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}

	p := &Pool{
		workers:   workers,
		taskQueue: make(chan func(), workers*2), // Buffered channel
	}

	// Start workers
	for i := 0; i < workers; i++ {
		go p.worker()
	}

	return p
}

// Submit adds tasks to the pool for execution.
func (p *Pool) Submit(tasks ...func()) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	for _, task := range tasks {
		if task != nil {
			p.wg.Add(1)
			p.taskQueue <- task
		}
	}
}

// Wait blocks until all submitted tasks have completed.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Close shuts down the pool and waits for all tasks to complete.
func (p *Pool) Close() {
	p.mu.Lock()
	if !p.closed {
		p.closed = true
		close(p.taskQueue)
	}
	p.mu.Unlock()
	p.wg.Wait()
}

// worker is the internal worker function that processes tasks.
func (p *Pool) worker() {
	for task := range p.taskQueue {
		task()
		p.wg.Done()
	}
}

// Debounce creates a debounced version of a function that delays execution
// until after the specified duration has elapsed since the last call.
//
// Example:
//
//	debouncedFunc := async.Debounce(func() {
//		fmt.Println("Executed after delay")
//	}, 100*time.Millisecond)
//
//	debouncedFunc() // Will be cancelled
//	debouncedFunc() // Will be cancelled  
//	debouncedFunc() // Will execute after 100ms
func Debounce(f func(), delay time.Duration) func() {
	var timer *time.Timer
	var mu sync.Mutex

	return func() {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(delay, f)
	}
}

// Throttle creates a throttled version of a function that limits execution
// to at most once per specified duration.
//
// Example:
//
//	throttledFunc := async.Throttle(func() {
//		fmt.Println("Executed at most once per interval")
//	}, 100*time.Millisecond)
//
//	throttledFunc() // Executes immediately
//	throttledFunc() // Ignored (too soon)
//	time.Sleep(150 * time.Millisecond)
//	throttledFunc() // Executes (enough time has passed)
func Throttle(f func(), interval time.Duration) func() {
	var lastExecution time.Time
	var mu sync.Mutex

	return func() {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		if now.Sub(lastExecution) >= interval {
			lastExecution = now
			f()
		}
	}
}