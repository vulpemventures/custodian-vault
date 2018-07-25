#!/bin/bash
echo "Removing previous installed version.."
sudo rm -rf /usr/local/go
rm -rf $HOME/go 

ARCH=`uname -m`
TARGET="go1.10.3.linux-"
case "$ARCH" in
  "x86")
    TARGET=$TARGET"386"
    ;;
  "x86_64")
    TARGET=$TARGET"amd64"
    ;;
  "armv6l" | "armv7l")
    TARGET=$TARGET"armv6l"
    ;;
  "armv8l")
    TARGET=$TARGET"arm64"
    ;;
  *)
    echo "unsupported OS"
    exit 1
    ;;
esac
TARGET=$TARGET".tar.gz"

echo "Downloading archive.."
wget https://dl.google.com/go/$TARGET

echo "Extracting files.."
sudo tar -C /usr/local -xzf $TARGET

rm -rf $TARGET

echo "export PATH=\$PATH:/usr/local/go/bin" >> $HOME/.profile

mkdir -p $HOME/go/bin $HOME/go/src $HOME/go/pkg
echo "export GOPATH=\$HOME/go" >> $HOME/.profile

echo -e "Installation complete\nGo installed at path: /usr/local/go\nWorkspace set at GOPATH=$GOPATH"