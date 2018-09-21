#!/bin/bash
ARCH=`uname -m`
TARGET="vault_0.11.1_darwin_"
case "$ARCH" in
  "x86")
    TARGET="386"
    ;;
  "x86_64")
    TARGET="amd64"
    ;;
  *)
    echo "unsupported OS"
    exit 1
    ;;
esac
TARGET=$TARGET".zip"

echo "Downloading archive.."
curl -O https://releases.hashicorp.com/vault/0.11.1/$TARGET

rm -rf $HOME/vault
mkdir -p $HOME/vault
echo "Extracting files.."
unzip $TARGET -d $HOME/vault
rm -rf $TARGET
echo "export PATH=\$PATH:$HOME/vault" >> $HOME/.bash_profile
echo "export VAULT_ADDR=http://localhost:8200" >> $HOME/.bash_profile

echo -e "Installation complete\nVault installed at path: $HOME/vault"