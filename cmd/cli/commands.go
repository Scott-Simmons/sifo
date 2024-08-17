package main

// Probably exposing too much to the user right now. But can reel things back later

import (
	"SecureSyncDrive/pkg/archive_encrypt_sync_prune"
	"SecureSyncDrive/pkg/config_create"
	"SecureSyncDrive/pkg/config_dump"
	"SecureSyncDrive/pkg/decrypt"
	"SecureSyncDrive/pkg/encrypt"
	"SecureSyncDrive/pkg/sync"
	"fmt"
	"log"
)

// ref: https://github.com/alecthomas/kong/blob/master/_examples/docker/commands.go#L6

type GenEncryptedTarCmd struct {
	SrcDir     string `help:"Directory to archive."`
	DestFile   string `help:"File name for encrypted archive."`
	PrivateKey string `help:"AES-256 private key file path."`
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
	fmt.Println("Generating key...")
	key, err := encrypt.GenerateKey()
	if err != nil {
		log.Fatalf("Error generating private key: %v", err)
		fmt.Printf("Key generated, %v", key)
	}
	return nil
}

type SyncToBackblazeCmd struct {
	FilePathToSync   string `help:"File to sync to remote"`
	BackblazeRemoteName string `help:"Name of the (backblaze) remote"`
	// TODO: Maybe better to specify as a connection string for flexibilty
	// See: https://rclone.org/docs/#connection-strings
}

func (s *SyncToBackblazeCmd) Run(globals *Globals) error {
	fmt.Println("Syncing file to backblaze...")
	fmt.Println("Don't forget to export RCLONE_CONFIG_PASS")

	client, err := sync.NewClient()
	if err != nil {
		return err
	}
	err = sync.SyncToBackblaze(client, s.FilePathToSync, s.BackblazeRemoteName)
	if err != nil {
		return err
	}
	fmt.Println("Syncing done. Check the remote for changes.")
	return nil
}

type SyncFromBackblazeCmd struct {
	LocalSyncDir     string `help:"Directory to sync remote to"`
	BackblazeRemoteName string `help:"Name of the (backblaze) remote"`
	// TODO: Maybe better to specify as a connection string for flexibilty
	// See: https://rclone.org/docs/#connection-strings
}

func (s *SyncFromBackblazeCmd) Run(globals *Globals) error {
	fmt.Println("Syncing backblaze to local directory...")
	fmt.Println("Don't forget to export RCLONE_CONFIG_PASS")

	client, err := sync.NewClient()
	if err != nil {
		return err
	}
	err = sync.SyncFromBackblaze(client, s.BackblazeRemoteName, s.LocalSyncDir)
	if err != nil {
		return err
	}
	fmt.Println("Syncing done. Check the local for changes.")
	return nil
}

type DecryptTarCmd struct {
	SrcFile    string `help:"File to decrypt."`
	PrivateKey string `help:"AES-256 private key file path."`
}

func (d *DecryptTarCmd) Run() error {
	fmt.Println("Starting archival and encryption")
	err := decrypt.Decrypt(
		d.SrcFile,
		d.PrivateKey,
	)
	if err != nil {
		log.Fatalf("Error decrypting: %v", err)
	}
	fmt.Printf("File %s decrypted", d.SrcFile)
	return nil
}

type ConfigDumpCmd struct{}

func (c *ConfigDumpCmd) Run() error {
	fmt.Println("Dumping rclone config...")
	client, err := sync.NewClient()
	out, err := config_dump.DumpConfig(client)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", out)
	return nil
}

type ConfigCreateCmd struct{}

func (c *ConfigCreateCmd) Run() error {
	fmt.Println("creating rclone config...")
	client, err := sync.NewClient()
	err = config_create.CreateConfig(client, "hi", ":hi", true)
	if err != nil {
		return err
	}
	return nil
}
