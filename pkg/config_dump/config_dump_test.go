package config_dump

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

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

// Test needs to be more comprehensive
func TestDumpConfig(t *testing.T) {
	client := new(MockRPCClient)
	client.On("Initialize", mock.Anything).Return(nil)
	client.On("RPC", "config/dump", mock.Anything).Return(`{"success": {}}`, 200)

	_, err := DumpConfig(client)
	if err != nil {
		t.Errorf("DumpConfig failed: %v", err)
	}
	client.AssertExpectations(t)
}
