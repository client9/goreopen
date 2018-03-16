#!/bin/sh
set -ex 
go build ./...
gometalinter \
    --vendor \
    --vendored-linters \
    --deadline=60s \
    --disable-all \
    --enable=goimports \
    --enable=vetshadow \
    --enable=varcheck \
    --enable=structcheck \
    --enable=deadcode \
    --enable=ineffassign \
    --enable=unconvert \
    --enable=goconst \
    --enable=golint \
    --enable=gosimple \
    --enable=gofmt \
    --enable=misspell \
    --enable=staticcheck \
    .
go test

