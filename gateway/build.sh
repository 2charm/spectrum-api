#!/bin/bash
GOOS=linux go build
chmod +x gateway
docker build -t 2charm/gateway .
go clean