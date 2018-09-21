#!/bin/bash
set -e

brew uninstall go
brew update
brew install golang

mkdir -p $HOME/go/bin $HOME/go/src $HOME/go/pkg

echo "export GOPATH=$HOME/go"
echo "export GOROOT=/usr/local/opt/go/libexec"
echo "export PATH=$PATH:$GOPATH/bin"
echo "export PATH=$PATH:$GOROOT/bin"

echo -e "Installation complete\nGo installed at path: $GOROOT\nWorkspace set at GOPATH=$GOPATH"