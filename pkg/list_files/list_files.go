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

type GoogleDriveFile struct {
	Path     string `json:"Path"`
	Name     string `json:"Name"`
	Size     int    `json:"Size"`
	MimeType string `json:"MimeType"`
	ModTime  string `json:"ModTime"`
	IsDir    bool   `json:"IsDir"`
	ID       string `json:"ID"`
}

type FileListOutput struct {
	List []GoogleDriveFile `json:"list"`
}

// Generics don't work well in golang but would be lovely to have a generic map.
func Map(files []GoogleDriveFile, f func(GoogleDriveFile) string) []string {
	ids := make([]string, len(files))
	for i, file := range files {
		ids[i] = f(file)
	}
	return ids
}

func ListRemoteFiles(client rpc_client.RPCClient, googleDriveRemote string) ([]GoogleDriveFile, error) {
	// Assume flat structure for now no recursive list. Will be flat encrypted anyway.
	// Ref: https://github.com/rclone/rclone/blob/1901bae4ebcbc4cdd82f6bb1862d0479f3fa386e/fs/operations/rc_test.go
	if err := client.Initialize(); err != nil {
		return nil, err
	}
	const listJsonMethod = "operations/list"
	listFilesRequest := listFilesRequest{
		FSrc:   googleDriveRemote,
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

func GetGoogleDriveFileIds(files []GoogleDriveFile) []string {
	return Map(files, func(f GoogleDriveFile) string { return f.ID })
}
