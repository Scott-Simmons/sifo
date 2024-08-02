#!/bin/bash

echo "Pruning revisions"

# To save on storage space
limit=1

GOOGLE_FILE_ID_PLACEHOLDER=ID_PLACEHOLDER
GOOGLE_REVISIONS_ENDPOINT=https://www.googleapis.com/drive/v3/files/\{"${GOOGLE_FILE_ID_PLACEHOLDER}"\}/revisions
echo "Using generic endpoint: ${GOOGLE_REVISIONS_ENDPOINT}"

echo "Gathering google drive file IDs..."
readarray -t ids <<< "$(rclone lsjson google-drive-backup: | jq -r '.[] | .ID')"
echo "${ids[@]}"


# Their could be 404 errors

for file_id in "${ids[@]}"; do
  echo "$file_id"
  auth_info=$(rclone config dump | jq '. | ."google-drive-backup"')
  access_token=$(echo "$auth_info" | jq -r '.token' | jq -r '.access_token')
  specific_endpoint="${GOOGLE_REVISIONS_ENDPOINT//"${GOOGLE_FILE_ID_PLACEHOLDER}"/"${file_id}"}"
  echo "Targeting ${specific_endpoint}"
  data=$( \
    curl \
      -H 'GData-Version: 3.0' \
      -H "Authorization: Bearer ${access_token}" \
      "${specific_endpoint}"
    )
  n_revisions=$(echo "$data" | jq '.revisions | length')
  number_to_chop=$(($n_revisions - $limit))

  rev_ids_to_delete=$( \
    echo "$data" \
    | jq --argjson N "${number_to_chop}" '
    {
      kind: .kind,
      revisions: (
        .revisions
        | sort_by(.modifiedTime)
        | .[:$N]
      )
    }
  ' | jq '.revisions[].id'
  )
  readarray -t rev_ids_to_delete <<< "${rev_ids_to_delete}"

  # Now delete the revisions
  for revision_id in "${rev_ids_to_delete[@]}"; do
    revision_id=$(echo "${revision_id}" | tr -d '"')
    echo Deleting revision "${revision_id}" for file "${file_id}"
    specific_delete_endpoint="${specific_endpoint}"/"${revision_id}"
    curl \
      -X DELETE \
      -H 'GData-Version: 3.0' \
      -H "Authorization: Bearer ${access_token}" \
      "${specific_delete_endpoint}"
  done
done



