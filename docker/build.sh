#!/bin/sh
set -e
set -x

if go version | grep '^go version go1.5.[0-9]' > /dev/null \
   && command -v glide > /dev/null; then
    export GO15VENDOREXPERIMENT=1
    glide install
else
    go get -v
fi

REPOSITORY=github.com/0rax/go-redirect
docker run -it --rm -v $GOPATH/src:/tmp/go/src \
                    -e GOPATH=/tmp/go \
                    -w /tmp/go/src/$REPOSITORY \
                    golang:1.5-alpine \
                    sh -c 'go build -v'
