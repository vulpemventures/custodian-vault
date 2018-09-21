#!/bin/bash
ARCH=`uname -m`
TARGET="vault_0.10.3_linux_"
case "$ARCH" in
  "x86")
    TARGET=$TARGET"386"
    ;;
  "x86_64")
    TARGET=$TARGET"amd64"
    ;;
  "armv6l" | "armv7l")
    TARGET=$TARGET"arm"
    ;;
  "armv8l")
    TARGET=$TARGET"arm64"
    ;;
  *)
    echo "unsupported OS"
    exit 1
    ;;
esac
TARGET=$TARGET".zip"

echo "Installing unzip.."
sudo apt-get update
sudo apt-get install unzip

echo "Downloading archive.."
wget https://releases.hashicorp.com/vault/0.11.1/$TARGET

rm -rf $HOME/vault
mkdir -p $HOME/vault
echo "Extracting files.."
unzip $TARGET -d $HOME/vault
rm -rf $TARGET
echo "export PATH=\$PATH:$HOME/vault" >> $HOME/.profile
echo -e "Installation complete\nVault installed at path: $HOME/vault"