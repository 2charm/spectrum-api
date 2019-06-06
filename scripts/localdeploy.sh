#!/bin/bash
export ADDR=':443'
export KEY=$(cat '../news_api.key')
go run main.go
