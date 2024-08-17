package sync

// Ref: https://rclone.org/docs/#connection-strings
// Ref: https://forum.rclone.org/t/use-rclone-golang-to-transfer-files/34983/23?page=2
// Ref: https://github.com/alankritkharbanda/rclone/blob/5975b7d27728f5ba0c3c670759fe9cc3dfb65ff2/librclone/README.md

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	_ "github.com/rclone/rclone/backend/drive"
	_ "github.com/rclone/rclone/backend/local" // local backend and backblaze backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/fs/sync"
)

type syncRequest struct {
	SrcFs string `json:"srcFs"`
	DstFs string `json:"dstFs"` // best to keep struct as PascalCase. Not accessible otherwise
}

func NewClient() (*rpc_client.RealRPCClient, error) {
	client := &rpc_client.RealRPCClient{}
	return client, nil
}

// NOTE: This is destructive... it makes dst look exactly like src.
func sync(client rpc_client.RPCClient, src string, dst string) error {
	if err := client.Initialize(); err != nil {
		return err
	}

	// Ref: https://github.com/rclone/rclone/blob/master/fs/sync/sync.go
	const syncMethod = "sync/sync"

	syncRequest := syncRequest{
		SrcFs: src,
		DstFs: dst,
	}

	syncRequestJson, err := json.Marshal(syncRequest)
	if err != nil {
		return err
	}
	fmt.Println(string(syncRequestJson))
	out, status := client.RPC(syncMethod, string(syncRequestJson))

	if status != 200 {
		return fmt.Errorf("Error status: %d and error output: %s", status, out)
	} else {
		fmt.Printf("Success: %s\n", out)
	}
	fmt.Println(out)
	fmt.Println(status)
	return nil
}

// TODO: Validation
func SyncToBackblaze(client rpc_client.RPCClient, srcFilePath string, backblazeFilePath string) error {
	if err := sync(client, srcFilePath, backblazeFilePath); err != nil {
		return err
	}
	return nil
}

// TODO: Validation
func SyncFromBackblaze(client rpc_client.RPCClient, backblazeFilePath string, localTargetDir string) error {
	if err := sync(client, backblazeFilePath, localTargetDir); err != nil {
		return err
	}
	return nil
}
