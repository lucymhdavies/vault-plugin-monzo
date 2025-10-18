#!/bin/bash
cd $(dirname "$0")
set -eu

export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=root


# Source env vars if they exist
if [ -f .env ]; then
	source .env
fi

# Disable the plugin if it already exists
vault secrets disable monzo

# Then mount it
vault secrets enable monzo

# And configure, using the creds we either got from env or vlt
if [ -z ${REDIRECT_BASE_URL:-} ] ; then
	vault write monzo/config client_id=${MONZO_CLIENT_ID} client_secret=${MONZO_CLIENT_SECRET}
else
	vault write monzo/config client_id=${MONZO_CLIENT_ID} client_secret=${MONZO_CLIENT_SECRET} redirect_base_url="${REDIRECT_BASE_URL}"
fi

vault read monzo/config
