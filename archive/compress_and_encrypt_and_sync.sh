#!/bin/bash
set -euo pipefail

GOOGLE_DRIVE_REMOTE=google-drive-backup
DEFAULT_GPG_ID=E422078FC8B46B08282BC8514A1C546EBA5C02D7
PATH_TO_BACKUP=$1
DIR_TO_BACKUP_TO=$2
GPG_ID=${GPG_ID:-${DEFAULT_GPG_ID}}

if [ -z "$1" ]; then
    echo "Error: Missing argument. Please provide the source directory to compress" >&2
    exit 1
fi
if [ -z "$2" ]; then
    echo "Error: Missing argument. Please provide the target directory to compress to." >&2
    exit 1
fi
if [[ ! -e "${PATH_TO_BACKUP}" ]]; then
  echo "Error: Folder to backup:" "${PATH_TO_BACKUP}" "was not found on machine" >&2
  exit 1
fi

mkdir -p "$DIR_TO_BACKUP_TO"
echo "Backing up the following: ${PATH_TO_BACKUP}"

folder_to_archive=$(basename "${PATH_TO_BACKUP}")
context_dir=$(dirname "${PATH_TO_BACKUP}")
backup_path="$DIR_TO_BACKUP_TO/${folder_to_archive}.tar.gz"
touch "${backup_path}"

### Compression
echo Compressing "${folder_to_archive}" to "${backup_path}" using context "${context_dir}"
echo tar -czf "${backup_path}" -C "${context_dir}" "${folder_to_archive}"
tar -czf "${backup_path}" -C "${context_dir}" "${folder_to_archive}" &
pid=$!
wait $pid
echo Compression complete


### Encryption
echo Encrypting backup data "${backup_path}" to "${backup_path}.gpg"
gpg \
  --yes \
  --output "${backup_path}".gpg \
  --encrypt \
  --recipient "${GPG_ID}" \
  --cipher-algo AES256 \
  "${backup_path}" &

pid=$!
wait $pid

if [ $? -eq 0 ]; then
  echo Encryption successful
  rm -rf "${backup_path}"
  echo Compressed folder deleted
else
  echo "Encryption failed" >&2
  exit 1
fi

echo Syncing encrypted files with google drive
echo rclone sync "${backup_path}".gpg "${GOOGLE_DRIVE_REMOTE}":/ &
rclone sync "${backup_path}".gpg "${GOOGLE_DRIVE_REMOTE}":/ &

pid=$!
wait $pid
