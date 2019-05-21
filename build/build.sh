#!/bin/bash
GOOS=linux go build ./...
chmod +x gateway
docker build -t 2charm/gateway .
go clean
docker build -t 2charm/sql ../db