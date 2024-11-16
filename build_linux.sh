#!/usr/bin/env sh

GOOS=linux GOARCH=amd64 go build -o bin/make_index make_index/main.go
GOOS=linux GOARCH=amd64 go build -o bin/dedup dedup/main.go
