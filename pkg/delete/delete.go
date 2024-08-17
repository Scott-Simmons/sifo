package delete

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	"path/filepath"

	_ "github.com/rclone/rclone/backend/local" // local backend and backblaze backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/fs/sync"
)

type deleteRequest struct {
	// These look the wrong way around but look at: https://github.com/rclone/rclone/blob/d1c84f9115fc7a782003e3c3f00216e5caf90605/fs/operations/rc_test.go#L187
	// It is counterintuitive as hell.
	RemoteName string `json:"fs"`
	FileName   string `json:"remote"`
}

func DeleteBackBlazeFile(client rpc_client.RPCClient, remoteName string, bucketName string, fileName string) error {
	return Delete(client, remoteName, filepath.Join(bucketName, fileName))
}
func Delete(client rpc_client.RPCClient, remoteName string, fileName string) error {
	if err := client.Initialize(); err != nil {
		return err
	}
	const deleteMethod = "operations/deletefile"
	deleteRequest := deleteRequest{
		RemoteName: remoteName,
		FileName:   fileName,
	}
	deleteRequestJson, err := json.Marshal(deleteRequest)
	if err != nil {
		return err
	}
	out, status := client.RPC(deleteMethod, string(deleteRequestJson))
	if status != 200 {
		return fmt.Errorf("Error status: %d and error output: %s", status, out)
	}
	return nil
}
