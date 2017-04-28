#!/bin/bash

#run inside the docker dir!
. app_env.sh

image_name="$HUB_REPO/$APP_NAME:$APP_BRANCH"
echo "Building $image_name"

# Make build folder
rm -rf build
mkdir -p build

# Make the final app folder
rm -rf release/app
mkdir -p release/app

# Build go app
mkdir -p build/go
cp -r ../go/src build/go
docker run --rm -v "$(pwd)/build/go":/go -w /go -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 golang:1.7 go build -o $APP_NAME $APP_NAME
mv build/go/$APP_NAME release/app/$APP_NAME

cd release
docker build -t $image_name .

echo "Finished building $image_name"
