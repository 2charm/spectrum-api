#!/bin/bash
docker login
docker rm -f gateway

#Set env variables
export ADDR=:443
export REDISADDR="redis_server:6379"
export NEWSADDR=:80
export APIKEY="`cat ./news_api.key`"
export SESSIONKEY="keykey"
export SQLADDR=3306
export MYSQL_ROOT_PASSWORD="sqlkey"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(sql_server:$SQLADDR)/mysql"
export TLSCERT="/etc/letsencrypt/live/api.spectrumnews.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.spectrumnews.me/privkey.pem"

#Gateway
docker pull 2charm/gateway
docker run -d --network service_network --name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=$ADDR \
-e NEWSADDR="news_service$NEWSADDR" \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
-e DSN=$DSN \
2charm/gateway