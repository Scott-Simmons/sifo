package extract

// TODO: Flesh out tests

//import (
//	"SecureSyncDrive/pkg/test_helpers"
//	"os"
//	"testing"
//)

//func TestExtractTar(t *testing.T) {
//	testFile := "testfile.txt"
//	test_helpers.CreateGzipTarFile(t, "testFile.tar", "This is a test file", testFile)
//	err := ExtractTar("testFile.tar", ".")
//	if err != nil {
//		t.Fatalf("Failed to extract tar %v", err)
//	}
//	extractedContent, err := os.ReadFile(testFile)
//	if err != nil {
//		t.Fatalf("Failed to read extracted file: %v", err)
//	}
//	expectedContent := "This is a test file"
//	if string(extractedContent) != expectedContent {
//		t.Errorf("Extracted file content does not match. Expected: %q, got: %q", expectedContent, string(extractedContent))
//	}
//}
