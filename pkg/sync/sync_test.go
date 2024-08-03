package sync

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

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
func TestSyncToGoogleDriveWithMock(t *testing.T) {
	client := new(MockRPCClient)
	client.On("Initialize", mock.Anything).Return(nil)
	client.On("RPC", "sync/sync", mock.Anything).Return(`{"success": {}}`, 200)

	err := SyncToGoogleDrive(client, "test_file.txt", "remote:path/to/drive/file.txt")
	if err != nil {
		t.Errorf("SyncToGoogleDrive failed: %v", err)
	}

	client.AssertExpectations(t)
}

// TODO: This is a weak test. Need to expand it... edge cases
func TestSyncFromGoogleDriveWithMock(t *testing.T) {
	client := new(MockRPCClient)
	client.On("Initialize", mock.Anything).Return(nil)
	client.On("RPC", "sync/sync", mock.Anything).Return(`{"success": {}}`, 200)

	err := SyncFromGoogleDrive(client, "remote:path/to/drive/file.txt", "foo_dir/")
	if err != nil {
		t.Errorf("SyncFromGoogleDrive failed: %v", err)
	}

	client.AssertExpectations(t)

}
