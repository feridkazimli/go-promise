package promise

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Unit Tests

func TestNew(t *testing.T) {
	p := New(func() (int, error) {
		return 42, nil
	})

	result, err := p.Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestNewWithError(t *testing.T) {
	expectedErr := errors.New("test error")
	p := New(func() (int, error) {
		return 0, expectedErr
	})

	_, err := p.Await()
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestThen(t *testing.T) {
	p := New(func() (int, error) {
		return 10, nil
	}).Then(func(val int) (int, error) {
		return val * 2, nil
	}).Then(func(val int) (int, error) {
		return val + 5, nil
	})

	result, err := p.Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 25 {
		t.Errorf("Expected 25, got %d", result)
	}
}

func TestThenWithError(t *testing.T) {
	expectedErr := errors.New("chain error")
	p := New(func() (int, error) {
		return 10, nil
	}).Then(func(val int) (int, error) {
		return 0, expectedErr
	}).Then(func(val int) (int, error) {
		return val + 5, nil
	})

	_, err := p.Await()
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestCatch(t *testing.T) {
	p := New(func() (int, error) {
		return 0, errors.New("initial error")
	}).Catch(func(err error) (int, error) {
		return 999, nil
	})

	result, err := p.Await()
	if err != nil {
		t.Errorf("Expected no error after catch, got %v", err)
	}
	if result != 999 {
		t.Errorf("Expected 999, got %d", result)
	}
}

func TestFinally(t *testing.T) {
	finallyCalled := false
	p := New(func() (int, error) {
		return 42, nil
	}).Finally(func() {
		finallyCalled = true
	})

	result, _ := p.Await()
	if !finallyCalled {
		t.Error("Finally was not called")
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestFinallyWithError(t *testing.T) {
	finallyCalled := false
	p := New(func() (int, error) {
		return 0, errors.New("error")
	}).Finally(func() {
		finallyCalled = true
	})

	_, err := p.Await()
	if !finallyCalled {
		t.Error("Finally was not called on error")
	}
	if err == nil {
		t.Error("Expected error to propagate")
	}
}

func TestAll(t *testing.T) {
	promises := []*Promise[int]{
		New(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 1, nil
		}),
		New(func() (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 2, nil
		}),
		New(func() (int, error) {
			time.Sleep(20 * time.Millisecond)
			return 3, nil
		}),
	}

	result, err := All(promises...).Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}
	expected := []int{1, 2, 3}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected result[%d] = %d, got %d", i, v, result[i])
		}
	}
}

func TestAllWithError(t *testing.T) {
	promises := []*Promise[int]{
		New(func() (int, error) {
			return 1, nil
		}),
		New(func() (int, error) {
			return 0, errors.New("error in promise 2")
		}),
		New(func() (int, error) {
			return 3, nil
		}),
	}

	_, err := All(promises...).Await()
	if err == nil {
		t.Error("Expected error from All")
	}
}

func TestAllSettled(t *testing.T) {
	promises := []*Promise[int]{
		New(func() (int, error) {
			return 1, nil
		}),
		New(func() (int, error) {
			return 0, errors.New("error")
		}),
		New(func() (int, error) {
			return 3, nil
		}),
	}

	results, err := AllSettled(promises...).Await()
	if err != nil {
		t.Errorf("AllSettled should not return error, got %v", err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
	if results[0].Error != nil {
		t.Error("First promise should succeed")
	}
	if results[1].Error == nil {
		t.Error("Second promise should fail")
	}
	if results[2].Error != nil {
		t.Error("Third promise should succeed")
	}
}

func TestAny(t *testing.T) {
	promises := []*Promise[int]{
		New(func() (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		}),
		New(func() (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 2, nil
		}),
		New(func() (int, error) {
			time.Sleep(50 * time.Millisecond)
			return 3, nil
		}),
	}

	result, err := Any(promises...).Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 2 {
		t.Errorf("Expected 2 (fastest), got %d", result)
	}
}

func TestAnyAllFail(t *testing.T) {
	promises := []*Promise[int]{
		New(func() (int, error) {
			return 0, errors.New("error 1")
		}),
		New(func() (int, error) {
			return 0, errors.New("error 2")
		}),
	}

	_, err := Any(promises...).Await()
	if err == nil {
		t.Error("Expected error when all promises fail")
	}
}

func TestRace(t *testing.T) {
	promises := []*Promise[int]{
		New(func() (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		}),
		New(func() (int, error) {
			time.Sleep(30 * time.Millisecond)
			return 2, nil
		}),
	}

	result, err := Race(promises...).Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 2 {
		t.Errorf("Expected 2 (fastest), got %d", result)
	}
}

func TestResolve(t *testing.T) {
	p := Resolve(42)
	result, err := p.Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestReject(t *testing.T) {
	expectedErr := errors.New("rejection error")
	p := Reject[int](expectedErr)
	_, err := p.Await()
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestWithTimeout(t *testing.T) {
	p := New(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	}).WithTimeout(50 * time.Millisecond)

	_, err := p.Await()
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestWithTimeoutSuccess(t *testing.T) {
	p := New(func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 42, nil
	}).WithTimeout(100 * time.Millisecond)

	result, err := p.Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	p := New(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	}).WithContext(ctx)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := p.Await()
	if err == nil {
		t.Error("Expected context cancellation error")
	}
}

func TestMap(t *testing.T) {
	p := New(func() (int, error) {
		return 10, nil
	})

	mapped := Map(p, func(val int) string {
		return "value: " + string(rune(val+'0'))
	})

	result, err := mapped.Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "value: :" {
		t.Errorf("Expected mapped value, got %s", result)
	}
}

func TestFlatMap(t *testing.T) {
	p := New(func() (int, error) {
		return 10, nil
	})

	flattened := FlatMap(p, func(val int) *Promise[int] {
		return New(func() (int, error) {
			return val * 2, nil
		})
	})

	result, err := flattened.Await()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 20 {
		t.Errorf("Expected 20, got %d", result)
	}
}

// Benchmarks

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := New(func() (int, error) {
			return 42, nil
		})
		p.Await()
	}
}

func BenchmarkThenChain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := New(func() (int, error) {
			return 10, nil
		}).Then(func(val int) (int, error) {
			return val * 2, nil
		}).Then(func(val int) (int, error) {
			return val + 5, nil
		})
		p.Await()
	}
}

