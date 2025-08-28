package must

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	t.Run("returns value when no error", func(t *testing.T) {
		value := 42
		result := Must(value, nil)
		if result != value {
			t.Errorf("Must() = %v, want %v", result, value)
		}
	})

	t.Run("panics when error is not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Must() did not panic when error is not nil")
			}
		}()

		Must(42, errors.New("test error"))
	})

	t.Run("works with different types", func(t *testing.T) {
		stringResult := Must("hello", nil)
		if stringResult != "hello" {
			t.Errorf("Must() string = %v, want %v", stringResult, "hello")
		}

		sliceResult := Must([]int{1, 2, 3}, nil)
		if len(sliceResult) != 3 || sliceResult[0] != 1 {
			t.Errorf("Must() slice = %v, want [1 2 3]", sliceResult)
		}
	})

	t.Run("preserves error type in panic", func(t *testing.T) {
		testErr := errors.New("specific error")

		defer func() {
			if r := recover(); r != testErr {
				t.Errorf("Must() panic = %v, want %v", r, testErr)
			}
		}()

		Must(42, testErr)
	})
}

func TestMust0(t *testing.T) {
	t.Run("does not panic when error is nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Must0() panicked when error is nil: %v", r)
			}
		}()

		Must0(nil)
	})

	t.Run("panics when error is not nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Must0() did not panic when error is not nil")
			}
		}()

		Must0(errors.New("test error"))
	})

	t.Run("preserves error type in panic", func(t *testing.T) {
		testErr := errors.New("specific error")

		defer func() {
			if r := recover(); r != testErr {
				t.Errorf("Must0() panic = %v, want %v", r, testErr)
			}
		}()

		Must0(testErr)
	})
}
