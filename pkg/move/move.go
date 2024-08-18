package move

// TODO: This file breaks with 401s

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	//"path/filepath"

	_ "github.com/rclone/rclone/backend/b2"    // local backend and backblaze backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/backend/local" // local backend and backblaze backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/fs/sync"
)

// TODO: Running asynchronous jobs with _async = true
type moveRequest struct {
	SrcFs     string `json:"srcFs"`
	SrcRemote string `json:"srcRemote"`
	DstFs     string `json:"dstFs"`
	DstRemote string `json:"dstRemote"`
}

func MoveWithinRemote(client rpc_client.RPCClient, remoteName string, bucketName string, srcFile string, dstFile string) error {
	srcFileName := bucketName + "/" + srcFile
	dstFileName := bucketName + "/" + dstFile
	return move(client, remoteName, remoteName, srcFileName, dstFileName)
}
func move(client rpc_client.RPCClient, srcRemoteName string, dstRemoteName string, srcFileName string, dstFileName string) error {
	if err := client.Initialize(); err != nil {
		return err
	}
	const moveMethod = "operations/movefile"
	moveRequest := moveRequest{
		SrcFs:     srcRemoteName,
		SrcRemote: srcFileName,
		DstFs:     dstRemoteName,
		DstRemote: dstFileName,
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
