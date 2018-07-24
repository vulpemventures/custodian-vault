#!/bin/bash

export VAULT_ADDR='http://127.0.0.1:8200'

vault login root

SHASUM=$(shasum -a 256 "$HOME/tmp/vault-plugins/custodian-vault" | cut -d " " -f1)

vault write sys/plugins/catalog/custodian-vault   sha_256="$SHASUM"   command="custodian-vault"

vault secrets enable -path=custodian -plugin-name=custodian-vault plugin