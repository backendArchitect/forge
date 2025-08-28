package conv

import (
	"testing"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantErr bool
	}{
		// Integer types
		{"int", 42, 42, false},
		{"int8", int8(42), 42, false},
		{"int16", int16(42), 42, false},
		{"int32", int32(42), 42, false},
		{"int64", int64(42), 42, false},
		{"uint", uint(42), 42, false},
		{"uint8", uint8(42), 42, false},
		{"uint16", uint16(42), 42, false},
		{"uint32", uint32(42), 42, false},
		{"uint64", uint64(42), 42, false},

		// Float types
		{"float32", float32(42.5), 42, false},
		{"float64", 42.9, 42, false},

		// String types
		{"valid string", "123", 123, false},
		{"negative string", "-123", -123, false},
		{"invalid string", "abc", 0, true},
		{"empty string", "", 0, true},
		{"whitespace string", "   ", 0, true},

		// Boolean types
		{"true bool", true, 1, false},
		{"false bool", false, 0, false},

		// Error cases
		{"nil", nil, 0, true},
		{"unsupported type", []int{1, 2, 3}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"int", 123, "123"},
		{"float", 123.45, "123.45"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"byte slice", []byte("hello"), "hello"},
		{"empty byte slice", []byte{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.input)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	t.Run("valid int slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		got, err := ToSlice[int](input)
		if err != nil {
			t.Fatalf("ToSlice() error = %v", err)
		}
		if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Errorf("ToSlice() = %v, want [1 2 3]", got)
		}
	})

	t.Run("different type slice", func(t *testing.T) {
		input := []string{"1", "2", "3"}
		_, err := ToSlice[int](input)
		if err == nil {
			t.Error("ToSlice() expected error for different types")
		}
	})

	t.Run("non-slice type", func(t *testing.T) {
		_, err := ToSlice[int](123)
		if err == nil {
			t.Error("ToSlice() expected error for non-slice type")
		}
	})

	t.Run("nil input", func(t *testing.T) {
		_, err := ToSlice[int](nil)
		if err == nil {
			t.Error("ToSlice() expected error for nil input")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got, err := ToSlice[int](input)
		if err != nil {
			t.Fatalf("ToSlice() error = %v", err)
		}
		if len(got) != 0 {
			t.Errorf("ToSlice() = %v, want empty slice", got)
		}
	})
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{"nil", nil, "null", false},
		{"string", "hello", `"hello"`, false},
		{"int", 123, "123", false},
		{"struct", struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{"John", 30}, `{"name":"John","age":30}`, false},
		{"map", map[string]int{"a": 1, "b": 2}, `{"a":1,"b":2}`, false},
		{"unmarshalable", make(chan int), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	t.Run("valid JSON to struct", func(t *testing.T) {
		type Person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		jsonStr := `{"name":"John","age":30}`
		got, err := FromJSON[Person](jsonStr)
		if err != nil {
			t.Fatalf("FromJSON() error = %v", err)
		}
		if got.Name != "John" || got.Age != 30 {
			t.Errorf("FromJSON() = %v, want {Name:John Age:30}", got)
		}
	})

	t.Run("valid JSON to map", func(t *testing.T) {
		jsonStr := `{"a":1,"b":2}`
		got, err := FromJSON[map[string]int](jsonStr)
		if err != nil {
			t.Fatalf("FromJSON() error = %v", err)
		}
		if got["a"] != 1 || got["b"] != 2 {
			t.Errorf("FromJSON() = %v, want map[a:1 b:2]", got)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		_, err := FromJSON[map[string]int](`{"invalid":}`)
		if err == nil {
			t.Error("FromJSON() expected error for malformed JSON")
		}
	})

	t.Run("empty JSON string", func(t *testing.T) {
		_, err := FromJSON[map[string]int]("")
		if err == nil {
			t.Error("FromJSON() expected error for empty JSON string")
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		_, err := FromJSON[map[string]int]("   ")
		if err == nil {
			t.Error("FromJSON() expected error for whitespace-only JSON string")
		}
	})
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		// Boolean types
		{"true bool", true, true, false},
		{"false bool", false, false, false},

		// Integer types
		{"int 1", 1, true, false},
		{"int 0", 0, false, false},
		{"int -1", -1, true, false},
		{"int8 1", int8(1), true, false},
		{"int16 0", int16(0), false, false},
		{"int32 42", int32(42), true, false},
		{"int64 0", int64(0), false, false},
		{"uint 1", uint(1), true, false},
		{"uint 0", uint(0), false, false},

		// Float types
		{"float32 1.0", float32(1.0), true, false},
		{"float32 0.0", float32(0.0), false, false},
		{"float64 -1.5", -1.5, true, false},
		{"float64 0.0", 0.0, false, false},

		// String types
		{"string true", "true", true, false},
		{"string false", "false", false, false},
		{"string True", "True", true, false},
		{"string FALSE", "FALSE", false, false},
		{"string t", "t", true, false},
		{"string f", "f", false, false},
		{"string yes", "yes", true, false},
		{"string no", "no", false, false},
		{"string y", "y", true, false},
		{"string n", "n", false, false},
		{"string 1", "1", true, false},
		{"string 0", "0", false, false},
		{"string on", "on", true, false},
		{"string off", "off", false, false},
		{"string empty", "", false, false},
		{"string whitespace", "   ", false, false},
		{"invalid string", "maybe", false, true},

		// Error cases
		{"nil", nil, false, true},
		{"unsupported type", []int{1, 2, 3}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		// Float types
		{"float64", 42.5, 42.5, false},
		{"float32", float32(42.5), 42.5, false},

		// Integer types
		{"int", 42, 42.0, false},
		{"int8", int8(42), 42.0, false},
		{"int16", int16(42), 42.0, false},
		{"int32", int32(42), 42.0, false},
		{"int64", int64(42), 42.0, false},
		{"uint", uint(42), 42.0, false},
		{"uint8", uint8(42), 42.0, false},
		{"uint16", uint16(42), 42.0, false},
		{"uint32", uint32(42), 42.0, false},
		{"uint64", uint64(42), 42.0, false},

		// String types
		{"valid string", "123.45", 123.45, false},
		{"integer string", "42", 42.0, false},
		{"negative string", "-123.45", -123.45, false},
		{"invalid string", "abc", 0.0, true},
		{"empty string", "", 0.0, true},
		{"whitespace string", "   ", 0.0, true},

		// Boolean types
		{"true bool", true, 1.0, false},
		{"false bool", false, 0.0, false},

		// Error cases
		{"nil", nil, 0.0, true},
		{"unsupported type", []int{1, 2, 3}, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
