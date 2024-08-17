package do_copy

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	"path/filepath"

	_ "github.com/rclone/rclone/backend/local" // local backend and backblaze backend imported to ensure they are registered. Otherwise errors out.
	_ "github.com/rclone/rclone/fs/sync"
)

type copyRequest struct {
	SrcFs       string `json:"srcFs"`
	DstFs       string `json:"dstFs"` // best to keep struct as PascalCase. Not accessible otherwise
	SrcFileName string `json:"srcRemote"`
	DstFileName string `json:"dstRemote"`
}

func NewClient() (*rpc_client.RealRPCClient, error) {
	client := &rpc_client.RealRPCClient{}
	return client, nil
}
func do_copy(client rpc_client.RPCClient, srcFs string, dstFs string, srcFileName string, dstFileName string) error {
	if err := client.Initialize(); err != nil {
		return err
	}
	// Ref: https://github.com/rclone/rclone/blob/master/fs/sync/sync.go
	const copyMethod = "operations/copyfile"
	copyRequest := copyRequest{
		SrcFs:       srcFs,
		DstFs:       dstFs,
		SrcFileName: srcFileName,
		DstFileName: dstFileName,
	}
	copyRequestJson, err := json.Marshal(copyRequest)
	if err != nil {
		return err
	}
	out, status := client.RPC(copyMethod, string(copyRequestJson))
	if status != 200 {
		return fmt.Errorf("Error status: %d and error output: %s", status, out)
	}
	return nil
}
func CopyFromBackblaze(client rpc_client.RPCClient, backblazeRemoteFilePath string, backblazeRemoteName string, backblazeBucketName string, localDstDir string) error {
	return do_copy(client, backblazeRemoteName, "", filepath.Join(backblazeBucketName, backblazeRemoteFilePath), localDstDir)
}
