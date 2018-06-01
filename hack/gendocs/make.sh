#!/usr/bin/env bash

pushd $GOPATH/src/github.com/kubeware/messenger/hack/gendocs
go run main.go
popd
