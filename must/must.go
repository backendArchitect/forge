// Package must provides utilities for panicking on errors, useful for initialization
// and situations where errors are not expected and should cause immediate failure.
package must

// Must takes a value and an error, and returns the value if the error is nil.
// If the error is not nil, it panics with the error.
// This is useful for initialization where errors are not expected.
//
// Example:
//
//	config := must.Must(loadConfig())
//	port := must.Must(strconv.Atoi(os.Getenv("PORT")))
//	fmt.Printf("Server starting on port %d with config %+v\n", port, config)
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Must0 takes only an error and panics if it is not nil.
// This is useful for operations that return only an error.
//
// Example:
//
//	file, err := os.Open("config.txt")
//	must.Must0(err)
//	defer file.Close()
//
//	must.Must0(json.NewEncoder(file).Encode(data))
func Must0(err error) {
	if err != nil {
		panic(err)
	}
}
