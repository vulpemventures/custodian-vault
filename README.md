# custodian-vault

Store securely Bitcoin and Ethereum hot wallet using Vault

## Install plugin

Move to your go workspace directory and create a path for the project

```sh
cd $GOPATH/src
mkdir -p  github.com/vulpemventures
cd github.com/vulpemventures/
```

Clone the project

```sh
cd $GOPATH/src/github.com/vulpemventures/
git clone git@github.com:vulpemventures/custodian-vault.git
```

Create a directory where to save the binary of the project

```sh
mkdir -p ~/tmp/vault-plugins
go build -o ~/tmp/vault-plugins/custodian-vault
```

Create a config file to point Vault at the plugin directory

```sh
tee ~/tmp/vault.hcl <<EOF
plugin_directory = "$HOME/tmp/vault-plugins"
EOF
```

Start the vault server in dev mode passing the config file

```sh
vault server -dev -dev-root-token-id="root" -config=$HOME/tmp/vault.hcl
```

Open another tab, always in the repo directory, make the plugin installer script executable and launch it

```sh
chmod u+x install_plugin.sh
./install_plugin.sh
```

## Usage

Create a wallet

```sh
vault write custodian/wallet/test network=testnet
```

Read wallet info

```sh
vault read custodian/wallet/test
```

Generate an auth token for wallet

```sh
vault read custodian/creds/test
# Expected output
# lease_id           custodian/creds/test/<salted_auth_token>
# lease_duration     1h
# lease_renewable    true
# token              <auth_token>
```

Generate a new address passing the new generated token

```sh
vault read custodian/address/test token=<auth_token>
```

To renew or revoke auth token

```sh
vault lease renew|revoke <lease_id>
```

To disable plugin run

```sh
vault secrets disable custodian
```

## Troubleshooting

If get "server gave HTTP response to HTTPS client" error

```sh
export VAULT_ADDR='http://127.0.0.1:8200'
```