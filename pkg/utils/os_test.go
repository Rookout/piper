package utils

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func createFileWithContent(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	return err
}

func TestGetFilesData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := os.TempDir()
	testDir := filepath.Join(tempDir, "test")
	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir) // Clean up the temporary directory

	// Create some dummy files in the test directory
	file1Path := filepath.Join(testDir, "file1.txt")
	err = createFileWithContent(file1Path, "File 1 data")
	if err != nil {
		t.Fatalf("Error creating file1: %v", err)
	}

	file2Path := filepath.Join(testDir, "file2.txt")
	err = createFileWithContent(file2Path, "File 2 data")
	if err != nil {
		t.Fatalf("Error creating file2: %v", err)
	}

	// Call the function being tested
	fileData, err := GetFilesData(testDir)
	if err != nil {
		t.Fatalf("Error calling GetFilesData: %v", err)
	}

	// Verify the results
	expectedData := map[string][]byte{
		"file1.txt": []byte("File 1 data"),
		"file2.txt": []byte("File 2 data"),
	}

	for fileName, expected := range expectedData {
		actual, ok := fileData[fileName]
		if !ok {
			t.Errorf("Missing file data for %s", fileName)
		}

		if !bytes.Equal(actual, expected) {
			t.Errorf("File data mismatch for %s: expected '%s', got '%s'", fileName, expected, actual)
		}
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetFilesInLinkDirectory(t *testing.T) {
	tempDir := t.TempDir()

	fileNames := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, name := range fileNames {
		filePath := filepath.Join(tempDir, name)
		_, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	linkPath := filepath.Join(tempDir, "symlink")
	err := os.Symlink(tempDir, linkPath)
	if err != nil {
		t.Fatalf("Failed to create symbolic link: %v", err)
	}

	realPath, result, err := GetFilesInLinkDirectory(linkPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedRealPath, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("Failed to evaluate symlink: %v", err)
	}
	if *realPath != expectedRealPath {
		t.Errorf("Unexpected real path. Expected: %s, got: %s", expectedRealPath, *realPath)
	}

	t.Logf("Result: %v", result) // Add this line for debugging

	if !stringSlicesEqual(result, fileNames) {
		t.Errorf("Unexpected file names. Expected: %v, got: %v", fileNames, result)
	}
}
