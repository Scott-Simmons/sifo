package archive_encrypt_sync_prune

import (
  "fmt"
	"SecureSyncDrive/pkg/archive"
	"SecureSyncDrive/pkg/encrypt"
	// "SecureSyncDrive/pkg/sync"
	// "SecureSyncDrive/pkg/prune"
)

func ArchiveEncryptSyncPrune(srcDir string, destFile string, keyPath string) error {
  // This will orchestrate the whole process. For now just archive and encrypt

  tempDestFile := "temp_" + destFile

  err := archive.ArchiveFolder(srcDir, tempDestFile)
  if err != nil {
    // should I do cleanup? Yes probably
    // TODO: Delete archive folder etc... make sure state is cleaned up
    return err
  }
  fmt.Println("Archive created successfully")

  err = encrypt.EncryptTarBall(tempDestFile, destFile, keyPath)
  if err != nil {
    // should I do cleanup...
    return err
  }
  fmt.Println("Archive encrypted successfully")

  // TODO: Sync to the remote.

  // TODO: Prune the revisions on the remote.



  return nil

}
