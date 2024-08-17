package push

import (
	"SecureSyncDrive/pkg/archive"
	"SecureSyncDrive/pkg/encrypt"
	"path/filepath"
	"SecureSyncDrive/pkg/sync"
	"fmt"
	"os"
)

func Push(srcDir string, bucketName string, keyPath string, remoteName string) error {
  // Archive, Encrypt, Sync to remote.

	archivedFileDest := srcDir + ".tar.gz"
	err := archive.ArchiveFolder(srcDir, archivedFileDest)
	if err != nil {
		// should I do cleanup? Yes. Maybe needs to be cleaned up inside function...
		// TODO: Delete archive folder etc... make sure state is cleaned up
		return err
	}
	fmt.Println("Archive created successfully")

  encryptedFileDest := archivedFileDest + ".enc"
	err = encrypt.EncryptTarBall(archivedFileDest, encryptedFileDest, keyPath)
	if err != nil {
		// Do cleanup... maybe do inside function itself.
		return err
	}
	fmt.Println("Archive encrypted successfully")
  os.Remove(archivedFileDest)

  // https://stackoverflow.com/questions/37932551/mkdir-if-not-exists-using-golang
  if err := os.Mkdir(bucketName, 0700); err != nil && !os.IsExist(err) {
    return err
  }
  err = os.Rename(encryptedFileDest, filepath.Join(bucketName, encryptedFileDest))
  if err != nil {
    return err
  }

  // TODO: This is in urgent need of becoming transactional. e.g. It needs to do a rename, sync, if success rename. If fail then delete and raise error. Very important if I am syncing large files.
  client, err := sync.NewClient()
  if err != nil {
    return err
  }

  err = sync.SyncToBackblaze(client, bucketName, remoteName, bucketName)
  if err != nil {
    os.RemoveAll(bucketName)
    return err
  }
	fmt.Println("Synced to remote successfully")
  os.RemoveAll(bucketName)

	return nil
}
