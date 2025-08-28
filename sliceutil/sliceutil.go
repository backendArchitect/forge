// Package sliceutil provides utilities for working with slices in Go.
// All functions are generic and work with any slice type, maintaining type safety.
package sliceutil

import (
	"math/rand"
)

// Unique returns a new slice containing only unique elements from the input slice.
// The order of first occurrence is preserved.
//
// Example:
//
//	slice := []int{1, 2, 2, 3, 1, 4}
//	unique := sliceutil.Unique(slice)
//	fmt.Println(unique) // Output: [1 2 3 4]
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate function.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5}
//	evens := sliceutil.Filter(slice, func(x int) bool { return x%2 == 0 })
//	fmt.Println(evens) // Output: [2 4]
func Filter[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

// Map applies a transformation function to each element and returns a new slice.
//
// Example:
//
//	slice := []int{1, 2, 3}
//	strings := sliceutil.Map(slice, func(x int) string { return fmt.Sprintf("num_%d", x) })
//	fmt.Println(strings) // Output: [num_1 num_2 num_3]
func Map[T, U any](slice []T, transform func(T) U) []U {
	if len(slice) == 0 {
		return []U{}
	}

	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = transform(item)
	}

	return result
}

// Reduce applies a reduction function to the slice elements and returns a single value.
//
// Example:
//
//	slice := []int{1, 2, 3, 4}
//	sum := sliceutil.Reduce(slice, 0, func(acc int, x int) int { return acc + x })
//	fmt.Println(sum) // Output: 10
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// Chunk splits a slice into chunks of the specified size.
// The last chunk may be smaller if the slice length is not evenly divisible.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5, 6, 7}
//	chunks := sliceutil.Chunk(slice, 3)
//	fmt.Println(chunks) // Output: [[1 2 3] [4 5 6] [7]]
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	if len(slice) == 0 {
		return [][]T{}
	}

	var result [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}

	return result
}

// Contains checks if a slice contains a specific element.
//
// Example:
//
//	slice := []string{"apple", "banana", "cherry"}
//	found := sliceutil.Contains(slice, "banana")
//	fmt.Println(found) // Output: true
func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Shuffle randomly reorders the elements in a slice.
// This function modifies the original slice and returns it for convenience.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5}
//	shuffled := sliceutil.Shuffle(slice)
//	fmt.Println(shuffled) // Output: [3 1 5 2 4] (random order)
func Shuffle[T any](slice []T) []T {
	if len(slice) <= 1 {
		return slice
	}

	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}

// Difference returns elements that are in the first slice but not in the second slice.
//
// Example:
//
//	s1 := []int{1, 2, 3, 4}
//	s2 := []int{2, 4, 6}
//	diff := sliceutil.Difference(s1, s2)
//	fmt.Println(diff) // Output: [1 3]
func Difference[T comparable](s1, s2 []T) []T {
	if len(s1) == 0 {
		return []T{}
	}

	// Create a set from s2 for O(1) lookup
	s2Set := make(map[T]struct{})
	for _, item := range s2 {
		s2Set[item] = struct{}{}
	}

	result := make([]T, 0, len(s1))
	for _, item := range s1 {
		if _, exists := s2Set[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}
