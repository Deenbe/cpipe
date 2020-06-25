#!/bin/bash

set -ex

BINARY=""
if [[ "$(uname)" = "Linux" ]]; then
  BINARY="cpipe_linux_amd64"
fi

if [[ "$(uname)" = "Darwin" ]]; then
  BINARY="cpipe_darwin_amd64"
fi

if [[ "$BINARY" = "" ]]; then
    echo "OS $(uname) is not supported"
    exit 1
fi

curl -o /usr/local/bin/cpipe -L "https://github.com/Deenbe/cpipe/releases/latest/download/$BINARY" 

