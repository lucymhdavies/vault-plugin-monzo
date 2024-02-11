#!/bin/bash
set -ex


# Accessing the secret
# (This is where I'd pull and use my HCP-VS creds)
echo "My secret value is: $MY_SECRET"
# Might be this...
# export "HCP_CLIENT_ID=foo" >> ~/.bashrc
# export "HCP_CLIENT_SECRET=foo" >> ~/.bashrc

# TODO: Also create a .vlt.json from the secret env vars


# Install Vault and VLT
# https://developer.hashicorp.com/vault/tutorials/getting-started/getting-started-install
sudo apt update && sudo apt install gpg wget
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
gpg --no-default-keyring --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg --fingerprint
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install vault vlt

# Vault connection details
echo "export VAULT_ADDR=http://127.0.0.1:8200" >> ~/.bashrc
echo "export VAULT_TOKEN=root" >> ~/.bashrc
# TODO: set up .bashrc