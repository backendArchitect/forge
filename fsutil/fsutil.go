// Package fsutil provides file system utilities for common operations.
// All functions handle errors gracefully and follow Go conventions.
package fsutil

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Exists checks if a file or directory exists at the given path.
//
// Example:
//
//	if fsutil.Exists("/path/to/file.txt") {
//		fmt.Println("File exists")
//	}
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsFile checks if the given path exists and is a regular file.
//
// Example:
//
//	if fsutil.IsFile("/path/to/file.txt") {
//		fmt.Println("Path is a file")
//	}
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDir checks if the given path exists and is a directory.
//
// Example:
//
//	if fsutil.IsDir("/path/to/directory") {
//		fmt.Println("Path is a directory")
//	}
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// ReadJSON reads a JSON file and unmarshals it into the provided value.
//
// Example:
//
//	type Config struct {
//		Host string `json:"host"`
//		Port int    `json:"port"`
//	}
//	var config Config
//	err := fsutil.ReadJSON("config.json", &config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Config: %+v\n", config)
func ReadJSON(path string, v any) error {
	if v == nil {
		return os.ErrInvalid
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(v)
}

// WriteJSON marshals the provided value to JSON and writes it to a file.
// The file is created if it doesn't exist, or truncated if it does.
//
// Example:
//
//	config := Config{Host: "localhost", Port: 8080}
//	err := fsutil.WriteJSON("config.json", config)
//	if err != nil {
//		log.Fatal(err)
//	}
func WriteJSON(path string, v any) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	return encoder.Encode(v)
}

// CopyFile copies a file from src to dst, preserving permissions.
// If dst already exists, it will be overwritten.
//
// Example:
//
//	err := fsutil.CopyFile("source.txt", "destination.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Set permissions to match source
	return os.Chmod(dst, srcInfo.Mode())
}

// EnsureDir creates a directory and all necessary parent directories.
// If the directory already exists, no error is returned.
//
// Example:
//
//	err := fsutil.EnsureDir("/path/to/nested/directory")
//	if err != nil {
//		log.Fatal(err)
//	}
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// ReadFile reads the entire content of a file and returns it as a string.
//
// Example:
//
//	content, err := fsutil.ReadFile("config.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(content)
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes string content to a file, creating it if it doesn't exist.
// If the file already exists, it will be truncated.
//
// Example:
//
//	err := fsutil.WriteFile("output.txt", "Hello, World!")
//	if err != nil {
//		log.Fatal(err)
//	}
func WriteFile(path, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}
