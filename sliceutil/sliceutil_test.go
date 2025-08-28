package sliceutil

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"with duplicates", []int{1, 2, 2, 3, 1, 4}, []int{1, 2, 3, 4}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty slice", []int{}, []int{}},
		{"all same", []int{5, 5, 5}, []int{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		isEven := func(x int) bool { return x%2 == 0 }
		got := Filter(input, isEven)
		want := []int{2, 4, 6}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	})

	t.Run("filter keeps all", func(t *testing.T) {
		input := []int{1, 2, 3}
		alwaysTrue := func(x int) bool { return true }
		got := Filter(input, alwaysTrue)
		if !reflect.DeepEqual(got, input) {
			t.Errorf("Filter() = %v, want %v", got, input)
		}
	})

	t.Run("filter keeps none", func(t *testing.T) {
		input := []int{1, 2, 3}
		alwaysFalse := func(x int) bool { return false }
		got := Filter(input, alwaysFalse)
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		isEven := func(x int) bool { return x%2 == 0 }
		got := Filter(input, isEven)
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Filter() = %v, want %v", got, want)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		input := []int{1, 2, 3}
		toString := func(x int) string { return fmt.Sprintf("num_%d", x) }
		got := Map(input, toString)
		want := []string{"num_1", "num_2", "num_3"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map() = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		toString := func(x int) string { return fmt.Sprintf("%d", x) }
		got := Map(input, toString)
		want := []string{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map() = %v, want %v", got, want)
		}
	})

	t.Run("double numbers", func(t *testing.T) {
		input := []int{1, 2, 3}
		double := func(x int) int { return x * 2 }
		got := Map(input, double)
		want := []int{2, 4, 6}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Map() = %v, want %v", got, want)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("sum reduction", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		sum := func(acc int, x int) int { return acc + x }
		got := Reduce(input, 0, sum)
		want := 10
		if got != want {
			t.Errorf("Reduce() = %v, want %v", got, want)
		}
	})

	t.Run("string concatenation", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		concat := func(acc string, x string) string { return acc + x }
		got := Reduce(input, "", concat)
		want := "abc"
		if got != want {
			t.Errorf("Reduce() = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		sum := func(acc int, x int) int { return acc + x }
		got := Reduce(input, 10, sum)
		want := 10
		if got != want {
			t.Errorf("Reduce() = %v, want %v", got, want)
		}
	})
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		chunkSize int
		want      [][]int
	}{
		{"perfectly divisible", []int{1, 2, 3, 4, 5, 6}, 3, [][]int{{1, 2, 3}, {4, 5, 6}}},
		{"not perfectly divisible", []int{1, 2, 3, 4, 5, 6, 7}, 3, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}},
		{"chunk larger than slice", []int{1, 2}, 5, [][]int{{1, 2}}},
		{"empty slice", []int{}, 3, [][]int{}},
		{"invalid chunk size", []int{1, 2, 3}, 0, [][]int{}},
		{"negative chunk size", []int{1, 2, 3}, -1, [][]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chunk(tt.input, tt.chunkSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name    string
		slice   []string
		element string
		want    bool
	}{
		{"element exists", []string{"apple", "banana", "cherry"}, "banana", true},
		{"element does not exist", []string{"apple", "banana", "cherry"}, "grape", false},
		{"empty slice", []string{}, "apple", false},
		{"single element - exists", []string{"apple"}, "apple", true},
		{"single element - not exists", []string{"apple"}, "banana", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.slice, tt.element)
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	t.Run("shuffle maintains elements", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		input := make([]int, len(original))
		copy(input, original)

		shuffled := Shuffle(input)

		// Check same length
		if len(shuffled) != len(original) {
			t.Errorf("Shuffle() changed length: got %d, want %d", len(shuffled), len(original))
		}

		// Check all elements are present
		sort.Ints(shuffled)
		sort.Ints(original)
		if !reflect.DeepEqual(shuffled, original) {
			t.Errorf("Shuffle() changed elements: got %v, want %v", shuffled, original)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		shuffled := Shuffle(input)
		if len(shuffled) != 0 {
			t.Errorf("Shuffle() of empty slice should remain empty")
		}
	})

	t.Run("single element", func(t *testing.T) {
		input := []int{42}
		shuffled := Shuffle(input)
		if len(shuffled) != 1 || shuffled[0] != 42 {
			t.Errorf("Shuffle() of single element should remain unchanged")
		}
	})
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name string
		s1   []int
		s2   []int
		want []int
	}{
		{"overlapping slices", []int{1, 2, 3, 4}, []int{2, 4, 6}, []int{1, 3}},
		{"non-overlapping slices", []int{1, 2, 3}, []int{4, 5, 6}, []int{1, 2, 3}},
		{"empty first slice", []int{}, []int{1, 2, 3}, []int{}},
		{"empty second slice", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{"both empty", []int{}, []int{}, []int{}},
		{"identical slices", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Difference(tt.s1, tt.s2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}