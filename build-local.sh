#!/bin/bash

if [ ! -f "/.dockerenv" ]; then
    docker pull golang:latest && docker run -ti --rm -v "${PWD}":"${PWD}" --workdir "${PWD}" golang:latest
fi;

apt-get update &&
apt-get install -yqq \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libglx-dev \
    libgl-dev \
    libxxf86vm-dev &&

git config --global --add safe.directory "${PWD}" &&

go build -v -ldflags "-s -w -X 'main.version=0.0.0' -X 'main.commit=dev' -X 'main.date=$(date '+%Y-%m-%d %H:%M:%S')' -X 'main.builtBy=manual'" .;

exit ${?}

# GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=amd64-linux-gnu-gcc CXX=amd64-linux-gnu-g++ HOST=amd64-linux-gnu go build -v
# x86_64-pc-linux-gcc
