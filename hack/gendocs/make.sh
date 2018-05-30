#!/usr/bin/env bash

pushd $GOPATH/src/github.com/appscode/messenger/hack/gendocs
go run main.go
popd
