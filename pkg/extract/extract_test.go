package extract

import (
	"SecureSyncDrive/pkg/test_helpers"
	"bytes"
	"io"
	"os"
	"testing"
)

func TestExtractTar(t *testing.T) {
	testFile := "testfile.txt"
	tarFile := test_helpers.CreateGzipTarFile(t, "testFile.tar", "This is a test file", testFile)
	defer os.Remove(tarFile.Name())

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, tarFile)
	if err != nil {
		t.Fatalf("Failed to write file content into memory")
	}
	tarData := buffer.Bytes()
	err = ExtractTar(tarData)
	if err != nil {
		t.Fatalf("Failed to extract tar %v", err)
	}

	extractedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	expectedContent := "This is a test file"
	if string(extractedContent) != expectedContent {
		t.Errorf("Extracted file content does not match. Expected: %q, got: %q", expectedContent, string(extractedContent))
	}

}
