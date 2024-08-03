package encrypt

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func setupTestDataDir() error {
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		if err := os.MkdirAll("testdata", 0755); err != nil {
			return err
		}
	}
	return nil
}

func setupTempKeyFile(t *testing.T, key []byte, pathToWriteKeyInto string) (*os.File, error) {
	tempKeyFile, err := os.Create(pathToWriteKeyInto)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	hexData := hex.EncodeToString(key)

	if _, err := tempKeyFile.Write([]byte(hexData)); err != nil {
		tempKeyFile.Close()
		return nil, err
	}

	return tempKeyFile, nil
}

func TestEncryptTarBall_Success(t *testing.T) {
	const bytesForKeySize = 32
	if err := setupTestDataDir(); err != nil {
		t.Fatalf("Failed to set up testdata directory: %v", err)
	}
	srcFilePath := "testdata/testfile.tar"
	encryptedFilePath := "testdata/testfile.tar.enc"
	key := make([]byte, bytesForKeySize)

	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	tempKeyFilePath := "tempkey.key"
	if _, err := setupTempKeyFile(t, key, tempKeyFilePath); err != nil {
		t.Fatalf("Failed to set up key file")
	}

	if err := createSampleTarFile(srcFilePath); err != nil {
		t.Fatalf("Failed to create sample tar file: %v", err)
	}
	defer os.Remove(srcFilePath)

	fmt.Printf("Using keypath: %s", tempKeyFilePath)
	content, _ := os.ReadFile(tempKeyFilePath)
	fmt.Println("Lets see the content")
	fmt.Println(string(content))
	if err := EncryptTarBall(srcFilePath, encryptedFilePath, tempKeyFilePath); err != nil {
		t.Fatalf("EncryptTarBall returned an error: %v", err)
	}
	defer os.Remove(encryptedFilePath)

	// Check if the encrypted file is created
	if _, err := os.Stat(encryptedFilePath); os.IsNotExist(err) {
		t.Fatalf("Encrypted file was not created")
	}
}

func TestEncryptTarBall_InvalidKeyLength(t *testing.T) {
	if err := setupTestDataDir(); err != nil {
		t.Fatalf("Failed to set up testdata directory: %v", err)
	}

	srcFilePath := "testdata/testfile.tar"
	encryptedFilePath := "testdata/testfile.tar.enc"
	key := make([]byte, 16) // Invalid key length

	// Create a sample tar file
	if err := createSampleTarFile(srcFilePath); err != nil {
		t.Fatalf("Failed to create sample tar file: %v", err)
	}
	defer os.Remove(srcFilePath)

	tempKeyFilePath := "tempkey.key"
	if _, err := setupTempKeyFile(t, key, tempKeyFilePath); err != nil {
		t.Fatalf("Failed to set up key file")
	}

	// Encrypt the tar file
	err := EncryptTarBall(srcFilePath, encryptedFilePath, tempKeyFilePath)
	if err == nil {
		t.Fatalf("Expected error due to invalid key length, got nil")
	}
}

// TODO: Right now the tar validation is not working, need to fix it up
//func TestEncryptTarBall_InvalidTar(t *testing.T) {
//  if err := setupTestDataDir(); err != nil {
//		t.Fatalf("Failed to set up testdata directory: %v", err)
//	}
//
//	srcFilePath := "testdata/invalidfile.txt"
//	encryptedFilePath := "testdata/invalidfile.txt.enc"
//	key := make([]byte, 32)
//	if _, err := rand.Read(key); err != nil {
//		t.Fatalf("Failed to generate random key: %v", err)
//	}
//
//	// Create a sample non-tar file
//	if err := os.WriteFile(srcFilePath, []byte("This is not a tar file"), 0644); err != nil {
//		t.Fatalf("Failed to create sample non-tar file: %v", err)
//	}
//	defer os.Remove(srcFilePath)
//
//	// Encrypt the non-tar file
//  tempKeyFilePath := "tempkey.key"
//  if _, err := setupTempKeyFile(t, key, tempKeyFilePath); err != nil {
//    t.Fatalf("Failed to set up key file")
//  }
//	err := EncryptTarBall(srcFilePath, encryptedFilePath, tempKeyFilePath)
//	if err == nil {
//		t.Fatalf("Expected error due to invalid tar file, got nil")
//	}
//}

func createSampleTarFile(path string) error {
	// Create a simple tar file for testing
	fileContent := "test data"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO: Make it a real tar with actual headers
	if _, err := f.Write([]byte(fileContent)); err != nil {
		return err
	}

	return nil
}

func TestPad(t *testing.T) {
	tests := []struct {
		data      []byte
		blockSize int
		expected  []byte
	}{
		{[]byte("test"), 8, []byte("test\x04\x04\x04\x04")},
		{[]byte("testdata"), 16, []byte("testdata\x08\x08\x08\x08\x08\x08\x08\x08")},
		{[]byte(""), 16, []byte("\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10")},
	}

	for _, tt := range tests {
		result := pad(tt.data, tt.blockSize)
		if !bytes.Equal(result, tt.expected) {
			t.Errorf("pad(%v, %d) = %v; want %v", tt.data, tt.blockSize, result, tt.expected)
		}
	}
}
