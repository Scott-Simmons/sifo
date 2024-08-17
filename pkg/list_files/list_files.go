package list_files

import (
	"SecureSyncDrive/pkg/rpc_client"
	"encoding/json"
	"fmt"
	_ "github.com/rclone/rclone/fs/operations"
	"log"
)

type listFilesRequest struct {
	FSrc   string `json:"fs"`
	Remote string `json:"remote"`
}

type BackblazeFile struct {
	Path     string `json:"Path"`
	Name     string `json:"Name"`
	Size     int    `json:"Size"`
	MimeType string `json:"MimeType"`
	ModTime  string `json:"ModTime"`
	IsDir    bool   `json:"IsDir"`
	ID       string `json:"ID"`
}

type FileListOutput struct {
	List []BackblazeFile `json:"list"`
}

// Generics don't work well in golang but would be lovely to have a generic map.
func Map(files []BackblazeFile, f func(BackblazeFile) string) []string {
	ids := make([]string, len(files))
	for i, file := range files {
		ids[i] = f(file)
	}
	return ids
}

func ListRemoteFiles(client rpc_client.RPCClient, BackblazeRemote string) ([]BackblazeFile, error) {
	// Assume flat structure for now no recursive list. Will be flat encrypted anyway.
	// Ref: https://github.com/rclone/rclone/blob/1901bae4ebcbc4cdd82f6bb1862d0479f3fa386e/fs/operations/rc_test.go
	if err := client.Initialize(); err != nil {
		return nil, err
	}
	const listJsonMethod = "operations/list"
	listFilesRequest := listFilesRequest{
		FSrc:   BackblazeRemote,
		Remote: "",
	}
	listFilesRequestJson, err := json.Marshal(listFilesRequest)
	if err != nil {
		return nil, err
	}
	out, status := client.RPC(listJsonMethod, string(listFilesRequestJson))
	if status != 200 {
		return nil, fmt.Errorf("Error status: %d and error output: %s", status, out)
	}
	var structuredOutput FileListOutput
	err = json.Unmarshal([]byte(out), &structuredOutput)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	return structuredOutput.List, nil
}

func GetBackblazeFileIds(files []BackblazeFile) []string {
	return Map(files, func(f BackblazeFile) string { return f.ID })
}
