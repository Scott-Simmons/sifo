package validate_config

import (
	"os"
	"testing"
)

func createTestFile(t *testing.T, content string) string {
	t.Helper()
	file, err := os.CreateTemp("", "config*.ini")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return file.Name()
}

func TestParseConfigSuccess(t *testing.T) {
	iniContent := `
[backblazeRemote1]
type = b2
account = keyname
key = keyvalue
hard_delete = true

[backblazeRemote2]
type = b2
account = keyname
key = keyvalue
hard_delete = true
`

	filePath := createTestFile(t, iniContent)
	defer os.Remove(filePath)

	err := ValidateConfig(filePath)
	if err != nil {
		t.Fatalf("ParseConfig() returned an error: %v", err)
	}
}

func TestParseConfigInvalidType(t *testing.T) {

	iniContent := `
[backblazeRemote1]
type = somethinginvalid
account = keyname
key = keyvalue
hard_delete = true
`

	filePath := createTestFile(t, iniContent)
	defer os.Remove(filePath)

	err := ValidateConfig(filePath)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expectedError := "Config validation failed for backblazeRemote1"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestParseConfigFileError(t *testing.T) {
	err := ValidateConfig("/invalid/path/to/config.ini")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestParseConfigINIUnmarshalError(t *testing.T) {
	invalidINIContent := `[invalid INI missing bracket`

	filePath := createTestFile(t, invalidINIContent)
	defer os.Remove(filePath)

	err := ValidateConfig(filePath)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
