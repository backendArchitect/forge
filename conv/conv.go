// Package conv provides type conversion utilities with comprehensive error handling.
// All functions follow Go conventions and handle edge cases gracefully.
package conv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ToInt converts any value to an integer with comprehensive type support.
// Supports int types, float types, string representations of numbers, and booleans.
//
// Example:
//
//	i, err := conv.ToInt("123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(i) // Output: 123
func ToInt(v any) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("cannot convert nil to int")
	}

	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		if strings.TrimSpace(val) == "" {
			return 0, fmt.Errorf("cannot convert empty string to int")
		}
		return strconv.Atoi(val)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// ToString converts any value to its string representation.
// Handles all basic types and byte slices appropriately.
//
// Example:
//
//	s := conv.ToString(123)
//	fmt.Println(s) // Output: "123"
func ToString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// ToSlice converts any value to a slice of the specified type T.
// Returns an error if the input is not a slice or if types don't match.
//
// Example:
//
//	slice, err := conv.ToSlice[int]([]int{1, 2, 3})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(slice) // Output: [1 2 3]
func ToSlice[T any](v any) ([]T, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot convert nil to slice")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value is not a slice, got %T", v)
	}

	// Check if we can directly assert to the target type
	if slice, ok := v.([]T); ok {
		return slice, nil
	}

	// Try to convert each element
	length := rv.Len()
	result := make([]T, length)
	var zero T
	expectedType := reflect.TypeOf(zero)

	for i := 0; i < length; i++ {
		elem := rv.Index(i).Interface()
		if !reflect.TypeOf(elem).AssignableTo(expectedType) {
			return nil, fmt.Errorf("element at index %d is not assignable to type %T", i, zero)
		}
		result[i] = elem.(T)
	}

	return result, nil
}

// ToJSON converts any value to its JSON string representation.
// Returns an error if the value cannot be marshaled to JSON.
//
// Example:
//
//	type Person struct {
//		Name string `json:"name"`
//		Age  int    `json:"age"`
//	}
//	p := Person{Name: "John", Age: 30}
//	jsonStr, err := conv.ToJSON(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(jsonStr) // Output: {"name":"John","age":30}
func ToJSON(v any) (string, error) {
	if v == nil {
		return "null", nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return string(data), nil
}

// FromJSON parses a JSON string and returns a value of the specified type T.
// Returns an error if the JSON is malformed or cannot be unmarshaled to type T.
//
// Example:
//
//	type Person struct {
//		Name string `json:"name"`
//		Age  int    `json:"age"`
//	}
//	jsonStr := `{"name":"John","age":30}`
//	person, err := conv.FromJSON[Person](jsonStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%+v\n", person) // Output: {Name:John Age:30}
func FromJSON[T any](jsonStr string) (T, error) {
	var result T

	if strings.TrimSpace(jsonStr) == "" {
		return result, fmt.Errorf("cannot unmarshal empty JSON string")
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}
