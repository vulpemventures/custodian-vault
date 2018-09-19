
#!/usr/bin/env bash
set -e

MNT_PATH="custodian"
PLUGIN_NAME="custodian-vault"
#
# Helper script for local development. Automatically builds and registers the
# plugin. Requires `vault` is installed and available on $PATH.
#

echo "==> Starting dev"

echo "--> Scratch dir"
echo "    Creating"
SCRATCH="$HOME/tmp"
mkdir -p "$SCRATCH/plugins"

echo "--> Vault server"
echo "    Writing config"
tee "$SCRATCH/vault.hcl" > /dev/null <<EOF
plugin_directory = "$SCRATCH/plugins"
EOF

echo "    Envvars"
export VAULT_DEV_ROOT_TOKEN_ID="root"
export VAULT_ADDR="http://127.0.0.1:8200"

echo "    Starting"
vault server \
  -dev \
  -log-level="debug" \
  -config="$SCRATCH/vault.hcl" \
  &
sleep 2
VAULT_PID=$!

function cleanup {
  echo ""
  echo "==> Cleaning up"
  kill -INT "$VAULT_PID"
  rm -rf "$SCRATCH"
}
trap cleanup EXIT

echo "    Authing"
vault auth root &>/dev/null

echo "--> Building"
go build -o "$SCRATCH/plugins/$PLUGIN_NAME"
SHASUM=$(shasum -a 256 "$SCRATCH/plugins/$PLUGIN_NAME" | cut -d " " -f1)

echo "    Registering plugin"
vault write sys/plugins/catalog/$PLUGIN_NAME \
  sha_256="$SHASUM" \
  command="$PLUGIN_NAME"

echo "    Mouting plugin"
vault secrets enable -path=$MNT_PATH -plugin-name=$PLUGIN_NAME plugin

echo "    Mounting 2fa plugin"
vault secrets enable totp

echo "==> Ready!"
wait $!
