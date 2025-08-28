// Package strutil provides string manipulation utilities.
// All functions handle edge cases gracefully and follow Go conventions.
package strutil

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Random generates a random string of the specified length using the provided charset.
// If no charset is provided, uses alphanumeric characters.
//
// Example:
//
//	random := strutil.Random(8)
//	fmt.Println(random) // Output: "aBc3DeF9" (random alphanumeric string)
//
//	custom := strutil.Random(5, "12345")
//	fmt.Println(custom) // Output: "42153" (random string from custom charset)
func Random(length int, charset ...string) string {
	if length <= 0 {
		return ""
	}

	chars := defaultCharset
	if len(charset) > 0 && charset[0] != "" {
		chars = charset[0]
	}

	if len(chars) == 0 {
		return ""
	}

	// Initialize random seed if not already done
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

// IsBlank returns true if the string is empty or contains only whitespace characters.
//
// Example:
//
//	fmt.Println(strutil.IsBlank(""))        // Output: true
//	fmt.Println(strutil.IsBlank("   "))     // Output: true
//	fmt.Println(strutil.IsBlank(" a "))     // Output: false
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Truncate truncates a string to the specified length and appends a suffix if truncation occurs.
//
// Example:
//
//	truncated := strutil.Truncate("Hello World", 5, "...")
//	fmt.Println(truncated) // Output: "Hello..."
//
//	noTruncate := strutil.Truncate("Hi", 5, "...")
//	fmt.Println(noTruncate) // Output: "Hi"
func Truncate(s string, length int, suffix string) string {
	if length < 0 {
		return ""
	}

	if len(s) <= length {
		return s
	}

	if length < len(suffix) {
		if length == 0 {
			return ""
		}
		return suffix[:length]
	}

	return s[:length-len(suffix)] + suffix
}

// CamelToSnake converts camelCase strings to snake_case.
//
// Example:
//
//	snake := strutil.CamelToSnake("myFunctionName")
//	fmt.Println(snake) // Output: "my_function_name"
func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}

	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteByte('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// SnakeToCamel converts snake_case strings to camelCase.
//
// Example:
//
//	camel := strutil.SnakeToCamel("my_function_name")
//	fmt.Println(camel) // Output: "myFunctionName"
func SnakeToCamel(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	if len(parts) <= 1 {
		return s
	}

	var result strings.Builder
	firstPart := true

	for i := 0; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			if firstPart {
				result.WriteString(parts[i])
				firstPart = false
			} else {
				result.WriteString(strings.Title(parts[i]))
			}
		}
	}

	return result.String()
}

// Template performs simple template substitution using map values.
// Template variables are specified as {{key}} in the template string.
//
// Example:
//
//	template := "Hello {{name}}, you are {{age}} years old"
//	data := map[string]any{"name": "John", "age": 30}
//	result, err := strutil.Template(template, data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result) // Output: "Hello John, you are 30 years old"
func Template(template string, data map[string]any) (string, error) {
	if template == "" {
		return "", nil
	}

	if data == nil {
		data = make(map[string]any)
	}

	// Regular expression to find {{key}} patterns
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	result := re.ReplaceAllStringFunc(template, func(match string) string {
		// Extract the key (remove {{ and }})
		key := match[2 : len(match)-2]
		key = strings.TrimSpace(key)

		if value, exists := data[key]; exists {
			return fmt.Sprintf("%v", value)
		}

		// Return the original match if key not found
		return match
	})

	// Check if there are any unresolved templates
	if re.MatchString(result) {
		unresolvedKeys := re.FindAllStringSubmatch(result, -1)
		var missingKeys []string
		for _, match := range unresolvedKeys {
			if len(match) > 1 {
				key := strings.TrimSpace(match[1])
				if _, exists := data[key]; !exists {
					missingKeys = append(missingKeys, key)
				}
			}
		}
		if len(missingKeys) > 0 {
			return result, fmt.Errorf("missing template keys: %v", missingKeys)
		}
	}

	return result, nil
}
