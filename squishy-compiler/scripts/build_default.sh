#!/bin/bash

APP_NAME="squishy-compiler"
PROJECT_ROOT="./cmd/$APP_NAME"
OUTPUT_DIR="bin"

VERSION="1.0.0"

DO_ZIP="y"
DO_CHECKSUM="y"

./scripts/build.sh <<EOF
$APP_NAME
$PROJECT_ROOT
$OUTPUT_DIR
$VERSION
""
$DO_ZIP
$DO_CHECKSUM
EOF