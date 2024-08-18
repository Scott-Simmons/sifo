package pull

import (
	"SecureSyncDrive/pkg/decrypt"
	"SecureSyncDrive/pkg/do_copy"
	"SecureSyncDrive/pkg/extract"
	"SecureSyncDrive/pkg/sync"
	"fmt"
	"os"
)

func Pull(dstDirectory string, remoteFileName string, bucketName string, keyPath string, remoteName string) error {

	client, clientError := sync.NewClient()
	if clientError != nil {
		localRemoveDirErr := os.RemoveAll(bucketName)
		return fmt.Errorf(
			"Client creation error: %v; Local bucket cleanup dir error: %v",
			clientError, localRemoveDirErr,
		)
	}

	if remoteCopyToLocalErr := do_copy.CopyFromBackblaze(client, remoteFileName, remoteName, bucketName, remoteFileName); remoteCopyToLocalErr != nil {
		localRemoteDirErr := os.RemoveAll(remoteFileName)
		return fmt.Errorf(
			"Copy to local error: %v; Local cleanup error: %v",
			remoteCopyToLocalErr, localRemoteDirErr,
		)
	}

	aPrioriDecryptedFileName := remoteFileName + ".dec"
	if decryptErr := decrypt.Decrypt(remoteFileName, keyPath); decryptErr != nil {
		localRemoveErr := os.RemoveAll(remoteFileName)
		localRemoveEncErr := os.RemoveAll(aPrioriDecryptedFileName)
		return fmt.Errorf(
			"Decrypt to local error: %v; Local cleanup errors: %v; %v",
			decryptErr, localRemoveErr, localRemoveEncErr,
		)
	}
	os.RemoveAll(remoteFileName)

	r, err := os.Open(aPrioriDecryptedFileName)
	if err != nil {
		return err
	}
	defer r.Close()

	extractDataErr := extract.ExtractTar(aPrioriDecryptedFileName, dstDirectory)
	if extractDataErr != nil {
		return extractDataErr
	}
	os.RemoveAll(aPrioriDecryptedFileName)

	return nil
}
