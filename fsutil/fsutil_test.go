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

func TestEnsureDir(t *testing.T) {
	t.Run("create new directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		dirPath := filepath.Join(tmpDir, "new", "nested", "directory")

		err := EnsureDir(dirPath)
		if err != nil {
			t.Fatalf("EnsureDir() error = %v", err)
		}

		if !IsDir(dirPath) {
			t.Error("EnsureDir() did not create directory")
		}
	})

	t.Run("directory already exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		err := EnsureDir(tmpDir)
		if err != nil {
			t.Fatalf("EnsureDir() error on existing directory = %v", err)
		}

		if !IsDir(tmpDir) {
			t.Error("EnsureDir() should not affect existing directory")
		}
	})
}

func TestReadFile(t *testing.T) {
	t.Run("read existing file", func(t *testing.T) {
		content := "Hello, World!\nThis is a test file."
		tmpFile, err := os.CreateTemp("", "readtest")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(content)
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		result, err := ReadFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}

		if result != content {
			t.Errorf("ReadFile() = %q, want %q", result, content)
		}
	})

	t.Run("read non-existent file", func(t *testing.T) {
		_, err := ReadFile("/path/that/does/not/exist")
		if err == nil {
			t.Error("ReadFile() expected error for non-existent file")
		}
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("write to new file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		content := "Hello, World!\nThis is a test."

		err := WriteFile(filePath, content)
		if err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if !IsFile(filePath) {
			t.Error("WriteFile() did not create file")
		}

		result, err := ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read written file: %v", err)
		}

		if result != content {
			t.Errorf("Written content = %q, want %q", result, content)
		}
	})

	t.Run("write to nested path", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "nested", "dir", "test.txt")
		content := "Nested file content"

		err := WriteFile(filePath, content)
		if err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if !IsFile(filePath) {
			t.Error("WriteFile() did not create nested file")
		}

		result, err := ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read nested file: %v", err)
		}

		if result != content {
			t.Errorf("Nested file content = %q, want %q", result, content)
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "writetest")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		originalContent := "Original content"
		newContent := "New content"

		_, err = tmpFile.WriteString(originalContent)
		if err != nil {
			t.Fatalf("Failed to write original content: %v", err)
		}
		tmpFile.Close()

		err = WriteFile(tmpFile.Name(), newContent)
		if err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		result, err := ReadFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to read overwritten file: %v", err)
		}

		if result != newContent {
			t.Errorf("Overwritten content = %q, want %q", result, newContent)
		}
	})
}
