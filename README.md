### Google drive backups

There is a need for my local drive to be backed up on the cloud:

1. For accessibilty: I want to be able to access my work from multiple machines via google drive.
2. For backups: I want to be able to be confident that my machine dying will not compromise my data

However, I am not in love with the idea of uploading all of my data to the cloud unencrypted.

1. Because servers can get hacked.
2. Because I cannot trust that cloud providers aren't using my data to train LLMs or segment/cluster the population based on our files.

For this reason, this google drive folder stores encrypted folders.

The way that this works on my machine is with a backup and sync script:

1. Archive the folder
2. Encrypt the folder
3. Move the folder to the google drive
4. Run this every N minutes.

If I want to access my files from another machine, I can login via my cloud provider and download the archive.

Alternatively, I can clone this repo and restore using the utilities in this script.

This has the added benefit of being compressed so that I can stay on the free tier of google drive.

