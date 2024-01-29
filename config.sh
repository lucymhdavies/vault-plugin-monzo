#!/bin/bash
set -eu

cd $(dirname "$0")


# Source env vars if they exist
if [ -f .env ]; then
	source .env
fi
vault write monzo/config client_id=${MONZO_CLIENT_ID} client_secret=${MONZO_CLIENT_SECRET}
