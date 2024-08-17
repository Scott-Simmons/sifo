# Push and pull encrypted data to Backblaze

![Build Status](https://img.shields.io/github/actions/workflow/status/Scott-Simmons/backup-system/ci.yml?branch=main)

Implements a "push" and "pull" functionality using Rclone + AES-256 encryption.

The purpose is for periodic backups: I want to be able to be confident that my machine dying will not compromise my data.

However, I am not in love with the idea of uploading all of my data to the cloud unencrypted.

The way that this works on my machine is via "push" functionality.

1. Archive the folder
2. Encrypt the folder
3. Move the folder to the backblaze

The encrypted data on backblaze can be restored to local disk via "pull" functionality.

1. Sync a single backblaze file to a local directory. (TODO: Validate this at execution time)
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


### Why full backups?

Can split backup options by:

- Incremental backups
- Differential backups
- Full backups (what this package implements. Simple, Reliable, Inefficient).

There are trade offs here. Backup frequency, storage space, backup time, recovery time, reliablity, complexity.

Can split backup options by:

- Overwrite
- Maintain last N versions (versioning).


### What this application does not do:

A file was mistakenly deleted a week ago, having a full backup from a week ago enables you to recover that file. 

If a file on local disk becomes corrupted and I want to restore it.

Historical data access.

### What this application does do:

If my box bricks itself, I can restore the most recent state of my file tree on a new box.

I do not care about recent data corruption or accidental deletions. Just replicate whate what is on my local disk to the cloud.


### Tough problem: How to do integration and end-to-end tests

- Need to test recovery
- Need to test backup


## Another note: Backblaze offers encryption. Rclone offers a crypt wrapper that does encryption. But I want to implement it myself for learning.

## Another note: You can setup server side encryption in backblaze. But this is managed by backblaze.

## Another note: No snapshotting. Not part of use case.

## Triple layer protection: Backblaze, Rclone Crypt Wrapper, Manual Encruption. 3 keys. For now just implement Manual.
