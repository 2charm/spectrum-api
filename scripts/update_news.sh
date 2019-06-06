#!/bin/bash
docker login
docker rm -f news_service

#Set env variables
export NEWSADDR=:80
export APIKEY="`cat ./news_api.key`"
export SQLADDR=3306
export MYSQL_ROOT_PASSWORD="sqlkey"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(sql_server:$SQLADDR)/mysql"

#News Service
docker pull 2charm/news_service
docker run -d --network service_network --name news_service \
-p 80:80 \
-e ADDR=$NEWSADDR \
-e DSN=$DSN \
-e APIKEY=$APIKEY \
2charm/news_service