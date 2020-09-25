#!/bin/bash
APP_DIR=$GOPATH/src/github.com/Armatorix/BitLyke
set -a
. $APP_DIR/configs/env.local.api
set +a

go run $APP_DIR/cmd/bitlyke/main.go