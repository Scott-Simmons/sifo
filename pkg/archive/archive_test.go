package archive

import (
	"SecureSyncDrive/pkg/test_helpers"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func readTarFile(tarPath string) ([]byte, error) {
	data, err := os.ReadFile(tarPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func createGzipFile(t *testing.T) *os.File {
	gzipFile, err := os.CreateTemp("", "testfile.gz")
	if err != nil {
		t.Fatalf("Failed to create temp gzip file: %v", err)
	}

	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()

	if _, err := gzipWriter.Write([]byte("This is a gzip file")); err != nil {
		t.Fatalf("Failed to write gzip content: %v", err)
	}
	gzipWriter.Close()

	// Reset file offset to the beginning
	gzipFile.Seek(0, io.SeekStart)
	return gzipFile
}

func TestFileIsTarWithTar(t *testing.T) {
	tarFile := test_helpers.CreateTarFile(t, "testing.tar", "this is content", "test.txt")
	defer os.Remove(tarFile.Name())

	file, err := os.Open(tarFile.Name())
	if err != nil {
		t.Fatalf("Failed to open tar file: %v", err)
	}
	defer file.Close()

	isTar, err := FileIsTar(file)
	if err != nil {
		t.Fatalf("fileIsTar returned an error: %v", err)
	}
	if !isTar {
		t.Error("Expected tar file to be identified as a tar file")
	}
}

func TestFileIsTarWithGzip(t *testing.T) {
	gzipFile := createGzipFile(t)
	defer os.Remove(gzipFile.Name())

	file, err := os.Open(gzipFile.Name())
	if err != nil {
		t.Fatalf("Failed to open gzip file: %v", err)
	}
	defer file.Close()

	isTar, err := FileIsTar(file)
	if err != nil {
		t.Fatalf("fileIsTar returned an error: %v", err)
	}
	if isTar {
		t.Error("Expected gzip file to be identified as not a tar file")
	}
}

func TestAddSimpleFileToTar(t *testing.T) {
	tarPath := "test_archive.tar"
	testFilePath := "test_file.txt"
	expectedTestContent := "test content"

	tarFile, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("Failed to create tar file: %v", err)
	}
	defer os.Remove(tarPath)
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFilePath) // remove file
	defer testFile.Close()        // release resources

	_, err = testFile.Write([]byte(expectedTestContent))
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	err = addFileToTar(tarWriter, testFilePath, testFilePath)
	if err != nil {
		t.Fatalf("Failed to add file to tar archive: %v", err)
	}

	tarData, err := readTarFile(tarPath)
	if err != nil {
		t.Fatalf("Failed to read tar file: %v", err)
	}

	tarReader := tar.NewReader(bytes.NewReader(tarData))
	header, err := tarReader.Next()
	if err != nil {
		t.Fatalf("Failed to read next file from tar: %v", err)
	}

	if header.Name != testFilePath {
		t.Errorf("Expected file name %s, got %s", testFilePath, header.Name)
	}

	content, err := io.ReadAll(tarReader)
	if err != nil {
		t.Fatalf("Failed to read file content from tar: %v", err)
	}

	if string(content) != expectedTestContent {
		t.Errorf("Expected file content %s, got %s", expectedTestContent, content)
	}
}

func TestAddNestedFileToTar(t *testing.T) {
	tarPath := "test_archive.tar"
	testFilePath := "test_file.txt"
	expectedTestContent := "test content"
	expectedFullTestFilePathInHeader := "directory/" + testFilePath

	tarFile, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("Failed to create tar file: %v", err)
	}
	defer os.Remove(tarPath)
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFilePath) // remove file
	defer testFile.Close()        // release resources

	_, err = testFile.Write([]byte(expectedTestContent))
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	err = addFileToTar(tarWriter, testFilePath, expectedFullTestFilePathInHeader)
	if err != nil {
		t.Fatalf("Failed to add file to tar archive: %v", err)
	}

	tarData, err := readTarFile(tarPath)
	if err != nil {
		t.Fatalf("Failed to read tar file: %v", err)
	}

	tarReader := tar.NewReader(bytes.NewReader(tarData))
	header, err := tarReader.Next()
	if err != nil {
		t.Fatalf("Failed to read next file from tar: %v", err)
	}

	if header.Name != expectedFullTestFilePathInHeader {
		t.Errorf("Expected file name %s, got %s", testFilePath, header.Name)
	}

	content, err := io.ReadAll(tarReader)
	if err != nil {
		t.Fatalf("Failed to read file content from tar: %v", err)
	}

	if string(content) != expectedTestContent {
		t.Errorf("Expected file content %s, got %s", expectedTestContent, content)
	}
}

func TestAddDirToTar(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Set up a directory structure with some files
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	filePaths := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(subDir, "file2.txt"),
		filepath.Join(tempDir, "file3.txt"),
	}

	// Create files with some content
	for _, path := range filePaths {
		err = os.WriteFile(path, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create a buffer to hold the tar file data
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	// Call the function to add the directory to the tar
	err = addDirToTar(tarWriter, tempDir, tempDir)
	if err != nil {
		t.Fatalf("Failed to add directory to tar: %v", err)
	}

	// Close the tar writer
	tarWriter.Close()

	// Read the tar file data
	tarData := buf.Bytes()
	if len(tarData) == 0 {
		t.Fatal("Tar file is empty")
	}

	// Verify the tar file contains the expected files
	tarReader := tar.NewReader(bytes.NewReader(tarData))

	expectedFiles := map[string]bool{
		"file1.txt":        true,
		"subdir/file2.txt": true,
		"file3.txt":        true,
	}

	for header, err := tarReader.Next(); err == nil; header, err = tarReader.Next() {
		if _, ok := expectedFiles[header.Name]; !ok {
			t.Fatalf("Unexpected file in tar: %s", header.Name)
		}
		content, err := io.ReadAll(tarReader)
		if err != nil {
			t.Fatalf("Failed to read file content from tar: %v", err)
		}
		if string(content) != "content" {
			t.Fatalf("Expected file content 'content', got '%s'", content)
		}
		// Remove from map once found
		delete(expectedFiles, header.Name)
	}

	// Ensure all expected files were found
	if len(expectedFiles) > 0 {
		t.Fatalf("Some expected files were not found in the tar: %v", expectedFiles)
	}
}

// TestArchiveFolder tests the ArchiveFolder function.
func TestArchiveFolder(t *testing.T) {
	srcDir := "testdata/src"
	destTar := "testdata/archive.tar.gz"

	defer os.RemoveAll(srcDir)
	defer os.Remove(destTar)

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("could not create source directory: %v", err)
	}

	testFile := filepath.Join(srcDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("could not create test file: %v", err)
	}

	if err := ArchiveFolder(srcDir, destTar); err != nil {
		t.Fatalf("ArchiveFolder returned an error: %v", err)
	}

	if err := verifyTarball(destTar, srcDir); err != nil {
		t.Fatalf("tarball verification failed: %v", err)
	}
}

func verifyTarball(tarPath, expectedDir string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		expectedFile := filepath.Join(expectedDir, header.Name)
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			return err
		}

		if header.Typeflag == tar.TypeReg {
			content, err := io.ReadAll(tarReader)
			if err != nil {
				return err
			}
			expectedContent, err := os.ReadFile(expectedFile)
			if err != nil {
				return err
			}
			if string(content) != string(expectedContent) {
				return err
			}
		}
	}

	return nil
}
