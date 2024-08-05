package config_dump

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	"log"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
}

type GoogleDriveBackup struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	TeamDrive    string `json:"team_drive"`
	Token        Token  `json:"-"`     // this is a hack.
	TokenString  string `json:"token"` // this is a hack.
	Type         string `json:"type"`
}

type RcloneConfig struct {
	GoogleDriveBackup GoogleDriveBackup `json:"google-drive-backup"`
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
	// Need to unmarshal the token separately... not sure a better way exists
	var token Token
	err = json.Unmarshal([]byte(structuredOutput.GoogleDriveBackup.TokenString), &token)
	if err != nil {
		fmt.Println("Error unmarshalling token JSON:", err)
		return RcloneConfig{}, err
	}
	structuredOutput.GoogleDriveBackup.Token = token
	return structuredOutput, nil
}

// use *Config to avoid copying large struct...and avoid unneccassry mem allocations
// passing into a function copys the entire struct
func GetGoogleDriveAccessToken(config *RcloneConfig) string {
	return config.GoogleDriveBackup.Token.AccessToken
}
