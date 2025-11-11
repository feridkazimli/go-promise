# Go Promise 🚀

A JavaScript Promise-like library for Go, making goroutines more manageable and composable with a familiar API.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/feridkazimli/go-promise)](https://goreportcard.com/report/github.com/feridkazimli/go-promise)

## ✨ Features

- 🎯 **Familiar API** - JavaScript Promise-like interface for Go developers
- 🔄 **Composable** - Chain promises with `.Then()`, `.Catch()`, and `.Finally()`
- ⚡ **Concurrent** - Built-in support for `All()`, `Any()`, `Race()`, and `AllSettled()`
- 🛡️ **Type-Safe** - Full generic support for type safety
- ⏱️ **Timeout & Context** - Built-in timeout and context cancellation support
- 🔧 **Functional** - `Map()` and `FlatMap()` for functional composition
- 📦 **Zero Dependencies** - Only uses Go standard library

## 📦 Installation

```bash
go get github.com/feridkazimli/go-promise
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/feridkazimli/go-promise"
)

func main() {
    // Create a simple promise
    p := promise.New(func() (int, error) {
        time.Sleep(100 * time.Millisecond)
        return 42, nil
    })

    result, err := p.Await()
    fmt.Println("Result:", result) // Output: Result: 42
}
```

## 📖 Usage Examples

### Basic Promise

```go
// Create and await a promise
p := promise.New(func() (string, error) {
    return "Hello, World!", nil
})

result, err := p.Await()
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // Output: Hello, World!
```

### Chaining with Then()

```go
result, _ := promise.New(func() (int, error) {
    return 10, nil
}).Then(func(val int) (int, error) {
    return val * 2, nil
}).Then(func(val int) (int, error) {
    return val + 5, nil
}).Await()

fmt.Println(result) // Output: 25
```

### Error Handling with Catch()

```go
result, _ := promise.New(func() (int, error) {
    return 0, errors.New("something went wrong")
}).Catch(func(err error) (int, error) {
    log.Println("Error caught:", err)
    return 999, nil // Return default value
}).Await()

fmt.Println(result) // Output: 999
```

### Finally()

```go
promise.New(func() (int, error) {
    return 42, nil
}).Finally(func() {
    fmt.Println("Cleanup executed")
}).Await()
```

### Promise.All() - Wait for All Promises

```go
promises := []*promise.Promise[int]{
    promise.New(func() (int, error) {
        time.Sleep(100 * time.Millisecond)
        return 1, nil
    }),
    promise.New(func() (int, error) {
        time.Sleep(50 * time.Millisecond)
        return 2, nil
    }),
    promise.New(func() (int, error) {
        time.Sleep(150 * time.Millisecond)
        return 3, nil
    }),
}

results, err := promise.All(promises...).Await()
fmt.Println(results) // Output: [1 2 3]
```

### Promise.Any() - First Successful Promise

```go
result, err := promise.Any(
    promise.New(func() (int, error) {
        time.Sleep(100 * time.Millisecond)
        return 1, nil
    }),
    promise.New(func() (int, error) {
        time.Sleep(30 * time.Millisecond)
        return 2, nil
    }),
    promise.New(func() (int, error) {
        time.Sleep(50 * time.Millisecond)
        return 3, nil
    }),
).Await()

fmt.Println(result) // Output: 2 (fastest to complete)
```

### Promise.Race() - First Completed Promise

```go
result, err := promise.Race(promises...).Await()
// Returns the first promise to complete (success or error)
```

### Promise.AllSettled() - All Results

```go
promises := []*promise.Promise[int]{
    promise.New(func() (int, error) { return 1, nil }),
    promise.New(func() (int, error) { return 0, errors.New("error") }),
    promise.New(func() (int, error) { return 3, nil }),
}

results, _ := promise.AllSettled(promises...).Await()
for i, result := range results {
    fmt.Printf("Promise %d: Value=%d, Error=%v\n", i, result.Value, result.Error)
}
```

### Timeout Support

```go
p := promise.New(func() (int, error) {
    time.Sleep(2 * time.Second)
    return 42, nil
}).WithTimeout(500 * time.Millisecond)

_, err := p.Await()
if err != nil {
    fmt.Println("Timeout!") // Will timeout
}
```

### Context Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

p := promise.New(func() (int, error) {
    time.Sleep(2 * time.Second)
    return 42, nil
}).WithContext(ctx)

// Cancel after 100ms
go func() {
    time.Sleep(100 * time.Millisecond)
    cancel()
}()

_, err := p.Await()
if err != nil {
    fmt.Println("Cancelled!") // Will be cancelled
}
```

### Map and FlatMap

```go
// Map - transform promise value
stringPromise := promise.Map(
    promise.Resolve(42),
    func(val int) string {
        return fmt.Sprintf("Value: %d", val)
    },
)

// FlatMap - chain promises
result := promise.FlatMap(
    promise.Resolve(10),
    func(val int) *promise.Promise[int] {
        return promise.New(func() (int, error) {
            return val * 2, nil
        })
    },
)
```

### Resolve and Reject

```go
// Create resolved promise
p1 := promise.Resolve(42)

// Create rejected promise
p2 := promise.Reject[int](errors.New("error"))
```

## 🎯 API Reference

### Core Methods

| Method        | Description                                     |
| ------------- | ----------------------------------------------- |
| `New(fn)`     | Creates a new promise                           |
| `Await()`     | Waits for promise completion and returns result |
| `Then(fn)`    | Chains transformation on success                |
| `Catch(fn)`   | Handles errors                                  |
| `Finally(fn)` | Executes cleanup regardless of outcome          |

### Static Methods

| Method                    | Description                                   |
| ------------------------- | --------------------------------------------- |
| `All(promises...)`        | Waits for all promises (fails on first error) |
| `Any(promises...)`        | Returns first successful promise              |
| `Race(promises...)`       | Returns first completed promise               |
| `AllSettled(promises...)` | Waits for all promises (collects all results) |
| `Resolve(value)`          | Creates resolved promise                      |
| `Reject(error)`           | Creates rejected promise                      |

### Advanced Methods

| Method                  | Description                       |
| ----------------------- | --------------------------------- |
| `WithTimeout(duration)` | Adds timeout to promise           |
| `WithContext(ctx)`      | Adds context cancellation support |
| `Map(promise, fn)`      | Transforms promise value          |
| `FlatMap(promise, fn)`  | Chains promises                   |

## 🧪 Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark

# Run all checks (fmt, vet, lint, test)
make check
```

## 📊 Benchmarks

```
BenchmarkNew-8                    500000      2843 ns/op      448 B/op      6 allocs/op
BenchmarkThenChain-8              200000      8234 ns/op     1344 B/op     18 allocs/op
BenchmarkAll100-8                   5000    285421 ns/op    44928 B/op    600 allocs/op
BenchmarkAny10-8                   50000     34567 ns/op     5376 B/op     71 allocs/op
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by JavaScript's Promise API
- Built with Go's powerful concurrency primitives

## 📬 Contact

Your Name - [@yourtwitter](https://twitter.com/yourtwitter)

Project Link: [https://github.com/feridkazimli/go-promise](https://github.com/feridkazimli/go-promise)

---

⭐ If you find this project useful, please consider giving it a star!

## 🗺️ Roadmap

- [ ] Promise pooling for resource management
- [ ] Retry mechanism with exponential backoff
- [ ] Progress tracking for long-running operations
- [ ] Promise cancellation tokens
- [ ] Integration with popular Go frameworks

## 💡 Why Use Go Promise?

### Before (Traditional Go)

```go
var wg sync.WaitGroup
results := make([]int, 3)
errors := make([]error, 3)

wg.Add(3)
go func() {
    defer wg.Done()
    results[0], errors[0] = fetchData1()
}()
go func() {
    defer wg.Done()
    results[1], errors[1] = fetchData2()
}()
go func() {
    defer wg.Done()
    results[2], errors[2] = fetchData3()
}()
wg.Wait()

// Check for errors...
```

### After (With Go Promise)

```go
results, err := promise.All(
    promise.New(fetchData1),
    promise.New(fetchData2),
    promise.New(fetchData3),
).Await()

if err != nil {
    log.Fatal(err)
}
```

Clean, concise, and easier to reason about! 🎉
