#!/bin/bash
GOOS=linux go build -o ../cmd/news/news ../cmd/news
chmod +x ../cmd/news/news
docker build -t 2charm/news_service ../cmd/news
go clean ../cmd/news