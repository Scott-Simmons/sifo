package push

import (
	"SecureSyncDrive/pkg/archive"
	"SecureSyncDrive/pkg/delete"
	"SecureSyncDrive/pkg/encrypt"
	//	"SecureSyncDrive/pkg/move"
	"SecureSyncDrive/pkg/sync"
	"fmt"
	"os"
	"path/filepath"
)

func Push(srcDir string, bucketName string, keyPath string, remoteName string) error {

	archivedFileDest := srcDir + ".tar.gz"
	if archiveErr := archive.ArchiveFolder(srcDir, archivedFileDest); archiveErr != nil {
		localArchiveRemoveErr := os.Remove(archivedFileDest)
		return fmt.Errorf(
			"Local archive error: %v; Local Cleanup error: %v",
			archiveErr, localArchiveRemoveErr,
		)
	}

	encryptedFileDest := archivedFileDest + ".enc"
	if encryptErr := encrypt.EncryptTarBall(archivedFileDest, encryptedFileDest, keyPath); encryptErr != nil {
		localArchiveRemoveErr := os.Remove(archivedFileDest)
		localEncryptedRemoveErr := os.Remove(encryptedFileDest)
		return fmt.Errorf(
			"Local encrypt error: %v; Local archive cleanup error: %v; Local encrypted cleanup error: %v",
			encryptErr, localArchiveRemoveErr, localEncryptedRemoveErr,
		)
	}

	if localArchiveRemoveErr := os.Remove(archivedFileDest); localArchiveRemoveErr != nil {
		return fmt.Errorf("Local archive cleanup error: %v", localArchiveRemoveErr)
	}

	// Reference idiom: https://stackoverflow.com/questions/37932551/mkdir-if-not-exists-using-golang
	//tempRemoteDestName := encryptedFileDest + ".tmp"
	tempRemoteDestName := encryptedFileDest // the function "move" is borked
	if directoryCreateErr := os.Mkdir(bucketName, 0700); directoryCreateErr != nil && !os.IsExist(directoryCreateErr) {
		localEncryptedRemoveErr := os.Remove(encryptedFileDest)
		return fmt.Errorf(
			"Local dir creation error: %v; Local encrypted cleanup error: %v",
			directoryCreateErr, localEncryptedRemoveErr,
		)
	}
	if fileMoveError := os.Rename(encryptedFileDest, filepath.Join(bucketName, tempRemoteDestName)); fileMoveError != nil {
		localEncryptedRemoveErr := os.Remove(encryptedFileDest)
		return fmt.Errorf(
			"Local file move error: %v; Local encrypted cleanup error: %v",
			fileMoveError, localEncryptedRemoveErr,
		)
	}

	client, clientError := sync.NewClient()
	if clientError != nil {
		localRemoveDirErr := os.RemoveAll(bucketName)
		return fmt.Errorf(
			"Client creation error: %v; Local bucket cleanup dir error: %v",
			clientError, localRemoveDirErr,
		)
	}

	if remoteSyncErr := sync.SyncToBackblaze(client, bucketName, remoteName, bucketName); remoteSyncErr != nil {
		localRemoveErr := os.RemoveAll(bucketName)
		remoteDeleteErr := delete.DeleteBackBlazeFile(client, remoteName, bucketName, tempRemoteDestName)
		return fmt.Errorf(
			"Remote sync error: %v; Local cleanup error: %v; Remote cleanup error: %v",
			remoteSyncErr, localRemoveErr, remoteDeleteErr,
		)
	}

	//if remoteMoveErr := move.MoveWithinRemote(client, remoteName, bucketName, tempRemoteDestName, encryptedFileDest); remoteMoveErr != nil {
	//localRemoveErr := os.RemoveAll(bucketName)
	//remoteDeleteErr := delete.DeleteBackBlazeFile(client, remoteName, bucketName, tempRemoteDestName)
	//remoteDeleteErr := "nah"
	//return fmt.Errorf(
	//	"Remote move error: %v; Local cleanup error: %v; Remote cleanup error: %v",
	//	remoteMoveErr, localRemoveErr, remoteDeleteErr,
	//)
	//}

	if localRemoveErr := os.RemoveAll(bucketName); localRemoveErr != nil {
		return localRemoveErr
	}

	return nil
}
