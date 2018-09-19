# custodian-vault

This repository contains a plugin for Vault that let's you store your crypto in a secure way. It exploits Vault's capabilities of storing encrypted data to securely create new random wallets, get addresses and sign raw transactions and setting strict policy rules to access and manage your funds.  
At the moment only bitcoin is supported.

## Table of Contents

* [Prerequisites](#prerequisites)
* [Installation](#installation)
* [Usage](#usage)
* [Tests](#tests)
* [Troubleshooting](#troubleshooting)

## Prerequisites

* [Golang](https://golang.org/)
* [Vault](https://www.vaultproject.io/)

If you have already installed them on your machine, skip this step.

Clone the project

```sh
git clone https://github.com/vulpemventures/custodian-vault.git && cd custodian-vault
```

Run `./scripts/go_installer.sh` to install Go. It will be installed at `/usr/local/go` and will export environment variable `GOPATH=$HOME/go`.  
Run `./scripts/vault_installer.sh` to install Vault. It will be installed at `$HOME/vault`.

NOTICE: These scripts are intended to install packages only into Linux systems.

Delete these folders to uninstall the packages.  

## Installation

Create the path for the project in your go workspace `GOPATH`

```sh
mkdir -p  $GOPATH/src/github.com/vulpemventures
```

### If you followed the previous step

Move the folder of the project into the path

```sh
mv ../custodian-vault $GOPATH/src/github.com/vulpemventures
```

### If you skipped the previous step

Clone the project in the directory

```sh
git clone https://github.com/vulpemventures/custodian-vault.git $GOPATH/src/github.com/vulpemventures/custodian-vault
```

You can now start an instance of Vault and automatically install the plugin launching the script

```sh
./scripts/start_dev.sh
```

This starts Vault in dev mode, in order to let you test the features of the custodian plugin.

## Usage

### Create a BIP44 wallet

```sh
vault write custodian/wallet/<name> network=<testnet|mainnet>
```

This creates a random mnemonic that's returned in the response object. It is used to generate the seed as stated in [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki).  
NOTICE: In order to enforce to backup the mnemonic, it is showed to the end user only once, at creation time, and won't be accessible anymore in the future.

To get info about a wallet:

```sh
vault read custodian/wallet/<name>
```

It returns a response object that contains:

* `network` either bitcoin main net or test net
* extended public key at path `m/44'/0'/0'/0` for main net, that is the master key from which all receiving address are generated (see [BIP-0032](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) for more details).

### Create a BIP49 wallet

```sh
vault write custodian/wallet/<name> network=<testnet|mainnet> segwit=true
```

Again, this creates a mnemonic visible only at creation time, but derivation path is set to `m/49'/0'/0'/0` (for mainnet). This kind of wallet is used for segwit backward compatibility.

### Create multisig wallet

```sh
vault write custodian/wallet/multisig/<name> m=<number> n=<number> pubkeys=<list,of,pubkeys> network=<testnet|mainnet>
```

This creates a new key used in an `m` of `n` multisg wallet. It requires a list of `n-1` comma separated public keys that are used to create the `redeem script`. This will be subject of breaking changes, since it's planned to add support for [BIP-0045](https://github.com/bitcoin/bips/blob/master/bip-0045.mediawiki).

To get the `redeem script` along with the other info:

```sh
vault read custodian/wallet/multisig/<name>
```

### Create a BIP84 wallet (native segwit)

```sh
vault write custodian/wallet/segwit/<name> network=<mainnet|testnet>
```
This creates a native segwit wallet.

### Generate an `auth_token` for a wallet

```sh
vault read custodian/creds/<wallet_name>
# Expected output
# lease_id           <lease_id>
# lease_duration     5m
# lease_renewable    false
# token              <auth_token>
```

This creates an `auth_token` to allow the end user to interact with the created wallet. Operations like generating a new receiving address or creating a signature for a raw transaction require a token to be passed in the request.

NOTICE: Auth tokens expire right after a request is satisfied and have a default Time To Live of 5 minutes.

To manually revoke an `auth_token`:

```sh
vault lease revoke <lease_id>
```

### Generate an `auth_token` for a multisig wallet

```sh
vault read custodian/creds/multisig/<wallet_name>
# Same response object as above
```

### Generate an `auth_token` for a native segwit wallet

```sh
vault reas custodian/creds/segwit/<wallet_name>
# Same response object as above
```

### Derive a new receiving address for a wallet

```sh
vault write custodian/address/<wallet_name> token=<auth_token>
```

This derives new P2PKH addresses for BIP44 wallet and P2WPKH nested in P2SH for BIP49 wallet.

### Get receiving address for a multisig wallet

```sh
vault write custodian/address/multisig/<wallet_name> token=<auth_token>
```

This returns the receiving address of a previous created multisig.  
The address won't change since it's the `base58check` encode of the hash of the `redeem script`.  
Also this feature will change when BIP-0045 will be supported.

### Derive a new receiving address for native segwit wallet

```sh
vault write custodian/address/segwit/<wallet_name> token=<auth_token>
```

This derives new P2WPKH addresses in `Bech32` format.

### Sign raw transactions

```sh
vault write custodian/transaction/<wallet_name> mode=<standard|multisig|segwit> rawTx=<string> token=<auth_token>
```

This will create a signature for the passed raw transaction.  
You need to specify if the wallet is a multisig, native segwit or standard (either BIP44 and BIP49) type, this flag is set to `standard` by default.
The produced signature is deterministic, which means that the same message and the same key yield the same signature, and canonical in accordance with [RFC6979](https://tools.ietf.org/html/rfc6979) and [BIP-0062](https://github.com/bitcoin/bips/blob/master/bip-0062.mediawiki) respectively.

## Tests

// TODO

## Troubleshooting

If you get a `server gave HTTP response to HTTPS client` error:

```sh
export VAULT_ADDR='http://127.0.0.1:8200'
```
