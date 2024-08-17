package config_create

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
)

// Ref: https://github.com/rclone/rclone/blob/master/backend/pikpak/api/types.go#L253
type CreateConfigStructure struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Parameters struct {
		TestKey string `json:"test_key"`
	} `json:"parameters"`
}

// Ref: https://rclone.org/commands/rclone_config_create/
// Ref: https://www.youtube.com/watch?v=f8K-V3HHDA0

// Need to also not do auto-config
// Need to be able to init everything without issues...
// Basically, just want to pass everything in from a json config file to set up config except for the password.

func CreateConfig(client rpc_client.RPCClient, remoteName string, remoteType string, overwriteIfExists bool) error {
	// PSUEDOCODE....
	// 1. Check if remote name exists and if it does, maybe overwrite.
	// 2. Validate config name with MakeConfigName which wil lmake sure the new remote name string is legal.
	// 3.
	// Important password ref: https://github.com/rclone/rclone/blob/master/fs/config/ui.go#L780
	// Also useful: https://github.com/rclone/rclone/blob/master/fs/config/crypt.go#L245

	if err := client.Initialize(); err != nil {
		return err
	}
	const configCreateMethod = "config/create"
	const configPasswordMethod = "config/password"
	configCreateRequest := CreateConfigStructure{}
	configCreateRequestJson, err := json.Marshal(configCreateRequest)
	if err != nil {
		return err
	}
	out, status := client.RPC(configCreateMethod, string(configCreateRequestJson))
	if status != 200 {
		return fmt.Errorf("Error status: %d and error output: %s", status, out)
	}
	return nil
}
