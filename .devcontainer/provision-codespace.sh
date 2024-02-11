#!/bin/bash

echo "Hello from provision-codespace.sh"
echo "This script is executed when the Codespace starts."

# Accessing the secret
echo "My secret value is: $MY_SECRET"


touch codespace.txt
