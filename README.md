# forge - A Modular Go Utility Library

[![CI](https://github.com/backendArchitect/forge/actions/workflows/ci.yml/badge.svg)](https://github.com/backendArchitect/forge/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/backendArchitect/forge.svg)](https://pkg.go.dev/github.com/backendArchitect/forge)
[![Go Report Card](https://goreportcard.com/badge/github.com/backendArchitect/forge)](https://goreportcard.com/report/github.com/backendArchitect/forge)

forge is a modular Go utility library that provides common functionality with zero external dependencies. The library follows Go best practices and offers excellent performance, comprehensive documentation, and high test coverage.

## Features

- **Zero Dependencies**: Only uses the Go standard library
- **Modular Design**: Import only what you need
- **High Performance**: Optimized for efficiency
- **Comprehensive Testing**: High test coverage with edge case handling
- **Excellent Documentation**: Every exported function has clear documentation and examples
- **Type Safety**: Leverages Go generics for type-safe operations

## Installation

```bash
go get github.com/backendArchitect/forge
```

## Packages

### conv - Type Conversion Utilities

```go
import "github.com/backendArchitect/forge/conv"

// Convert various types to integers
age, err := conv.ToInt("25")
count := conv.Must(conv.ToInt("42")) // Using with must package

// Convert anything to string
str := conv.ToString(123)
str = conv.ToString(true)

// Convert to boolean with comprehensive support
isActive, err := conv.ToBool("true")  // true
enabled := conv.ToBool(1)             // true
disabled := conv.ToBool(0)            // false

// Convert to float64 with type safety
price, err := conv.ToFloat64("123.45")
ratio := conv.ToFloat64(42)           // 42.0

// JSON operations
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

jsonStr, _ := conv.ToJSON(Person{Name: "John", Age: 30})
person, _ := conv.FromJSON[Person](jsonStr)

// Type-safe slice conversion
numbers, _ := conv.ToSlice[int]([]int{1, 2, 3})
```

### sliceutil - Slice Operations

```go
import "github.com/backendArchitect/forge/sliceutil"

numbers := []int{1, 2, 2, 3, 1, 4}

// Remove duplicates
unique := sliceutil.Unique(numbers) // [1, 2, 3, 4]

// Filter elements
evens := sliceutil.Filter(numbers, func(x int) bool { return x%2 == 0 })

// Transform elements
strings := sliceutil.Map(numbers, func(x int) string { 
    return fmt.Sprintf("num_%d", x) 
})

// Reduce to single value
sum := sliceutil.Reduce(numbers, 0, func(acc, x int) int { return acc + x })

// Split into chunks
chunks := sliceutil.Chunk(numbers, 2) // [[1, 2], [2, 3], [1, 4]]

// Check containment
hasTwo := sliceutil.Contains(numbers, 2) // true

// Shuffle elements
shuffled := sliceutil.Shuffle(numbers)

// Find differences
diff := sliceutil.Difference([]int{1, 2, 3}, []int{2, 4}) // [1, 3]

// Find intersections
common := sliceutil.Intersection([]int{1, 2, 3}, []int{2, 3, 4}) // [2, 3]

// Reverse elements in place
reversed := sliceutil.Reverse([]int{1, 2, 3, 4}) // [4, 3, 2, 1]
```

### strutil - String Utilities

```go
import "github.com/backendArchitect/forge/strutil"

// Generate random strings
randomID := strutil.Random(8) // "aBc3DeF9"
customRandom := strutil.Random(5, "12345") // "31425"

// Check if string is blank
isEmpty := strutil.IsBlank("   ") // true

// Truncate strings
short := strutil.Truncate("Hello World", 5, "...") // "He..."

// Case conversions
snake := strutil.CamelToSnake("myFunctionName") // "my_function_name"
camel := strutil.SnakeToCamel("my_function_name") // "myFunctionName"

// Template substitution
template := "Hello {{name}}, you are {{age}} years old"
data := map[string]any{"name": "John", "age": 30}
result, _ := strutil.Template(template, data)
// "Hello John, you are 30 years old"

// Pad strings to specific length
padded := strutil.Pad("hello", 10, ' ')    // "hello     "
centered := strutil.Pad("hi", 6, '*')      // "hi****"

// Capitalize first letter
title := strutil.Capitalize("hello WORLD") // "Hello world"

// Reverse string characters
backwards := strutil.Reverse("hello")      // "olleh"
```

### async - Concurrency Utilities

```go
import "github.com/backendArchitect/forge/async"

// Parallel processing with order preservation
numbers := []int{1, 2, 3, 4, 5}
squares := async.ParallelMap(numbers, func(x int) int {
    time.Sleep(100 * time.Millisecond) // Simulate work
    return x * x
})

// Error group for managing goroutines
eg := &async.ErrGroup{}
for i := 0; i < 5; i++ {
    i := i
    eg.Go(func() error {
        // Do work that might return an error
        return processItem(i)
    })
}
if err := eg.Wait(); err != nil {
    log.Printf("Error: %v", err)
}

// Worker pool
pool := async.NewPool(3)
defer pool.Close()

tasks := []func(){
    func() { fmt.Println("Task 1") },
    func() { fmt.Println("Task 2") },
    func() { fmt.Println("Task 3") },
}
pool.Submit(tasks...)
pool.Wait()

// Debounce function calls
debouncedSave := async.Debounce(func() {
    fmt.Println("Saving...")
}, 300*time.Millisecond)

// Throttle function calls
throttledUpdate := async.Throttle(func() {
    fmt.Println("Updating...")
}, 100*time.Millisecond)

// Execute functions with timeout
err := async.Timeout(func() error {
    // Some operation that might take too long
    return doSomething()
}, 5*time.Second)

// Retry operations with exponential backoff
result, err := async.Retry(func() (string, error) {
    return fetchDataFromAPI()
}, 3, 100*time.Millisecond)
```

### must - Panic on Error

```go
import "github.com/backendArchitect/forge/must"

// Panic if error is not nil, otherwise return value
config := must.Must(loadConfig()) 
port := must.Must(strconv.Atoi(os.Getenv("PORT")))

// Panic if error is not nil (for functions that only return error)
file, err := os.Open("config.txt")
must.Must0(err)
defer file.Close()
```

### fsutil - File System Utilities

```go
import "github.com/backendArchitect/forge/fsutil"

// Check file/directory existence
if fsutil.Exists("/path/to/file") {
    fmt.Println("File exists")
}

if fsutil.IsFile("/path/to/file") {
    fmt.Println("It's a file")
}

if fsutil.IsDir("/path/to/directory") {
    fmt.Println("It's a directory")
}

// JSON file operations
type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

config := Config{Host: "localhost", Port: 8080}
must.Must0(fsutil.WriteJSON("config.json", config))

var loadedConfig Config
must.Must0(fsutil.ReadJSON("config.json", &loadedConfig))

// Copy files
must.Must0(fsutil.CopyFile("source.txt", "destination.txt"))

// Ensure directory exists
must.Must0(fsutil.EnsureDir("/path/to/nested/directory"))

// Read and write text files
content, err := fsutil.ReadFile("config.txt")
must.Must0(fsutil.WriteFile("output.txt", "Hello, World!"))
```

## Contributing

Contributions are welcome! Please ensure that:

1. Code follows Go conventions and best practices
2. All functions have comprehensive documentation with examples
3. Tests cover happy paths, edge cases, and error conditions
4. No external dependencies are introduced
5. CI pipeline passes (linting, testing, and building)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
