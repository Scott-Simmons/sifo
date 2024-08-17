package config_dump

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	"log"
)

type BackblazeBackup struct {
	Account    string `json:"account"`
	HardDelete string `json:"hard_delete"`
	Key        string `json:"key"`
	Type       string `json:"type"`
}

type RcloneConfig struct {
	BackblazeBackup BackblazeBackup `json:"backblaze"`
}

type configDumpRequest struct{}

// TODO: should move this
func NewClient() (*rpc_client.RealRPCClient, error) {
	client := &rpc_client.RealRPCClient{}
	return client, nil
}

func DumpConfig(client rpc_client.RPCClient) (RcloneConfig, error) {
	if err := client.Initialize(); err != nil {
		return RcloneConfig{}, err
	}
	const configDumpMethod = "config/dump"
	configDumpRequest := configDumpRequest{}
	configDumpRequestJson, err := json.Marshal(configDumpRequest)
	if err != nil {
		return RcloneConfig{}, err
	}
	out, status := client.RPC(configDumpMethod, string(configDumpRequestJson))
	if status != 200 {
		return RcloneConfig{}, fmt.Errorf("Error status: %d and error output: %s", status, out)
	}
	var structuredOutput RcloneConfig
	err = json.Unmarshal([]byte(out), &structuredOutput)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
	return structuredOutput, nil
}