func BenchmarkAll10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promises := make([]*Promise[int], 10)
		for j := 0; j < 10; j++ {
			val := j
			promises[j] = New(func() (int, error) {
				return val, nil
			})
		}
		All(promises...).Await()
	}
}

func BenchmarkAll100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promises := make([]*Promise[int], 100)
		for j := 0; j < 100; j++ {
			val := j
			promises[j] = New(func() (int, error) {
				return val, nil
			})
		}
		All(promises...).Await()
	}
}

func BenchmarkAny10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promises := make([]*Promise[int], 10)
		for j := 0; j < 10; j++ {
			val := j
			promises[j] = New(func() (int, error) {
				return val, nil
			})
		}
		Any(promises...).Await()
	}
}

func BenchmarkRace10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		promises := make([]*Promise[int], 10)
		for j := 0; j < 10; j++ {
			val := j
			promises[j] = New(func() (int, error) {
				return val, nil
			})
		}
		Race(promises...).Await()
	}
}

func BenchmarkWithTimeout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := New(func() (int, error) {
			time.Sleep(1 * time.Millisecond)
			return 42, nil
		}).WithTimeout(100 * time.Millisecond)
		p.Await()
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := New(func() (int, error) {
			return 10, nil
		})
		mapped := Map(p, func(val int) int {
			return val * 2
		})
		mapped.Await()
	}
}

func BenchmarkFlatMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := New(func() (int, error) {
			return 10, nil
		})
		flattened := FlatMap(p, func(val int) *Promise[int] {
			return New(func() (int, error) {
				return val * 2, nil
			})
		})
		flattened.Await()
	}
}

func BenchmarkConcurrentPromises(b *testing.B) {
	b.Run("10_concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			promises := make([]*Promise[int], 10)
			for j := 0; j < 10; j++ {
				val := j
				promises[j] = New(func() (int, error) {
					time.Sleep(1 * time.Millisecond)
					return val, nil
				})
			}
			All(promises...).Await()
		}
	})

	b.Run("100_concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			promises := make([]*Promise[int], 100)
			for j := 0; j < 100; j++ {
				val := j
				promises[j] = New(func() (int, error) {
					time.Sleep(1 * time.Millisecond)
					return val, nil
				})
			}
			All(promises...).Await()
		}
	})

	b.Run("1000_concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			promises := make([]*Promise[int], 1000)
			for j := 0; j < 1000; j++ {
				val := j
				promises[j] = New(func() (int, error) {
					return val, nil
				})
			}
			All(promises...).Await()
		}
	})
}

func BenchmarkResolveReject(b *testing.B) {
	b.Run("Resolve", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := Resolve(42)
			p.Await()
		}
	})

	b.Run("Reject", func(b *testing.B) {
		err := errors.New("test error")
		for i := 0; i < b.N; i++ {
			p := Reject[int](err)
			p.Await()
		}
	})
}
