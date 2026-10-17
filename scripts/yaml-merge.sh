#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'EOF'
usage: yaml-merge.sh <config-file>

Merges the OpenAPI YAML documents listed in <config-file>, in listed order
(later files win on scalar/array conflicts; objects are merged key by key),
and writes the result as YAML to the configured output. If an inputFile is
a directory, every *.yaml file found recursively under it is merged in
sorted path order at that point in the sequence.

Config format (YAML):

  inputs:
    - inputFile: ./service-a.yaml
    - inputFile: ./service-b-fragments/
  output: ./openapi.yaml

inputFile and output are resolved relative to the config file's directory.
EOF
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

config="$1"
if [[ ! -f "$config" ]]; then
  echo "yaml-merge: $config: no such file" >&2
  exit 1
fi

config_dir="$(cd "$(dirname "$config")" && pwd)"

resolve() {
  local p="$1"
  if [[ "$p" = /* ]]; then
    printf '%s' "$p"
  else
    printf '%s' "$config_dir/$p"
  fi
}

config_json="$(yq . "$config")"
out_path="$(resolve "$(jq -r '.output' <<<"$config_json")")"
n="$(jq '.inputs | length' <<<"$config_json")"

if [[ "$n" -eq 0 ]]; then
  echo "yaml-merge: $config: no inputs listed" >&2
  exit 1
fi

files=()
for ((i = 0; i < n; i++)); do
  f="$(resolve "$(jq -r ".inputs[$i].inputFile" <<<"$config_json")")"
  if [[ -d "$f" ]]; then
    while IFS= read -r sub; do
      files+=("$sub")
    done < <(find "$f" -name '*.yaml' | sort)
  elif [[ -f "$f" ]]; then
    files+=("$f")
  else
    echo "yaml-merge: $f: no such file or directory" >&2
    exit 1
  fi
done

if [[ ${#files[@]} -eq 0 ]]; then
  echo "yaml-merge: $config: no yaml files found among the configured inputs" >&2
  exit 1
fi

mkdir -p "$(dirname "$out_path")"
yq -y -s 'reduce .[] as $doc ({}; . * $doc)' "${files[@]}" >"$out_path"

echo "merged ${#files[@]} file(s) into $out_path" >&2
