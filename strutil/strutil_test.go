package strutil

import (
	"strings"
	"testing"
)

func TestRandom(t *testing.T) {
	t.Run("correct length", func(t *testing.T) {
		length := 10
		result := Random(length)
		if len(result) != length {
			t.Errorf("Random() length = %d, want %d", len(result), length)
		}
	})

	t.Run("zero length", func(t *testing.T) {
		result := Random(0)
		if result != "" {
			t.Errorf("Random(0) = %q, want empty string", result)
		}
	})

	t.Run("negative length", func(t *testing.T) {
		result := Random(-5)
		if result != "" {
			t.Errorf("Random(-5) = %q, want empty string", result)
		}
	})

	t.Run("custom charset", func(t *testing.T) {
		charset := "ABC"
		result := Random(10, charset)
		if len(result) != 10 {
			t.Errorf("Random() length = %d, want 10", len(result))
		}
		for _, char := range result {
			if !strings.ContainsRune(charset, char) {
				t.Errorf("Random() contains invalid character %c", char)
			}
		}
	})

	t.Run("empty charset", func(t *testing.T) {
		result := Random(5, "")
		// Should use default charset
		if len(result) != 5 {
			t.Errorf("Random() with empty charset length = %d, want 5", len(result))
		}
	})

	t.Run("default charset", func(t *testing.T) {
		result := Random(10)
		if len(result) != 10 {
			t.Errorf("Random() length = %d, want 10", len(result))
		}
		// Check that result contains only alphanumeric characters
		for _, char := range result {
			if !strings.ContainsRune(defaultCharset, char) {
				t.Errorf("Random() contains invalid character %c", char)
			}
		}
	})
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"single space", " ", true},
		{"multiple spaces", "   ", true},
		{"tab and newline", "\t\n", true},
		{"mixed whitespace", " \t\n ", true},
		{"string with content", " a ", false},
		{"normal string", "hello", false},
		{"string with spaces between", "a b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBlank(tt.input)
			if got != tt.want {
				t.Errorf("IsBlank(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		suffix string
		want   string
	}{
		{"no truncation needed", "Hello", 10, "...", "Hello"},
		{"exact length", "Hello", 5, "...", "Hello"},
		{"truncation needed", "Hello World", 5, "...", "He..."},
		{"empty suffix", "Hello World", 5, "", "Hello"},
		{"suffix longer than length", "Hello", 2, "...", ".."},
		{"zero length", "Hello", 0, "...", ""},
		{"negative length", "Hello", -1, "...", ""},
		{"empty string", "", 5, "...", ""},
		{"suffix same as length", "Hello World", 3, "...", "..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.length, tt.suffix)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.suffix, got, tt.want)
			}
		})
	}
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple camelCase", "myFunctionName", "my_function_name"},
		{"single word", "function", "function"},
		{"already lowercase", "alllowercase", "alllowercase"},
		{"multiple capitals", "XMLHttpRequest", "x_m_l_http_request"},
		{"starting with capital", "MyFunction", "my_function"},
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"single capital", "A", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CamelToSnake(tt.input)
			if got != tt.want {
				t.Errorf("CamelToSnake(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple snake_case", "my_function_name", "myFunctionName"},
		{"single word", "function", "function"},
		{"already camelCase", "myFunction", "myFunction"},
		{"multiple underscores", "my__function", "myFunction"},
		{"trailing underscore", "my_function_", "myFunction"},
		{"leading underscore", "_my_function", "myFunction"},
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"single underscore", "_", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SnakeToCamel(tt.input)
			if got != tt.want {
				t.Errorf("SnakeToCamel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTemplate(t *testing.T) {
	t.Run("valid template", func(t *testing.T) {
		template := "Hello {{name}}, you are {{age}} years old"
		data := map[string]any{"name": "John", "age": 30}
		got, err := Template(template, data)
		if err != nil {
			t.Fatalf("Template() error = %v", err)
		}
		want := "Hello John, you are 30 years old"
		if got != want {
			t.Errorf("Template() = %q, want %q", got, want)
		}
	})

	t.Run("template with missing keys", func(t *testing.T) {
		template := "Hello {{name}}, you are {{age}} years old"
		data := map[string]any{"name": "John"}
		_, err := Template(template, data)
		if err == nil {
			t.Error("Template() expected error for missing keys")
		}
	})

	t.Run("empty template", func(t *testing.T) {
		got, err := Template("", nil)
		if err != nil {
			t.Fatalf("Template() error = %v", err)
		}
		if got != "" {
			t.Errorf("Template() = %q, want empty string", got)
		}
	})

	t.Run("no placeholders", func(t *testing.T) {
		template := "Hello World"
		got, err := Template(template, nil)
		if err != nil {
			t.Fatalf("Template() error = %v", err)
		}
		if got != template {
			t.Errorf("Template() = %q, want %q", got, template)
		}
	})

	t.Run("nil data map", func(t *testing.T) {
		template := "Hello World"
		got, err := Template(template, nil)
		if err != nil {
			t.Fatalf("Template() error = %v", err)
		}
		if got != template {
			t.Errorf("Template() = %q, want %q", got, template)
		}
	})

	t.Run("template with spaces in keys", func(t *testing.T) {
		template := "Hello {{ name }}, you are {{ age }} years old"
		data := map[string]any{"name": "John", "age": 30}
		got, err := Template(template, data)
		if err != nil {
			t.Fatalf("Template() error = %v", err)
		}
		want := "Hello John, you are 30 years old"
		if got != want {
			t.Errorf("Template() = %q, want %q", got, want)
		}
	})

	t.Run("template with different value types", func(t *testing.T) {
		template := "Name: {{name}}, Age: {{age}}, Score: {{score}}, Active: {{active}}"
		data := map[string]any{
			"name":   "John",
			"age":    30,
			"score":  95.5,
			"active": true,
		}
		got, err := Template(template, data)
		if err != nil {
			t.Fatalf("Template() error = %v", err)
		}
		want := "Name: John, Age: 30, Score: 95.5, Active: true"
		if got != want {
			t.Errorf("Template() = %q, want %q", got, want)
		}
	})
}

func TestPad(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		length  int
		padChar rune
		want    string
	}{
		{"pad with spaces", "hello", 10, ' ', "hello     "},
		{"pad with asterisks", "hi", 6, '*', "hi****"},
		{"no padding needed", "hello", 5, '-', "hello"},
		{"no padding for longer string", "hello world", 5, '-', "hello world"},
		{"pad single char", "a", 3, 'x', "axx"},
		{"pad empty string", "", 5, 'z', "zzzzz"},
		{"zero length", "hello", 0, ' ', "hello"},
		{"negative length", "hello", -1, ' ', "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pad(tt.input, tt.length, tt.padChar)
			if got != tt.want {
				t.Errorf("Pad(%q, %d, %q) = %q, want %q", tt.input, tt.length, tt.padChar, got, tt.want)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase word", "hello", "Hello"},
		{"uppercase word", "HELLO", "Hello"},
		{"mixed case", "hELLo", "Hello"},
		{"sentence", "hello WORLD", "Hello world"},
		{"single char lowercase", "a", "A"},
		{"single char uppercase", "A", "A"},
		{"empty string", "", ""},
		{"with numbers", "hello123", "Hello123"},
		{"with spaces", "hello world", "Hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Capitalize(tt.input)
			if got != tt.want {
				t.Errorf("Capitalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple string", "hello", "olleh"},
		{"single char", "a", "a"},
		{"empty string", "", ""},
		{"palindrome", "racecar", "racecar"},
		{"with spaces", "hello world", "dlrow olleh"},
		{"with numbers", "abc123", "321cba"},
		{"unicode", "hello 🌍", "🌍 olleh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse(tt.input)
			if got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
