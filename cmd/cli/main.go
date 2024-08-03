package main

// Ref: https://github.com/alecthomas/kong
// Using the docker example: https://github.com/alecthomas/kong/blob/master/_examples/docker/main.go

import (
  "fmt"
	"github.com/alecthomas/kong"
)

type Globals struct {
  RcloneConfig    string    `help:"Location of the client rclone config file" default:"~/.config/rclone/rclone.conf" type:"path"`
  Version         VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type CLI struct {

  Globals

  GenKey GenKeyCmd `cmd:"" help:"Generate a key Using AES-256."`
  GenEncryptedTar GenEncryptedTarCmd `cmd:"" help:"Archive and encrypt a directory."`
  SyncToRemote SyncToGoogleDriveCmd `cmd:"" help:"Sync file to google drive remote."`
  DecryptTar DecryptTarCmd `cmd:"" help:"Decrypt an encrypted archive."`
}

// ref: https://github.com/alecthomas/kong/blob/master/_examples/shell/commandstring/main.go
// ref: https://github.com/alecthomas/kong/blob/master/_examples/docker/main.go
func main() {

  cli := CLI{
    Globals: Globals{
      Version: VersionFlag("0.0.1"),
    },
  }
  
  ctx := kong.Parse(&cli,
          kong.Name("SecureStoreSync"),
          kong.Description("CLI for host-level archival and encryption and sync to a google drive remote."),
          kong.UsageOnError(),
          kong.ConfigureHelp(kong.HelpOptions{
                  Compact: true,
                  Summary: true,
          }),
          kong.Vars{
                  "version": "0.0.1",
          })

  err := ctx.Run(&cli.Globals)
  ctx.FatalIfErrorf(err)
}

