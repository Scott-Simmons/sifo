package main

import (
	"SecureSyncDrive/pkg/config_dump"
	"SecureSyncDrive/pkg/encrypt"
	"SecureSyncDrive/pkg/pull"
	"SecureSyncDrive/pkg/push"
	"SecureSyncDrive/pkg/sync"
	"SecureSyncDrive/pkg/validate_config"
	"fmt"
	"log"
)

// Exposing only: Push, Pull, GenKey, DumpConfig, ValidateConfig
type PullCmd struct {
	DstDir                  string `help:"Name of local directory to dump out extracted remote files to."`
	BackblazeRemoteFilePath string `help:"Name of the (backblaze) remote file"`
	BackblazeRemoteName     string `help:"Name of the (backblaze) remote"`
	BackblazeBucketName     string `help:"Name of the (backblaze) bucket"`
	KeyPath                 string `help:"Path to the symmetric encryption key"`
}

func (p *PullCmd) Run(globals *Globals) error {
	fmt.Println("Pulling backblaze to local directory...")
	fmt.Println("Don't forget to export RCLONE_CONFIG_PASS")
	err := pull.Pull(p.DstDir, p.BackblazeRemoteFilePath, p.BackblazeBucketName, p.KeyPath, p.BackblazeRemoteName)
	if err != nil {
		return err
	}
	fmt.Println("Pulling done. Check local for changes.")
	return nil
}

// ref: https://github.com/alecthomas/kong/blob/master/_examples/docker/commands.go#L6
type PushCmd struct {
	SrcDir     string `help:"Directory to archive and encrypt."`
	PrivateKey string `help:"AES-256 private key file path."`
	BucketName string `help:"Name of the backblaze bucket to write directory to."`
	RemoteName string `help:"Name of the backblaze remote."`
}

func (g *PushCmd) Run() error {
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

type GenKeyCmd struct{}

func (g *GenKeyCmd) Run() error {
	key, err := encrypt.GenerateKey()
	if err != nil {
		log.Fatalf("Error generating private key: %v", err)
	}
	fmt.Println(key)
	return nil
}

type ConfigDumpCmd struct{}

func (c *ConfigDumpCmd) Run() error {
	client, err := sync.NewClient()
	out, err := config_dump.DumpConfig(client)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", out)
	return nil
}

type ConfigValidateCmd struct {
	ConfigPath string `help:"Path to the config"`
}

func (c *ConfigValidateCmd) Run() error {
	fmt.Println("Validating rclone config...")
	err := validate_config.ValidateConfig(c.ConfigPath)
	if err != nil {
		return err
	}
	return nil
}
