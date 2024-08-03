package cmd

import (
	"SecureSyncDrive/pkg/archive"
	"SecureSyncDrive/pkg/encrypt"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
  srcDir := flag.String("src", "", "Source directory to archive and encrypt")
  destTar := flag.String("dest", "", "Destination path for the encrypted tarball")
  privateKeyPath := flag.String("key", "", "Private key path to Encrypt with (only AES-256 supported). Generate with GenerateKey()")
  flag.Parse()

  if *srcDir == "" || *destTar == "" {
      fmt.Println("Source directory and destination path must be specified")
      flag.Usage()
      os.Exit(1)
  }

  intermediateTar := "tmp_" + *destTar

  err := archive.ArchiveFolder(*srcDir, intermediateTar)
  if err != nil {
      log.Fatalf("Error archiving folder: %v", err)
  }
  fmt.Println("Archive created successfully")

  err = encrypt.EncryptTarBall(intermediateTar, *destTar, *privateKeyPath)
  if err != nil {
    log.Fatalf("Error Encrypting folder: %v", err)
  }
  fmt.Println("Archive encrypted successfully")

  // TODO: Rclone sync


}
