#!/usr/bin/env bash
echo "building go server for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a
docker build -t newera/gateway .
go clean