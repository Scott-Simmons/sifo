# Sifo: Securely Push and Pull to and from Backblaze

![Coverage Badge](coverage/coverage.svg)
![Build Status](https://img.shields.io/github/actions/workflow/status/Scott-Simmons/backup-system/ci.yml?branch=main)

`sifo` enables cloud backup and restore functionality using `backblaze` as the cloud storage provider, and `rclone` as the syncronisation tool.

`sifo` also implements archive (`*.tar.gz`) and AES-256 (CBC) application-level local encryption.

##### Dependencies

None. Unless you are [building from source](#building-from-source).


`librclone` exports [shims](https://github.com/rclone/rclone/blob/master/librclone/librclone/librclone.go) that wrap over the `rclone` RPC. Hence `sifo`'s' `rclone` dependencies are included at compile time.

##### System Requirements

`sifo` currently targets:

- `amd64` and `arm64` architectures.
- `linux`, `windows`, and `darwin` operating systems.

Note: Only `linux/amd64` and `darwin/arm64` have been tested through repeated usage.

`sifo` avoids system calls to support ease of data restoration.

##### Installation

Download the appropriate binary for your system [here](LINK TO RELEASES), or use `make install` to [build from source](#building-from-source).

##### Building from source

Requirements

- `go >= version 1.20`
- `make`

Build and install. Default install location is `usr/local/bin` but can be changed with `INSTALL_PREFIX`
```bash
sudo env "PATH=$PATH:/usr/local/go/bin" "GOPATH=$HOME/go" make install
sifo --version
```

##### Configuration

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
sifo gen-key > key.txt
```

### Usage:

Push a folder from local machine to backblaze remote:
```bash
sifo push \
    --src-dir=<folder_name> \
    --private-key=<key_path> \
    --bucket-name=<backblaze_bucket_name> \
    --remote-name=<remote_name:>
```

Pull a folder from backblaze remote to local:
```bash
sifo pull \
    --backblaze-remote-file-path=<remote_file_path> \
    --backblaze-remote-name=<remote_name:> \
    --backblaze-bucket-name=<backblaze_bucket_name> \
    --key-path=<key_path> \
    --dst-dir=<dir_path>
```

##### Testing

Much of the functionality is unit tested. More work can be done on the test suite.

End-to-end testing is not implemented. There are practical challenges involved setting up a suitable test environment for the end-to-end tests that can interact with a real `backblaze` cloud environment.

##### Notes:

Currently, `sifo` does not support dependencyless configuration encryption. If you want to encrypt your configuration, you can do so by downloading `rclone` and setting it up by following [the docs](https://rclone.org/docs/#configuration-encryption).

Example of using `sifo` using encrypted config:

Encrypting existing config
```bash
> rclone config
Current remotes
s) Set configuration password
```

Set the password.

Using `sifo` to push and pull folders. `RCLONE_CONFIG_PASS` must be exported to be able to read the encrypted `rclone` configuration file.
```bash
export RCLONE_CONFIG_PASS=<password>
sifo <commands>
```

##### Why application-layer encryption

`backblaze` provides server-side encryption. Similarly, `rclone` supports a [crypt](https://rclone.org/crypt/) remote that provides encryption to the remote.

The rationale for implementing encryption outside of `backblaze` & `rclone` is to have complete control over the encryption process, for an independent guarentee of end-to-end protection, regardless of the `rclone` & `backblaze` encryption implementations.

###### Why full backups

For simplicitiy and reliabilty of restorations, full backups were chosen over differential or incremental backups.

##### Versioning the backups

Versioning can be conifgured in `backblaze`. Read more in [the docs](https://www.backblaze.com/docs/cloud-storage-lifecycle-rules).

