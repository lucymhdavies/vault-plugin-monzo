#!/bin/bash
set -eu

# Source env vars if they exist
source $(dirname "$0")/.env

vault write monzo/config client_id=${MONZO_CLIENT_ID} client_secret=${MONZO_CLIENT_SECRET}