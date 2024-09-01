package push

import (
	"SecureSyncDrive/pkg/archive"
	"SecureSyncDrive/pkg/delete"
	"SecureSyncDrive/pkg/encrypt"
	"SecureSyncDrive/pkg/sync"
  "path/filepath"
	"fmt"
	"os"
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

	// TODO: See issue: https://github.com/Scott-Simmons/backup-system/issues/23 ... errors are thrown when using `move`. This should be implemented in the future but it is non-essential. It ensures the file transfers are 100% successful before overwriting.

	if directoryCreateErr := os.Mkdir(bucketName, 0700); directoryCreateErr != nil && !os.IsExist(directoryCreateErr) {
		localEncryptedRemoveErr := os.Remove(encryptedFileDest)
		return fmt.Errorf(
			"Local dir creation error: %v; Local encrypted cleanup error: %v",
			directoryCreateErr, localEncryptedRemoveErr,
		)
	}

  if fileMoveError := os.Rename(encryptedFileDest, filepath.Join(bucketName, filepath.Base(encryptedFileDest))); fileMoveError != nil {
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
		remoteDeleteErr := delete.DeleteBackBlazeFile(client, remoteName, bucketName, encryptedFileDest)
		return fmt.Errorf(
			"Remote sync error: %v; Local cleanup error: %v; Remote cleanup error: %v",
			remoteSyncErr, localRemoveErr, remoteDeleteErr,
		)
	}

	if localRemoveErr := os.RemoveAll(bucketName); localRemoveErr != nil {
		return localRemoveErr
	}

	return nil
}
