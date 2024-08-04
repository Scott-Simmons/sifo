# Push and pull encrypted data to Google Drive

![Build Status](https://img.shields.io/github/actions/workflow/status/Scott-Simmons/backup-system/ci.yml?branch=main)

Implements a "push" and "pull" functionality using Rclone + AES-256 encryption.

The purpose is for periodic backups: I want to be able to be confident that my machine dying will not compromise my data.

However, I am not in love with the idea of uploading all of my data to the cloud unencrypted.

The way that this works on my machine is via "push" functionality.

1. Archive the folder
2. Encrypt the folder
3. Move the folder to the google drive

The encrypted data on google drive can be restored to local disk via "pull" functionality.

1. Sync a single google drive file to a local directory. (TODO: Validate this at execution time)
2. Unencrypt.
3. Decompress.

This tool is intended to be fully portable, with no system calls used.

## TODO:

### Important TODOs:

TODO: Fix decryption bug
TODO: Fully wrap Rclone config create
TODO: Pruning revisions logic

### Less important TODOs:
TODO: Consider cross compilation
TODO: Consider a deployment/releases workflow
TODOS: Consider refactoring rclone sync to be strict on validating that syncing to local wont delete files. It is a dangerous operation. Maybe switch to copy.
TODO: Make sure test resources are cleaned up.
