package main

// Probably exposing too much to the user right now. But can reel things back later

import (
	"SecureSyncDrive/pkg/config_create"
	"SecureSyncDrive/pkg/config_dump"
	"SecureSyncDrive/pkg/decrypt"
	"SecureSyncDrive/pkg/do_copy"
	"SecureSyncDrive/pkg/encrypt"
	"SecureSyncDrive/pkg/push"
	"SecureSyncDrive/pkg/sync"
	"fmt"
	"log"
)

// ref: https://github.com/alecthomas/kong/blob/master/_examples/docker/commands.go#L6
type Push struct {
	SrcDir     string `help:"Directory to archive and encrypt."`
	PrivateKey string `help:"AES-256 private key file path."`
	BucketName string `help:"Name of the backblaze bucket to write directory to."`
	RemoteName string `help:"Name of the backblaze remote."`
}

func (g *Push) Run() error {
	fmt.Println("Starting archival, encryption, and syncing to backblaze.")
	err := push.Push(
		g.SrcDir,
		g.BucketName,
		g.PrivateKey,
		g.RemoteName,
	)
	if err != nil {
		log.Fatalf("Error archiving and encrypting: %v", err)
	}
	fmt.Printf("Directory %s archived and encrypted to backblaze remote %s in bucket %s\n", g.SrcDir, g.RemoteName, g.BucketName)
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
	FilePathToSync      string `help:"File to sync to remote."`
	BackblazeRemoteName string `help:"Name of the (backblaze) remote."`
	BackblazeBucketName string `help:"Name of the (backblaze) bucket."`
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
	err = sync.SyncToBackblaze(client, s.FilePathToSync, s.BackblazeRemoteName, s.BackblazeBucketName)
	if err != nil {
		return err
	}
	fmt.Println("Syncing done. Check the remote for changes.")
	return nil
}

type CopyFromBackblazeCmd struct {
	LocalSyncDst            string `help:"Directory name to sync remote to"`
	BackblazeRemoteFilePath string `help:"Name of the (backblaze) remote file"`
	BackblazeRemoteName     string `help:"Name of the (backblaze) remote"`
	BackblazeBucketName     string `help:"Name of the (backblaze) bucket"`
	// TODO: Maybe better to specify as a connection string for flexibilty
	// See: https://rclone.org/docs/#connection-strings
}

func (c *CopyFromBackblazeCmd) Run(globals *Globals) error {
	fmt.Println("Syncing backblaze to local directory...")
	fmt.Println("Don't forget to export RCLONE_CONFIG_PASS")
	client, err := sync.NewClient()
	if err != nil {
		return err
	}
	err = do_copy.CopyFromBackblaze(client, c.BackblazeRemoteFilePath, c.BackblazeRemoteName, c.BackblazeBucketName, c.LocalSyncDst)
	if err != nil {
		return err
	}
	fmt.Println("Syncing done. Check local for changes.")
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
