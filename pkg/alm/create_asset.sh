#!/usr/bin/env bash

usage() { echo "usage: $0 -n <string> -d <string> -l <string> -u <string>" 1>&2; exit 1; }

while getopts ":n:d:l:u:" o; do
  case "${o}" in
    n)
      asset_name=${OPTARG}
      ;;
    d)
      asset_description=${OPTARG}
      ;;
    l)
      label_dir=${OPTARG}
      ;;
    u)
      csv_file=${OPTARG}
      ;;
    *)
      usage
      ;;
  esac
done
shift $((OPTIND-1))

if [ -z "${asset_name}" ] || [ -z "${asset_description}" ] || [ -z "${label_dir}" ] || [ -z "${csv_file}" ]; then
    usage
fi

tmp_label_path=$(mktemp)
asset_id=$($ALM_EXEC render -o "${tmp_label_path}" -n "${asset_name}" -d "${asset_description}")
echo "Generated asset ID: ${asset_id}"

label_path="${label_dir}/${asset_id}.svg"
mv "${tmp_label_path}" "${label_path}"

cp "${csv_file}" "${csv_file}.bak"
echo "Backed up CSV file to ${csv_file}.bak"

echo "Writing asset to CSV records."
$ALM_EXEC save_asset --output_path "${csv_file}" \
  --asset_id "${asset_id}" \
  --asset_name "${asset_name}" \
  --asset_description "${asset_description}"

echo "Done. Label written to ${label_path}"
