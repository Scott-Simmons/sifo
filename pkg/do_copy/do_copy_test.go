package do_copy

// TODO: Lots of test tidyings and improvement of test quality.
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
func TestSyncFromBackblazeDriveWithMock(t *testing.T) {
	client := new(MockRPCClient)
	client.On("Initialize", mock.Anything).Return(nil)
	client.On("RPC", "operations/copyfile", mock.Anything).Return(`{"success": {}}`, 200)

	err := CopyFromBackblaze(client, "path/to/drive/file.txt", "remote:", "bucket", "foo_dir/")
	if err != nil {
		t.Errorf("CopyFromBackblazeDrive failed: %v", err)
	}
	client.AssertExpectations(t)
}
