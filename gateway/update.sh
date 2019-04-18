#!/bin/bash
docker login
docker volume prune -f
docker rm -f gateway
docker network rm service_network
docker network create service_network
#Set env variables
export ADDR=:443
export TLSCERT="/etc/letsencrypt/live/api.spectrumnews.me/fullchain.pem"
export TLSKEY="/etc/letsencrypt/live/api.spectrumnews.me/privkey.pem"
export KEY="`cat ./news_api.key`"

docker pull 2charm/gateway
docker run -d --network service_network --name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e ADDR=$ADDR \
-e KEY=$KEY \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
2charm/gateway