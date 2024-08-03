package main

import (
  "fmt"
  "SecureSyncDrive/pkg/archive_encrypt_sync_prune"
  "SecureSyncDrive/pkg/encrypt"
  "log"
)

// ref: https://github.com/alecthomas/kong/blob/master/_examples/docker/commands.go#L6

type GenEncryptedTarCmd struct {
	SrcDir      string    `help:"Directory to archive."`
	DestFile    string    `help:"File name for encrypted archive."`
	PrivateKey  string    `help:"AES-256 private key file path."`
}

func (g *GenEncryptedTarCmd) Run() error {
  fmt.Println("Starting archival and encryption")
  err := archive_encrypt_sync_prune.ArchiveEncryptSyncPrune(
    g.SrcDir,
    g.DestFile,
    g.PrivateKey,
  )
  if err != nil {
    log.Fatalf("Error archiving and encrypting: %v", err)
  }
  fmt.Printf("Directory %s archived and encrypted to %s\n", g.SrcDir, g.DestFile)
  return nil
}


type GenKeyCmd struct {
  // Just to keep the pattern consistent.
}

func (g *GenKeyCmd) Run() error {
  fmt.Printf("Generating key...")
  key, err := encrypt.GenerateKey()
  if err != nil {
    log.Fatalf("Error generating private key: %v", err)
  fmt.Printf("Key generated, %v", key)
  }
  return nil
}

