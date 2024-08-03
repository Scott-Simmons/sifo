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

func TestSyncGoogleDriveWithMock(t *testing.T) {
    client := new(MockRPCClient)
    client.On("Initialize", mock.Anything).Return(nil)
    client.On("RPC", "sync/sync", mock.Anything).Return(`{"success": {}}`, 200)

    err := SyncGoogleDrive(client, "test_file.txt", "remote:path/to/drive/file.txt")
    if err != nil {
        t.Errorf("syncGoogleDrive failed: %v", err)
    }

    client.AssertExpectations(t)
}
