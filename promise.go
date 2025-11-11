package promise

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Promise pattern implementation for managing goroutines
type Promise[T any] struct {
	wg     sync.WaitGroup
	result T
	err    error
	done   chan struct{}
	once   sync.Once
}

// New creates a new Promise
func New[T any](fn func() (T, error)) *Promise[T] {
	p := &Promise[T]{
		done: make(chan struct{}),
	}
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		defer p.once.Do(func() { close(p.done) })
		p.result, p.err = fn()
	}()

	return p
}

// Await waits for the Promise to complete and returns the result
func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()
	return p.result, p.err
}

// Then executes the callback when the Promise is fulfilled
func (p *Promise[T]) Then(fn func(T) (T, error)) *Promise[T] {
	return New(func() (T, error) {
		result, err := p.Await()
		if err != nil {
			var zero T
			return zero, err
		}
		return fn(result)
	})
}

// Catch is used for error handling
func (p *Promise[T]) Catch(fn func(error) (T, error)) *Promise[T] {
	return New(func() (T, error) {
		result, err := p.Await()
		if err != nil {
			return fn(err)
		}
		return result, nil
	})
}

// Finally executes after the Promise is settled (either fulfilled or rejected)
func (p *Promise[T]) Finally(fn func()) *Promise[T] {
	return New(func() (T, error) {
		defer fn()
		return p.Await()
	})
}

// All waits for all Promises to complete (similar to Promise.all)
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	return New(func() ([]T, error) {
		results := make([]T, len(promises))
		for i, p := range promises {
			result, err := p.Await()
			if err != nil {
				return nil, err
			}
			results[i] = result
		}
		return results, nil
	})
}

// AllSettled waits for all Promises to settle (either fulfilled or rejected)
type SettledResult[T any] struct {
	Value T
	Error error
}

func AllSettled[T any](promises ...*Promise[T]) *Promise[[]SettledResult[T]] {
	return New(func() ([]SettledResult[T], error) {
		results := make([]SettledResult[T], len(promises))
		for i, p := range promises {
			value, err := p.Await()
			results[i] = SettledResult[T]{
				Value: value,
				Error: err,
			}
		}
		return results, nil
	})
}

// Any returns the first fulfilled Promise (similar to Promise.any)
func Any[T any](promises ...*Promise[T]) *Promise[T] {
	return New(func() (T, error) {
		if len(promises) == 0 {
			var zero T
			return zero, errors.New("no promises provided")
		}

		resultCh := make(chan T, len(promises))
		errorCh := make(chan error, len(promises))

		for _, p := range promises {
			go func(promise *Promise[T]) {
				result, err := promise.Await()
				if err != nil {
					errorCh <- err
				} else {
					resultCh <- result
				}
			}(p)
		}

		errorCount := 0
		for {
			select {
			case result := <-resultCh:
				return result, nil
			case <-errorCh:
				errorCount++
				if errorCount == len(promises) {
					var zero T
					return zero, errors.New("all promises rejected")
				}
			}
		}
	})
}

// Race returns the first settled Promise (similar to Promise.race)
func Race[T any](promises ...*Promise[T]) *Promise[T] {
	return New(func() (T, error) {
		if len(promises) == 0 {
			var zero T
			return zero, errors.New("no promises provided")
		}

		resultCh := make(chan T, 1)
		errorCh := make(chan error, 1)

		for _, p := range promises {
			go func(promise *Promise[T]) {
				result, err := promise.Await()
				if err != nil {
					select {
					case errorCh <- err:
					default:
					}
				} else {
					select {
					case resultCh <- result:
					default:
					}
				}
			}(p)
		}

		select {
		case result := <-resultCh:
			return result, nil
		case err := <-errorCh:
			var zero T
			return zero, err
		}
	})
}

// Resolve creates a fulfilled Promise
func Resolve[T any](value T) *Promise[T] {
	return New(func() (T, error) {
		return value, nil
	})
}

// Reject creates a rejected Promise
func Reject[T any](err error) *Promise[T] {
	return New(func() (T, error) {
		var zero T
		return zero, err
	})
}

// WithTimeout adds a timeout to the Promise
func (p *Promise[T]) WithTimeout(timeout time.Duration) *Promise[T] {
	return New(func() (T, error) {
		resultCh := make(chan struct {
			value T
			err   error
		}, 1)

		go func() {
			value, err := p.Await()
			resultCh <- struct {
				value T
				err   error
			}{value, err}
		}()

		select {
		case result := <-resultCh:
			return result.value, result.err
		case <-time.After(timeout):
			var zero T
			return zero, errors.New("promise timeout")
		}
	})
}

// WithContext creates a Promise that can be canceled using a context
func (p *Promise[T]) WithContext(ctx context.Context) *Promise[T] {
	return New(func() (T, error) {
		resultCh := make(chan struct {
			value T
			err   error
		}, 1)

		go func() {
			value, err := p.Await()
			resultCh <- struct {
				value T
				err   error
			}{value, err}
		}()

		select {
		case result := <-resultCh:
			return result.value, result.err
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	})
}

// Map transforms the result of a Promise
func Map[T, U any](p *Promise[T], fn func(T) U) *Promise[U] {
	return New(func() (U, error) {
		result, err := p.Await()
		if err != nil {
			var zero U
			return zero, err
		}
		return fn(result), nil
	})
}

// FlatMap is used for chaining Promises
func FlatMap[T, U any](p *Promise[T], fn func(T) *Promise[U]) *Promise[U] {
	return New(func() (U, error) {
		result, err := p.Await()
		if err != nil {
			var zero U
			return zero, err
		}
		return fn(result).Await()
	})
}
