#!/bin/bash
OS=`uname -s`
TARGET="vault_0.11.1_"
case "$OS" in
  "Darwin")
    BASHRC=".bash_profile"
    TARGET=$TARGET"darwin_"
    ;;
  "Linux")
    BASHRC=".profile"
    TARGET=$TARGET"linux_"
    ;;
  *)
    echo "Unsupported OS"
    exit 1
    ;;
esac

ARCH=`uname -m`
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

if [$OS -eq "Darwin"]
then
  echo "Installing unzip.."
  sudo apt-get update
  sudo apt-get install unzip
fi

echo "Downloading archive.."
curl -O https://releases.hashicorp.com/vault/0.11.1/$TARGET

rm -rf $HOME/vault
mkdir -p $HOME/vault
echo "Extracting files.."
unzip $TARGET -d $HOME/vault
rm -rf $TARGET
echo "export PATH=\$PATH:$HOME/vault" >> $HOME/$BASHRC
echo "export VAULT_ADDR=http://localhost:8200" >> $HOME/$BASHRC

echo -e "Installation complete\nVault installed at path: $HOME/vault"