#!/bin/bash
docker login
docker volume prune -f
docker rm -f gateway
docker rm -f redis_server
docker rm -f sql_server
docker network rm service_network
docker network create service_network

#Set env variables
export ADDR=:443
export REDISADDR="redis_server:6379"
export APIKEY="`cat ./news_api.key`"
export SESSIONKEY="keykey"
export SQLADDR=3306
export MYSQL_ROOT_PASSWORD="sqlkey"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(sql_server:$SQLADDR)/mysql"
export TLSCERT="/etc/letsencrypt/live/api.spectrumnews.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.spectrumnews.me/privkey.pem"

#Redis Server
docker run -d --network service_network --name redis_server redis
#SQL Server
docker pull 2charm/sql
docker run -d --network service_network --name sql_server \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=mysql \
2charm/sql

#Ensure server is up and running before api is running
sleep 20

#Gateway
docker pull 2charm/gateway
docker run -d --network service_network --name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=$ADDR \
-e REDISADDR=$REDISADDR \
-e APIKEY=$APIKEY \
-e SESSIONKEY=$SESSIONKEY \
-e DSN=$DSN \
2charm/gateway