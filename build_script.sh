#!/bin/bash

set -e

DATETAG=$(date +"%y%m%d%H%M%S")
IMAGE_NAME="registry.iguiyu.com/tcpforward:$DATETAG"

echo $IMAGE_NAME

cd $GOPATH/src/github.com/kobeld/tcp-forwarder

echo "Building application..."
CGO_ENABLED=0 GOOS=linux go build -o main .

echo "Building docker image..."
docker build -t $IMAGE_NAME .

echo "Cleanup resources..."
rm main

echo "Done"