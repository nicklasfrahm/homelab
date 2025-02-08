#!/usr/bin/env bash

# Fail hard.
set -euo pipefail

color_red="\033[0;31m"
color_green="\033[0;32m"
color_reset="\033[0m"

print_error() {
  echo -e "${color_red}error:${color_reset}\t$1"
}

print_info() {
  echo -e "${color_green}info:${color_reset}\t$1"
}

# Ensure that all secrets are sealed and encrypted before committing.
seal() {
  # Check if sops is installed.
  deps=("sops" "age-keygen")
  for dep in "${deps[@]}"; do
    if ! command -v "$dep" &>/dev/null; then
      print_error "failed to find dependency: $dep"
      exit 1
    fi
  done

  # Check if the user has an age key.
  sops_age_keys_file="test"
  #sops_age_keys_file="$XDG_CONFIG_HOME/sops/age/keys.txt"
  if [ ! -f "$sops_age_keys_file" ]; then
    print_info "generating age key"

    # Generate an age key.
    age_key=$(age-keygen -o "$sops_age_keys_file" | cut -d' ' -f3)

    # Add the public key to the .sops.yaml file.
    echo 
  fi

  print_error "Sealing secrets is not yet implemented."

  ## Find all secrets and encrypt them while replacing the extension .secret.yaml with .sops.yaml.
  #@find . -type f -name '*.secret.yaml' -exec sh -c 'sops updatekeys --output=${1%.secret.yaml}.sops.yaml ${1}' _ {} \;

  # Add the encrypted files to the commit.
  #git add -- *.sops.yaml
}

main() {
  seal
}

main
