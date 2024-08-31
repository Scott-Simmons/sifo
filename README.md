# Sifo: Securely Push and Pull to and from Backblaze

![Build Status](https://img.shields.io/github/actions/workflow/status/Scott-Simmons/backup-system/ci.yml?branch=main)

`sifo` enables cloud backup and restore functionality using Backblaze as the cloud storage provider, and Rclone as the syncronisation tool.

`sifo` also implements archive (tar.gz) and AES-256 (CBC) encryption at the application level on the local server.

##### Dependencies

None. Unless you are [building from source](link to more).

`librclone` exports [shims]() that wrap over the Rclone RPC. All rclone dependencies are included in the final binary at compile time for `sifo`.

##### System Requirements

`sifo` currently targets:

- the mainstream linux platforms: foo bar baz
- 

##### Installation

Download the appropriate binary for your system, or build from source.

##### Building from source

Requirements

- golang version xxx.
- make version xxx.

```bash
make install

```

##### Configuration

Setting the rclone configuration file.

```bash

```

Generating the encryption key
```bash

```

### Usage:

Initialising the backblaze remote:
```bash

```

Initialising the AES-256 (symmetric) encryption key:
```bash

```

Pushing a folder from local to backblaze remote:
```bash

```

Pulling a folder from backblaze remote to local:
```bash

```

##### Testing

Much of the functionality is unit tested, and some integration tests for the RPC calls. More work can be done on the test suite.

Currently, end-to-end testing is not implemented. This is due to the challenges in setting up a suitable test environment for the end-to-end tests to interact with a real backblaze cloud environment.

The unit test coverage is here:

INSERT COVERAGE

##### Notes:

###### The TLDR use case:

The use case is very simple. If a local machine bricks itself, `sifo` supports the restoration of the the state of the files onto a new machine.

Restoring accidental deletions on a local machine from the cloud is not part of the use case.

Restoring corrupted data on a local machine from the cloud is not part of the use case.

Historical data access is not part of the use case. Snapshotting is out of scope.

The use case is: Replicate the state of local disk to the cloud periodically.

Currently, `sifo` does not support dependencyless configuration encryption. If you want to encrypt your configuration, you can do so by downloading `rclone` and setting it up by following [the docs](https://rclone.org/docs/#configuration-encryption).

Example of using `sifo` using encrypted config:

Encrypting existing config
```bash

```

Using `sifo` to push and pull folders.
```bash

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


