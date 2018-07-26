# custodian-vault

Store Bitcoin and Ethereum hot wallet securely using Vault.

## Prerequisites

* [Golang](https://golang.org/)
* [Vault](https://www.vaultproject.io/)

If you have already installed them, skip this step.

Clone the project

```sh
git clone https://github.com/vulpemventures/custodian-vault.git && cd custodian-vault
```

Run `./scripts/go_installer.sh` to install Go. It will be installed at `/usr/local/go`.  
Run `./scripts/vault_installer.sh` to install Vault. it will be installed at `$HOME/vault`.

Delete these folders to uninstall the packages.  

## Install plugin

Create the path for the project in your go workspace `GOPATH`

```sh
mkdir -p  $GOPATH/src/github.com/vulpemventures
```

Clone the project in the directory

```sh
git clone https://github.com/vulpemventures/custodian-vault.git $GOPATH/src/github.com/vulpemventures/custodian-vault
```

Or move project folder if you already cloned the repo

```sh
mv ../custodian-vault $GOPATH/src/github.com/vulpemventures
```

Create the directory where to save the binary of the project

```sh
mkdir -p ~/tmp/vault-plugins

cd $GOPATH/src/github.com/vulpemventures/custodian-vault
go get .
go build -o ~/tmp/vault-plugins/custodian-vault
```

Create the config file to point the Vault server to the plugin directory

```sh
tee ~/tmp/vault.hcl <<EOF
plugin_directory = "$HOME/tmp/vault-plugins"
EOF
```

Start Vault server in `dev` mode passing the config file

```sh
vault server -dev -dev-root-token-id="root" -config=$HOME/tmp/vault.hcl
```

Open another tab in the terminal and install the custodian plugin

```sh
./scripts/plugin_installer.sh
```

## Usage

Create a wallet

```sh
vault write custodian/wallet/test network=testnet
```

Get wallet extended public key (this is planned to be executed passing `auth_token`)

```sh
vault read custodian/wallet/test
```

Generate an `auth_token` for wallet

```sh
vault read custodian/creds/test
# Expected output
# lease_id           <lease_id>
# lease_duration     1h
# lease_renewable    true
# token              <auth_token>
```

Get a receiving address passing the new generated token

```sh
vault write custodian/address/test token=<auth_token>
```

Renew or revoke `auth_token`

```sh
vault lease renew|revoke <lease_id>
```

## Troubleshooting

If you get a `server gave HTTP response to HTTPS client` error

```sh
export VAULT_ADDR='http://127.0.0.1:8200'
```