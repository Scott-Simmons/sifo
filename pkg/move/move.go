package move

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	"path/filepath"

	_ "github.com/rclone/rclone/backend/local" // local backend and backblaze backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/fs/sync"
)

type moveRequest struct {
	SrcRemoteName string `json:"srcFs"`
	DstRemoteName string `json:"dstFs"`
	SrcFileName   string `json:"srcRemote"`
	DstFileName   string `json:"dstRemote"`
}

func MoveWithinRemote(client rpc_client.RPCClient, remoteName string, bucketName string, srcFile string, dstFile string) error {
	// TODO: Validation.
	// Usage:
	// remoteName is like backblaze:
	// bucketName is like LinuxFileTreeBackup
	// srcFile is like logs.tar.gz.enc.tmp
	// dstFile is like logs.tar.gz.enc
	return move(client, remoteName, remoteName, filepath.Join(bucketName, srcFile), filepath.Join(bucketName, dstFile))
}
func move(client rpc_client.RPCClient, srcRemoteName string, dstRemoteName string, srcFileName string, dstFileName string) error {
	if err := client.Initialize(); err != nil {
		return err
	}
	const moveMethod = "operations/movefile"
	moveRequest := moveRequest{
		SrcRemoteName: srcRemoteName,
		DstRemoteName: dstRemoteName,
		SrcFileName:   srcFileName,
		DstFileName:   dstFileName,
	}
	moveRequestJson, err := json.Marshal(moveRequest)
	if err != nil {
		return err
	}
	out, status := client.RPC(moveMethod, string(moveRequestJson))
	if status != 200 {
		return fmt.Errorf("Error status: %d and error output: %s", status, out)
	}
	return nil
}
