#!/bin/bash
GOOS=linux go build -o ../cmd/gateway/gateway ../cmd/gateway
chmod +x ../cmd/gateway/gateway
docker build -t 2charm/gateway ../cmd/gateway
go clean ../cmd/gateway