# Sifo: Securely Push and Pull to and from Backblaze

![Coverage Badge](coverage/coverage.svg)
![Build Status](https://img.shields.io/github/actions/workflow/status/Scott-Simmons/backup-system/ci.yml?branch=main)

`sifo` enables cloud backup and restore functionality using `backblaze` as the cloud storage provider, and `rclone` as the syncronisation tool.

`sifo` also implements archive (`*.tar.gz`) and AES-256 (CBC) encryption at the application level on the local server.

##### Dependencies

None. Unless you are [building from source]().

`librclone` exports [shims](https://github.com/rclone/rclone/blob/master/librclone/librclone/librclone.go) that wrap over the `rclone` RPC. `sifo`'s' `rclone` dependencies are included at compile time.

##### System Requirements

`sifo` currently targets:

- the mainstream linux platforms: foo bar baz
- TODO: Cross compile......... 

##### Installation

Download the appropriate binary for your system, or [build from source]().

##### Building from source

Requirements

- `go >= version 1.20`
- `make`

```bash
make install
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

Generating the AES-256 encryption key
```bash
sifo gen-key > key.txt
```

### Usage:

Pushing a folder from local machine to backblaze remote:
```bash
sifo push \
    --src-dir=<folder_name> \
    --private-key=<key_path> \
    --bucket-name=<backblaze_bucket_name> \
    --remote-name=<remote_name:>
```

Pulling a folder from backblaze remote to local:
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

Currently, end-to-end testing is not implemented. This is due to the challenges in setting up a suitable test environment for the end-to-end tests to interact with a real `backblaze` cloud environment.

The test coverage can be found by running `make test`

##### Notes:

###### The TLDR use case: TODO CLEAN THIS UP

The use case is very simple. If a local machine bricks itself, `sifo` supports the restoration of the the state of the files onto a new machine.

Restoring accidental deletions on a local machine from the cloud is not part of the use case.

Restoring corrupted data on a local machine from the cloud is not part of the use case.

Historical data access is not part of the use case. Snapshotting is out of scope.

The use case is: Replicate the state of local disk to the cloud periodically.

Currently, `sifo` does not support dependencyless configuration encryption. If you want to encrypt your configuration, you can do so by downloading `rclone` and setting it up by following [the docs](https://rclone.org/docs/#configuration-encryption).

Example of using `sifo` using encrypted config:

Encrypting existing config
```bash
> rclone config
Current remotes
s) Set configuration password
```

Set the password.

Using `sifo` to push and pull folders. Need to export `RCLONE_CONFIG_PASS` to be able to read the encrypted `rclone` configuration file.
```bash
export RCLONE_CONFIG_PASS=<password>
sifo push ...
```

##### Why application-layer encryption

Backblaze also offers server-side encryption. Similarly, Rclone also supports a [crypt](https://rclone.org/crypt/) wrapper which can apply its own encryption, 

The rationale for implementing encryption outside of backblaze/rclone is for complete control over the encryption process, guarenteeing end-to-end protection - independent of the rclone & backblaze encryption mechanisms.

If desired, the different encryption layers can work together to provide redundancy.

###### Why full backups

For simplicitiy and reliabilty of restorations, full backups were chosen over differential or incremental backups. See the use case.

##### Versioning the backups

Backblaze offers versioning. For my use case, I turn this off, using the `xxx` config. More in [the docs](). Refer to the use case.

You can set versioning paramaters from within backblaze.


`sifo` avoids system calls to ensure robust security, wide portabilty, and simplified deployment for ease of data restoration.


