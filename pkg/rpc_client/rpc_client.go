package rpc_client

import (
	_ "github.com/rclone/rclone/backend/drive"
	_ "github.com/rclone/rclone/backend/local" // local backend and google drive backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/fs/sync"
	"github.com/rclone/rclone/librclone/librclone"
)

type RPCClient interface {
	Initialize() error
	RPC(method string, params string) (string, int)
}

type RealRPCClient struct{}

func (c *RealRPCClient) RPC(method string, params string) (string, int) {
	return librclone.RPC(method, params)
}
func (c *RealRPCClient) Initialize() error {
	librclone.Initialize()
	return nil
}

