# Sifo: Securely Push and Pull to and from Backblaze
![Build Status](https://img.shields.io/github/actions/workflow/status/Scott-Simmons/backup-system/ci.yml?branch=main)

`sifo` enables cloud backup and restore functionality using `backblaze` as the cloud storage provider, and `rclone` as the syncronisation tool.

`sifo` also implements archive (`*.tar.gz`) and AES-256 (CBC) application-level local encryption.

## Dependencies

None. Unless you are [building from source](#building-from-source).

## System Requirements

`sifo` currently targets:

- `amd64` and `arm64` architectures.
- `linux`, `windows`, and `darwin` operating systems.

Note: Only `linux/amd64` and `darwin/arm64` have been tested through repeated usage.

## Installation

Download the appropriate binary for your system [here](https://github.com/Scott-Simmons/backup-system/releases/), or use `make install` to [build from source](#building-from-source).

### Building from source

Requirements

- `go >= version 1.20`
- `make`

Build and install. The default install location is `usr/local/bin`. This can be configured with `INSTALL_PREFIX`
```bash
make build
sudo env "PATH=$PATH:/usr/local/go/bin" "GOPATH=$HOME/go" make install
sifo --version
```

## Uninstalling

```bash
sudo make uninstall
```

## Configuration

To set the rclone configuration file with a backblaze remote, create a file named `~/.config/rclone/rclone.conf`.

```ini
[backblaze]
type = b2
account = <bucket account name here>
key = <bucket account key here>
hard_delete = true
```

Get the bucket account value and key by creating a b2 bucket [here](https://secure.backblaze.com).

Validate the rclone configuration file.
```bash
sifo config-validate --config-path=~/.config/rclone/rclone.conf           
```

Generate a AES-256 encryption key. Keep this secure.
```bash
sifo gen-key > ~/.sifo_key
```

## Usage:

Example of a folder called `/home/foo/Documents/Notes` for a bucket called `FileBackups` on a remote called `backblaze`

Push a folder from local machine to backblaze remote:
```bash
sifo push \
    --src-dir=/home/foo/Documents/Notes/ \
    --private-key=~/.sifo_key \
    --bucket-name=FileBackups \
    --remote-name=backblaze:
```

Pull a folder from backblaze remote to local. The remote directory will be compressed and suffixed with `tar.gz.enc`. In this example, the files will be restored to `/user/home/restored_notes`.
```bash
sifo pull \
    --backblaze-remote-file-path=Notes.tar.gz.enc \
    --backblaze-remote-name=backblaze: \
    --backblaze-bucket-name=FileBackups \
    --key-path=~/.sifo_key \
    --dst-dir=/home/foo/restored_notes
```

## Testing

Much of the functionality is unit tested. More work can be done on the test suite.

End-to-end testing is not implemented. There are practical challenges involved setting up a suitable test environment for the end-to-end tests that can interact with a real `backblaze` cloud environment.

## Notes:

### Encrypted `rclone` config

Currently, `sifo` does not support dependencyless configuration encryption. If you want to encrypt your configuration, you can do so by downloading `rclone` and setting it up by following [the docs](https://rclone.org/docs/#configuration-encryption).

Once the config is encrypted, `sifo` can be used push and pull folders, provided that `RCLONE_CONFIG_PASS` is exported.

### Versioning the backups

Versioning can be conifgured in `backblaze`. Read more in [the docs](https://www.backblaze.com/docs/cloud-storage-lifecycle-rules).

### Rationale for application-layer encryption

`backblaze` provides server-side encryption. Similarly, `rclone` supports a [crypt](https://rclone.org/crypt/) remote that provides encryption to the remote.

The rationale for implementing encryption outside of `backblaze` & `rclone` is to have complete control over the encryption process, independent from the `rclone` & `backblaze` encryption implementations.

### Rationale for full backups

For simplicitiy and reliabilty of restorations, full backups were chosen over differential or incremental backups.

### Implementation details

`sifo` is fully statically linked with no runtime dependencies except for the optional [encrypted rclone config](#encrypted-rclone-config)

`librclone` exports [shims](https://github.com/rclone/rclone/blob/master/librclone/librclone/librclone.go) that wrap over the `rclone` RPC. Hence `sifo`'s' `rclone` dependencies are included at compile time.
