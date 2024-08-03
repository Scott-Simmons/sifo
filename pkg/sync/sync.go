package sync

// Ref: https://rclone.org/docs/#connection-strings
// Ref: https://forum.rclone.org/t/use-rclone-golang-to-transfer-files/34983/23?page=2
// Ref: https://github.com/alankritkharbanda/rclone/blob/5975b7d27728f5ba0c3c670759fe9cc3dfb65ff2/librclone/README.md

import (
	"encoding/json"
	"fmt"
	_ "github.com/rclone/rclone/fs/sync"
	"github.com/rclone/rclone/librclone/librclone"
  _ "github.com/rclone/rclone/backend/local" // local backend and google drive backend imported to ensure they are registered. This results in 
  _ "github.com/rclone/rclone/backend/drive"
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

type syncRequest struct {
  SrcFs     string    `json:"srcFs"`
  DstFs     string    `json:"dstFs"`  // best to keep struct as PascalCase. Not accessible otherwise
}

func NewClient() (*RealRPCClient, error) {
  client := &RealRPCClient{}
  return client, nil
}

func SyncGoogleDrive(client RPCClient, srcFilePath string, googleDriveFilePath string) error {
  if err := client.Initialize(); err != nil {
    return err
  }

  // Ref: https://github.com/rclone/rclone/blob/master/fs/sync/sync.go
  const syncMethod = "sync/sync"

  syncRequest := syncRequest {
    SrcFs: srcFilePath,
    DstFs: googleDriveFilePath,
  }
  fmt.Printf("Sync re: %+v\n", syncRequest)

  syncRequestJson, err := json.Marshal(syncRequest)
  if err != nil {
    return err
  }
  fmt.Println(syncRequestJson)
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

