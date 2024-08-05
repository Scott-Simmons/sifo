package list_files

import (
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

// TODO: Get rid of mocking duplication later

func GenGoogleDriveFilesTestData() []GoogleDriveFile {
	return []GoogleDriveFile{
		{
			Path:     "foo/bar.txt",
			Name:     "bar.txt",
			Size:     5,
			MimeType: "text/plain",
			ModTime:  "2024-08-04T14:43:02.022Z",
			IsDir:    false,
			ID:       "1YImqqzHic8Acx0-4kyiENxki6i6Sgf5z",
		},
		{
			Path:     "files/yoo.tar.gzip",
			Name:     "yoo.tar.gzip",
			Size:     160,
			MimeType: "application/octet-stream",
			ModTime:  "2024-08-04T14:49:55.776Z",
			IsDir:    false,
			ID:       "1RUD_wmTvjg9g-tFBjEbEFb2Sr2jXmgPi",
		},
	}
}

type MockLibrclone struct {
	mock.Mock
}

type MockRPCClient struct {
	mock.Mock
}

func (m *MockRPCClient) Initialize() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRPCClient) RPC(method string, params string) (string, int) {
	args := m.Called(method, params)
	return args.String(0), args.Int(1)
}

// TODO: This is a weak test. Need to expand it... edge cases
func TestListFiles(t *testing.T) {
	// Needs a real integration test to test output
	client := new(MockRPCClient)
	client.On("Initialize", mock.Anything).Return(nil)
	client.On("RPC", "operations/list", mock.Anything).Return(`{"success": {}}`, 200)
	// TODO: Test output properly
	_, err := ListRemoteFiles(client, "remote:")
	if err != nil {
		t.Errorf("List files failed: %v", err)
	}
	client.AssertExpectations(t)
}

func TestGetGoogleDriveFileIds(t *testing.T) {
	expectedIds := []string{
		"1YImqqzHic8Acx0-4kyiENxki6i6Sgf5z",
		"1RUD_wmTvjg9g-tFBjEbEFb2Sr2jXmgPi",
	}
	actualIds := GetGoogleDriveFileIds(GenGoogleDriveFilesTestData())
	if !reflect.DeepEqual(actualIds, expectedIds) {
		t.Errorf("Expected IDs %v, but got %v", expectedIds, actualIds)
	}
}
