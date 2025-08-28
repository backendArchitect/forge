package fsutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExists(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_exists_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_exists_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing file", tmpFile.Name(), true},
		{"existing directory", tmpDir, true},
		{"non-existent path", "/path/that/does/not/exist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exists(tt.path)
			if got != tt.want {
				t.Errorf("Exists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_isfile_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_isfile_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"regular file", tmpFile.Name(), true},
		{"directory", tmpDir, false},
		{"non-existent path", "/path/that/does/not/exist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFile(tt.path)
			if got != tt.want {
				t.Errorf("IsFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_isdir_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_isdir_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"directory", tmpDir, true},
		{"regular file", tmpFile.Name(), false},
		{"non-existent path", "/path/that/does/not/exist", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDir(tt.path)
			if got != tt.want {
				t.Errorf("IsDir(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestReadWriteJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	original := TestStruct{Name: "John", Age: 30}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_json_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	t.Run("write and read JSON", func(t *testing.T) {
		// Write JSON
		err := WriteJSON(tmpFile.Name(), original)
		if err != nil {
			t.Fatalf("WriteJSON() error = %v", err)
		}

		// Read JSON back
		var result TestStruct
		err = ReadJSON(tmpFile.Name(), &result)
		if err != nil {
			t.Fatalf("ReadJSON() error = %v", err)
		}

		if !reflect.DeepEqual(result, original) {
			t.Errorf("ReadJSON() = %v, want %v", result, original)
		}
	})

	t.Run("read non-existent file", func(t *testing.T) {
		var result TestStruct
		err := ReadJSON("/path/that/does/not/exist.json", &result)
		if err == nil {
			t.Error("ReadJSON() expected error for non-existent file")
		}
	})

	t.Run("read with nil pointer", func(t *testing.T) {
		err := ReadJSON(tmpFile.Name(), nil)
		if err == nil {
			t.Error("ReadJSON() expected error for nil pointer")
		}
	})

	t.Run("write to nested directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test_nested_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		nestedPath := filepath.Join(tmpDir, "nested", "data.json")
		err = WriteJSON(nestedPath, original)
		if err != nil {
			t.Fatalf("WriteJSON() to nested path error = %v", err)
		}

		// Verify file was created
		if !Exists(nestedPath) {
			t.Error("WriteJSON() did not create nested file")
		}
	})
}

func TestCopyFile(t *testing.T) {
	// Create a temporary source file
	srcFile, err := os.CreateTemp("", "test_copy_src_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(srcFile.Name())

	testContent := "Hello, World!"
	_, err = srcFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	srcFile.Close()

	// Set specific permissions
	err = os.Chmod(srcFile.Name(), 0644)
	if err != nil {
		t.Fatalf("Failed to set permissions: %v", err)
	}

	t.Run("copy file successfully", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test_copy_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		dstPath := filepath.Join(tmpDir, "copied.txt")
		err = CopyFile(srcFile.Name(), dstPath)
		if err != nil {
			t.Fatalf("CopyFile() error = %v", err)
		}

		// Verify destination file exists
		if !Exists(dstPath) {
			t.Error("CopyFile() did not create destination file")
		}

		// Verify content
		content, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("Failed to read copied file: %v", err)
		}

		if string(content) != testContent {
			t.Errorf("CopyFile() content = %q, want %q", string(content), testContent)
		}

		// Verify permissions
		srcInfo, _ := os.Stat(srcFile.Name())
		dstInfo, _ := os.Stat(dstPath)
		if srcInfo.Mode() != dstInfo.Mode() {
			t.Errorf("CopyFile() permissions = %v, want %v", dstInfo.Mode(), srcInfo.Mode())
		}
	})

	t.Run("copy to nested directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test_copy_nested_*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		dstPath := filepath.Join(tmpDir, "nested", "dir", "copied.txt")
		err = CopyFile(srcFile.Name(), dstPath)
		if err != nil {
			t.Fatalf("CopyFile() to nested path error = %v", err)
		}

		if !Exists(dstPath) {
			t.Error("CopyFile() did not create nested destination file")
		}
	})

	t.Run("copy non-existent source", func(t *testing.T) {
		err := CopyFile("/path/that/does/not/exist", "/tmp/dest")
		if err == nil {
			t.Error("CopyFile() expected error for non-existent source")
		}
	})
}
