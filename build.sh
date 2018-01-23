#!/bin/bash

GO_PATH="$HOME/go"
SRC_PATH="${GO_PATH}/src/github.com/hjf288/sgemu"
PACKAGES="Data/Extractor LoginServer GameServer"
BUILD_PATH="$HOME/sgemu/binaries"

export GOPATH="${GO_PATH}"

if [[ ! -d $GO_PATH ]]; then
    echo "GO_PATH not found, creating at $GO_PATH"
    mkdir $GO_PATH
fi

if [[ ! -d $BUILD_PATH ]]; then
    echo "BUILD_PATH not found, creating at $BUILD_PATH"
    mkdir $BUILD_PATH
fi

# Get source

for PACKAGE in $PACKAGES; do
    echo "Grabbing source: $PACKAGE"
    go get -u github.com/hjf288/sgemu/"${PACKAGE}" || exit "Failed"
done

# Do Build

for PACKAGE in $PACKAGES; do
    if [[ $PACKAGE = "Data/Extractor" ]]; then
        OUTPUT_NAME="Extractor"
    else
        OUTPUT_NAME="${PACKAGE}"
    fi

    if [[ ! -d "$BUILD_PATH/${OUTPUT_NAME}" ]]; then
        echo "${BUILD_PATH}/${OUTPUT_NAME} not found, creating"
        mkdir "${BUILD_PATH}/${OUTPUT_NAME}"
    fi

    echo "Building ${OUTPUT_NAME}"
    cd "$SRC_PATH/${PACKAGE}/main"
    go build -o "${BUILD_PATH}/${OUTPUT_NAME}/${OUTPUT_NAME}"
done
