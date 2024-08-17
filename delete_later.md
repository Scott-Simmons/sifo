
I need to move bucket setup to this package. Rclone needs to be fully encapsulated by this package.

Steps:

1. rclone config (create new backblaze remote).
2. Specify name: <foo> e.g. backblaze_backup
3. Specify type (backblaze... this should be fixed and standard...)
4. Now specify account ID and application key ID.
5. Other configuration (defaults)
- Do not delete files on remote removal.

Note you will need to set up a bucket manually with backblaze.

This is (1) set up a bucket. and (2) set up application key.


I also need to see what is in that bucket:

1. Wrap around rclone ls <bucket_name>

I also need to see what buckets I have

1. Wrap around rclone show buckets

My app key:
My app key id:
