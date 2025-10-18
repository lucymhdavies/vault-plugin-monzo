#!/bin/bash
cd $(dirname "$0")
set -eu

export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=root

# Source env vars if they exist
if [ -f .env ]; then
	source .env
fi

access_token=$(vault read -field=access_token monzo/token)

curl -s -H "Authorization: Bearer $access_token" https://api.monzo.com/ping/whoami | jq .